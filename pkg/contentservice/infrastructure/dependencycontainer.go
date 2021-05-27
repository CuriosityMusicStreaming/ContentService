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
	TrustedContentQueryService() query.ContentQueryService
	AuthorizedContentQueryService(userDescriptor commonauth.UserDescriptor) query.ContentQueryService
	UserDescriptorSerializer() commonauth.UserDescriptorSerializer
}

func NewDependencyContainer(
	client commonmysql.TransactionalClient,
	logger logger.Logger,
	userServiceClient userserviceapi.UserServiceClient,
) DependencyContainer {

	userDescriptorSerializer := userDescriptorSerializer()

	return &dependencyContainer{
		contentService: contentService(
			unitOfWorkFactory(client),
			eventDispatcher(logger),
			authorizationService(
				userServiceClient,
				userDescriptorSerializer,
			),
		),
		trustedContentQueryService: trustedContentQueryService(client),
		userDescriptorSerializer:   userDescriptorSerializer,
	}
}

type dependencyContainer struct {
	contentService             service.ContentService
	trustedContentQueryService query.ContentQueryService
	userDescriptorSerializer   commonauth.UserDescriptorSerializer
}

func (container *dependencyContainer) ContentService() service.ContentService {
	return container.contentService
}

func (container *dependencyContainer) TrustedContentQueryService() query.ContentQueryService {
	return container.trustedContentQueryService
}

func (container *dependencyContainer) AuthorizedContentQueryService(userDescriptor commonauth.UserDescriptor) query.ContentQueryService {
	return query.NewAuthorizedContentQueryService(container.TrustedContentQueryService(), userDescriptor)
}

func (container *dependencyContainer) UserDescriptorSerializer() commonauth.UserDescriptorSerializer {
	return container.userDescriptorSerializer
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

func contentService(
	unitOfWork service.UnitOfWorkFactory,
	eventDispatcher domain.EventDispatcher,
	authorizationService auth.AuthorizationService,
) service.ContentService {
	return service.NewContentService(
		unitOfWork,
		eventDispatcher,
		authorizationService,
	)
}

func trustedContentQueryService(client commonmysql.TransactionalClient) query.ContentQueryService {
	return infrastructurequery.NewContentQueryService(client)
}

func userDescriptorSerializer() commonauth.UserDescriptorSerializer {
	return commonauth.NewUserDescriptorSerializer()
}

func authorizationService(
	userServiceClient userserviceapi.UserServiceClient,
	userDescriptorSerializer commonauth.UserDescriptorSerializer,
) auth.AuthorizationService {
	return userserviceadapter.NewAuthorizationService(
		userServiceClient,
		userDescriptorSerializer,
	)
}
