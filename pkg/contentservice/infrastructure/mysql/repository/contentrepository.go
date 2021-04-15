package repository

import (
	"contentservice/pkg/contentservice/domain"
	"database/sql"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

func NewContentRepository(db *sqlx.DB) domain.ContentRepository {
	return &contentRepository{
		client: db,
	}
}

type contentRepository struct {
	client *sqlx.DB
}

func (repo *contentRepository) NewID() domain.ContentID {
	return domain.ContentID(uuid.New())
}

func (repo *contentRepository) Find(contentID domain.ContentID) (domain.Content, error) {
	const selectSql = `SELECT * from content WHERE content_id = ?`

	binaryUUID, err := uuid.UUID(contentID).MarshalBinary()
	if err != nil {
		return domain.Content{}, err
	}

	var content sqlxContent

	err = repo.client.Get(&content, selectSql, binaryUUID)
	if err != nil {
		if err == sql.ErrNoRows {
			return domain.Content{}, domain.ErrContentNotFound
		}
		return domain.Content{}, errors.WithStack(err)
	}

	return domain.Content{
		ID:          domain.ContentID(content.ContentID),
		Name:        content.Name,
		ContentType: domain.ContentType(content.Type),
	}, nil
}

func (repo *contentRepository) Store(content domain.Content) error {
	const insertSql = `
		INSERT INTO content (content_id, name, author_id, type, availability_type) VALUES(?, ?, ?, ?, ?)
		ON DUPLICATE KEY 
		UPDATE content_id=VALUES(content_id), name=VALUES(name), author_id=VALUES(author_id), type=VALUES(type), availability_type=VALUES(availability_type)
	`

	binaryUUID, err := uuid.UUID(content.ID).MarshalBinary()
	if err != nil {
		return errors.WithStack(err)
	}

	_, err = repo.client.Exec(insertSql, binaryUUID, content.Name, content.ContentType)
	return err
}

func (repo *contentRepository) Remove(contentID domain.ContentID) error {
	const deleteSql = `DELETE FROM content WHERE content_id = ?`

	binaryUUID, err := uuid.UUID(contentID).MarshalBinary()
	if err != nil {
		return errors.WithStack(err)
	}

	_, err = repo.client.Exec(deleteSql, binaryUUID)
	return err

}

type sqlxContent struct {
	ContentID uuid.UUID `db:"content_id"`
	Name      string    `db:"name"`
	Type      int       `db:"type"`
}
