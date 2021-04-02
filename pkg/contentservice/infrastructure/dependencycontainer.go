package infrastructure

import (
	"contentservice/pkg/contentservice/app/service"
	"contentservice/pkg/contentservice/domain"
	"contentservice/pkg/contentservice/infrastructure/mysql/repository"
)

type DependencyContainer interface {
	ContentService() service.ContentService
}

func NewDependencyContainer() DependencyContainer {
	return &dependencyContainer{}
}

type dependencyContainer struct {
}

func (container *dependencyContainer) ContentService() service.ContentService {
	return service.NewContentService(container.domainContentService())
}

func (container *dependencyContainer) domainContentService() domain.ContentService {
	return domain.NewContentService(container.contentRepository())
}

func (container *dependencyContainer) contentRepository() domain.ContentRepository {
	return repository.NewContentRepository()
}
