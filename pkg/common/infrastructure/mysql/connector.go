package mysql

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/cenkalti/backoff"
	_ "github.com/go-sql-driver/mysql" // provides MySQL driver
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file" // provides filesystem source
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

const (
	dbDriverName = "mysql"

	maxReconnectWaitingTime = 15 * time.Second
)

type DSN struct {
	User     string
	Password string
	Host     string
	Database string
}

func (dsn *DSN) String() string {
	return fmt.Sprintf("%s:%s@(%s)/%s?charset=utf8mb4", dsn.User, dsn.Password, dsn.Host, dsn.Database)
}

type Connector interface {
	Open(dsn DSN, maxConnections int) error
	MigrateUp(dsn DSN, migrationsDir string) error
	Client() *sqlx.DB
	Close() error
}

type connector struct {
	db *sqlx.DB
}

func NewConnector() Connector {
	return &connector{}
}

func (c *connector) MigrateUp(dsn DSN, migrationsDir string) error {
	// Db connections will be closed when migration object is closed, so new connection must be opened
	db, err := openDb(dsn, 1)
	if err != nil {
		return errors.WithStack(err)
	}

	m, err := createMigrator(db, migrationsDir)
	if err != nil {
		return errors.WithStack(err)
	}
	defer m.Close()

	err = m.Up()
	if err == migrate.ErrNoChange {
		return nil
	}

	return errors.Wrap(err, "failed to migrate")
}

func (c *connector) Open(dsn DSN, maxConnections int) error {
	var err error
	c.db, err = openDb(dsn, maxConnections)
	return errors.WithStack(err)
}

func (c *connector) Close() error {
	err := c.db.Close()
	return errors.Wrap(err, "failed to disconnect")
}

func (c *connector) Client() *sqlx.DB {
	return c.db
}

func createMigrator(db *sqlx.DB, migrationsDir string) (*migrate.Migrate, error) {
	_, err := os.Stat(migrationsDir)
	if err != nil {
		return nil, errors.Wrapf(err, "cannot use migrations from %s", migrationsDir)
	}
	migrationsDir, err = filepath.Abs(migrationsDir)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	migrationsURL := fmt.Sprintf("file://%s", migrationsDir)

	var driver database.Driver
	err = backoff.Retry(func() error {
		var tryError error
		driver, tryError = mysql.WithInstance(db.DB, &mysql.Config{})
		return tryError
	}, newExponentialBackOff())
	if err != nil {
		return nil, errors.Wrapf(err, "cannot create migrations driver")
	}

	m, err := migrate.NewWithDatabaseInstance(migrationsURL, dbDriverName, driver)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create migrator")
	}
	return m, nil
}

func openDb(dsn DSN, maxConnections int) (*sqlx.DB, error) {
	db, err := sqlx.Open(dbDriverName, dsn.String())
	if err != nil {
		return nil, errors.Wrapf(err, "failed to open database")
	}

	// Limit max connections count,
	//  next goroutine will wait once reached limit.
	db.SetMaxOpenConns(maxConnections)

	err = backoff.Retry(func() error {
		tryError := db.Ping()
		return tryError
	}, newExponentialBackOff())
	if err != nil {
		dbCloseErr := db.Close()
		if dbCloseErr != nil {
			err = errors.Wrap(err, dbCloseErr.Error())
		}
		return nil, errors.Wrapf(err, "failed to ping database")
	}
	return db, errors.WithStack(err)
}

func newExponentialBackOff() *backoff.ExponentialBackOff {
	exponentialBackOff := backoff.NewExponentialBackOff()
	exponentialBackOff.MaxElapsedTime = maxReconnectWaitingTime
	return exponentialBackOff
}
