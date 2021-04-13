package infrastructure

import (
	"contentservice/pkg/contentservice/app/query"
	"contentservice/pkg/contentservice/app/service"
	"contentservice/pkg/contentservice/domain"
	"contentservice/pkg/contentservice/infrastructure/integration"
	infrastructurequery "contentservice/pkg/contentservice/infrastructure/mysql/query"
	"contentservice/pkg/contentservice/infrastructure/mysql/repository"
	"github.com/CuriosityMusicStreaming/ComponentsPool/pkg/app/auth"
	"github.com/CuriosityMusicStreaming/ComponentsPool/pkg/app/logger"
	"github.com/jmoiron/sqlx"
)

type DependencyContainer interface {
	ContentService() service.ContentService
	ContentQueryService() query.ContentQueryService
	UserDescriptorSerializer() auth.UserDescriptorSerializer
}

func NewDependencyContainer(client *sqlx.DB, logger logger.Logger) DependencyContainer {
	return &dependencyContainer{
		client:          client,
		logger:          logger,
		eventDispatcher: eventDispatcher(logger),
	}
}

type dependencyContainer struct {
	client          *sqlx.DB
	logger          logger.Logger
	eventDispatcher domain.EventDispatcher
}

func (container *dependencyContainer) ContentService() service.ContentService {
	return service.NewContentService(container.domainContentService())
}

func (container *dependencyContainer) ContentQueryService() query.ContentQueryService {
	return infrastructurequery.NewContentQueryService(container.client)
}

func (container dependencyContainer) UserDescriptorSerializer() auth.UserDescriptorSerializer {
	return auth.NewUserDescriptorSerializer()
}

func (container *dependencyContainer) domainContentService() domain.ContentService {
	return domain.NewContentService(container.contentRepository(), container.eventDispatcher)
}

func (container *dependencyContainer) contentRepository() domain.ContentRepository {
	return repository.NewContentRepository(container.client)
}

func eventDispatcher(logger logger.Logger) domain.EventDispatcher {
	eventPublisher := domain.NewEventPublisher()

	{
		handler := integration.NewIntegrationEventHandler(logger)
		eventPublisher.Subscribe(handler)
	}

	return eventPublisher
}
