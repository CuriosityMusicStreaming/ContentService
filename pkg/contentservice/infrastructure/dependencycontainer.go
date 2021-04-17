package infrastructure

import (
	"contentservice/pkg/contentservice/app/query"
	"contentservice/pkg/contentservice/app/service"
	"contentservice/pkg/contentservice/domain"
	"contentservice/pkg/contentservice/infrastructure/integration"
	"contentservice/pkg/contentservice/infrastructure/mysql"
	infrastructurequery "contentservice/pkg/contentservice/infrastructure/mysql/query"
	"github.com/CuriosityMusicStreaming/ComponentsPool/pkg/app/auth"
	"github.com/CuriosityMusicStreaming/ComponentsPool/pkg/app/logger"
	commonmysql "github.com/CuriosityMusicStreaming/ComponentsPool/pkg/infrastructure/mysql"
)

type DependencyContainer interface {
	ContentService() service.ContentService
	ContentQueryService() query.ContentQueryService
	UserDescriptorSerializer() auth.UserDescriptorSerializer
}

func NewDependencyContainer(client commonmysql.TransactionalClient, logger logger.Logger) DependencyContainer {
	return &dependencyContainer{
		client:            client,
		logger:            logger,
		eventDispatcher:   eventDispatcher(logger),
		unitOfWorkFactory: unitOfWorkFactory(client),
	}
}

type dependencyContainer struct {
	client            commonmysql.TransactionalClient
	logger            logger.Logger
	eventDispatcher   domain.EventDispatcher
	unitOfWorkFactory service.UnitOfWorkFactory
}

func (container *dependencyContainer) ContentService() service.ContentService {
	return service.NewContentService(container.unitOfWorkFactory, container.eventDispatcher)
}

func (container *dependencyContainer) ContentQueryService() query.ContentQueryService {
	return infrastructurequery.NewContentQueryService(container.client)
}

func (container *dependencyContainer) UserDescriptorSerializer() auth.UserDescriptorSerializer {
	return auth.NewUserDescriptorSerializer()
}

func unitOfWorkFactory(client commonmysql.TransactionalClient) service.UnitOfWorkFactory {
	return mysql.NewUnitOfFactory(client)
}

func eventDispatcher(logger logger.Logger) domain.EventDispatcher {
	eventPublisher := domain.NewEventPublisher()

	{
		handler := integration.NewIntegrationEventHandler(logger)
		eventPublisher.Subscribe(handler)
	}

	return eventPublisher
}
