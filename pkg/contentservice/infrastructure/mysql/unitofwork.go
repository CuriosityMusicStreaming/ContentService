package mysql

import (
	"contentservice/pkg/contentservice/app/service"
	"contentservice/pkg/contentservice/domain"
	"contentservice/pkg/contentservice/infrastructure/mysql/repository"
	"github.com/CuriosityMusicStreaming/ComponentsPool/pkg/infrastructure/mysql"
	"github.com/pkg/errors"
)

func NewUnitOfFactory(client mysql.TransactionalClient) service.UnitOfWorkFactory {
	return &unitOfWorkFactory{client: client}
}

type unitOfWorkFactory struct {
	client mysql.TransactionalClient
}

func (factory *unitOfWorkFactory) NewUnitOfWork(_ string) (service.UnitOfWork, error) {
	transaction, err := factory.client.BeginTransaction()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &unitOfWork{transaction: transaction}, nil
}

type unitOfWork struct {
	transaction mysql.Transaction
}

func (u *unitOfWork) ContentRepository() domain.ContentRepository {
	return repository.NewContentRepository(u.transaction)
}

func (u *unitOfWork) Complete(err error) error {
	if err != nil {
		err2 := u.transaction.Rollback()
		if err2 != nil {
			return errors.Wrap(err, err2.Error())
		}
	}

	return errors.WithStack(u.transaction.Commit())
}
