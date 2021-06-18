package mysql

import (
	"github.com/CuriosityMusicStreaming/ComponentsPool/pkg/infrastructure/mysql"
	"github.com/pkg/errors"

	"contentservice/pkg/contentservice/app/service"
	"contentservice/pkg/contentservice/domain"
	"contentservice/pkg/contentservice/infrastructure/mysql/repository"
)

func NewUnitOfFactory(client mysql.TransactionalClient) service.UnitOfWorkFactory {
	return &unitOfWorkFactory{client: client}
}

type unitOfWorkFactory struct {
	client mysql.TransactionalClient
}

func (factory *unitOfWorkFactory) NewUnitOfWork(lockName string) (service.UnitOfWork, error) {
	transaction, err := factory.client.BeginTransaction()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var lock *mysql.Lock
	if lockName != "" {
		l := mysql.NewLock(factory.client, lockName)
		lock = &l
		err = lock.Lock()
		if err != nil {
			return nil, errors.Wrap(transaction.Rollback(), err.Error())
		}
	}

	return &unitOfWork{transaction: transaction, lock: lock}, nil
}

type unitOfWork struct {
	transaction mysql.Transaction
	lock        *mysql.Lock
}

func (u *unitOfWork) ContentRepository() domain.ContentRepository {
	return repository.NewContentRepository(u.transaction)
}

func (u *unitOfWork) Complete(err error) error {
	if u.lock != nil {
		lockErr := u.lock.Unlock()
		if err != nil {
			if lockErr != nil {
				err = errors.Wrap(err, lockErr.Error())
			}
		} else {
			err = lockErr
		}
	}

	if err != nil {
		err2 := u.transaction.Rollback()
		if err2 != nil {
			return errors.Wrap(err, err2.Error())
		}
	}

	return errors.WithStack(u.transaction.Commit())
}
