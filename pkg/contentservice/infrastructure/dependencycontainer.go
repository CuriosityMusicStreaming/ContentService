package infrastructure

import (
	userserviceapi "contentservice/api/userservice"
	"contentservice/pkg/contentservice/app/auth"
	"contentservice/pkg/contentservice/app/query"
	"contentservice/pkg/contentservice/app/service"
	"contentservice/pkg/contentservice/domain"
	"contentservice/pkg/contentservice/infrastructure/integration"
	"contentservice/pkg/contentservice/infrastructure/mysql"
	infrastructurequery "contentservice/pkg/contentservice/infrastructure/mysql/query"
	"contentservice/pkg/contentservice/infrastructure/userserviceadapter"
	commonauth "github.com/CuriosityMusicStreaming/ComponentsPool/pkg/app/auth"
	"github.com/CuriosityMusicStreaming/ComponentsPool/pkg/app/logger"
	commonmysql "github.com/CuriosityMusicStreaming/ComponentsPool/pkg/infrastructure/mysql"
)

type DependencyContainer interface {
	ContentService() service.ContentService
	ContentQueryService() query.ContentQueryService
	UserDescriptorSerializer() commonauth.UserDescriptorSerializer
}

func NewDependencyContainer(
	client commonmysql.TransactionalClient,
	logger logger.Logger,
	userServiceClient userserviceapi.UserServiceClient,
) DependencyContainer {
	return &dependencyContainer{
		client:            client,
		logger:            logger,
		userServiceClient: userServiceClient,
		eventDispatcher:   eventDispatcher(logger),
		unitOfWorkFactory: unitOfWorkFactory(client),
	}
}

type dependencyContainer struct {
	client            commonmysql.TransactionalClient
	logger            logger.Logger
	userServiceClient userserviceapi.UserServiceClient
	eventDispatcher   domain.EventDispatcher
	unitOfWorkFactory service.UnitOfWorkFactory
}

func (container *dependencyContainer) ContentService() service.ContentService {
	return service.NewContentService(
		container.unitOfWorkFactory,
		container.eventDispatcher,
		container.authorizationService(),
	)
}

func (container *dependencyContainer) ContentQueryService() query.ContentQueryService {
	return infrastructurequery.NewContentQueryService(container.client)
}

func (container *dependencyContainer) UserDescriptorSerializer() commonauth.UserDescriptorSerializer {
	return commonauth.NewUserDescriptorSerializer()
}

func (container dependencyContainer) authorizationService() auth.AuthorizationService {
	return userserviceadapter.NewAuthorizationService(
		container.userServiceClient,
		container.UserDescriptorSerializer(),
	)
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
