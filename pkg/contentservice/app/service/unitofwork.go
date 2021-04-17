package service

import "contentservice/pkg/contentservice/domain"

type UnitOfWorkFactory interface {
	NewUnitOfWork(lockName string) (UnitOfWork, error)
}

type RepositoryProvider interface {
	ContentRepository() domain.ContentRepository
}

type UnitOfWork interface {
	RepositoryProvider
	Complete(err error) error
}
