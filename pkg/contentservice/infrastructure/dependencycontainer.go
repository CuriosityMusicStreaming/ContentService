package infrastructure

import (
	"contentservice/pkg/contentservice/app/service"
	"contentservice/pkg/contentservice/domain"
	"contentservice/pkg/contentservice/infrastructure/mysql/repository"
	"github.com/CuriosityMusicStreaming/ComponentsPool/pkg/app/auth"
	"github.com/jmoiron/sqlx"
)

type DependencyContainer interface {
	ContentService() service.ContentService
	UserDescriptorSerializer() auth.UserDescriptorSerializer
}

func NewDependencyContainer(client *sqlx.DB) DependencyContainer {
	return &dependencyContainer{
		client: client,
	}
}

type dependencyContainer struct {
	client *sqlx.DB
}

func (container *dependencyContainer) ContentService() service.ContentService {
	return service.NewContentService(container.domainContentService())
}

func (container dependencyContainer) UserDescriptorSerializer() auth.UserDescriptorSerializer {
	return auth.NewUserDescriptorSerializer()
}

func (container *dependencyContainer) domainContentService() domain.ContentService {
	return domain.NewContentService(container.contentRepository())
}

func (container *dependencyContainer) contentRepository() domain.ContentRepository {
	return repository.NewContentRepository(container.client)
}
