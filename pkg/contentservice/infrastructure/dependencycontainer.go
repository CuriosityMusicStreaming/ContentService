package infrastructure

import (
	commonauth "github.com/CuriosityMusicStreaming/ComponentsPool/pkg/app/auth"
	"github.com/CuriosityMusicStreaming/ComponentsPool/pkg/app/logger"
	commonstoredevent "github.com/CuriosityMusicStreaming/ComponentsPool/pkg/app/storedevent"
	commonmysql "github.com/CuriosityMusicStreaming/ComponentsPool/pkg/infrastructure/mysql"

	"contentservice/api/authorizationservice"
	"contentservice/pkg/contentservice/app/auth"
	"contentservice/pkg/contentservice/app/query"
	"contentservice/pkg/contentservice/app/service"
	"contentservice/pkg/contentservice/app/storedevent"
	"contentservice/pkg/contentservice/domain"
	"contentservice/pkg/contentservice/infrastructure/mysql"
	infrastructurequery "contentservice/pkg/contentservice/infrastructure/mysql/query"
	"contentservice/pkg/contentservice/infrastructure/transport/client"
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
	authorizationServiceClient authorizationservice.AuthorizationServiceClient,
	eventStore commonstoredevent.Store,
	storedEventSenderCallback mysql.UnitOfWorkCompleteNotifier,
) DependencyContainer {

	userDescriptorSerializer := userDescriptorSerializer()

	unitOfWorkFactory, notifier := unitOfWorkFactory(client)

	notifier.subscribe(storedEventSenderCallback)

	return &dependencyContainer{
		contentService: contentService(
			unitOfWorkFactory,
			eventDispatcher(eventStore),
			authorizationService(
				authorizationServiceClient,
				userDescriptorSerializer,
			),
		),
		trustedContentQueryService: trustedContentQueryService(client),
		userDescriptorSerializer:   userDescriptorSerializer,
	}
}

type completeNotifier struct {
	subscribers []mysql.UnitOfWorkCompleteNotifier
}

func (notifier *completeNotifier) subscribe(subscriber mysql.UnitOfWorkCompleteNotifier) {
	notifier.subscribers = append(notifier.subscribers, subscriber)
}

func (notifier *completeNotifier) onComplete() {
	for _, subscriber := range notifier.subscribers {
		subscriber()
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

func unitOfWorkFactory(client commonmysql.TransactionalClient) (service.UnitOfWorkFactory, *completeNotifier) {
	notifier := &completeNotifier{}

	return mysql.NewNotifyingUnitOfWorkFactory(
		mysql.NewUnitOfFactory(client),
		notifier.onComplete,
	), notifier
}

func eventDispatcher(store commonstoredevent.Store) domain.EventDispatcher {
	eventPublisher := domain.NewEventPublisher()

	{
		handler := commonstoredevent.NewStoredDomainEventHandler(store, storedevent.NewEventSerializer())
		eventPublisher.Subscribe(domain.HandlerFunc(func(event domain.Event) error {
			return handler.Handle(event)
		}))
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
	authorizationServiceClient authorizationservice.AuthorizationServiceClient,
	userDescriptorSerializer commonauth.UserDescriptorSerializer,
) auth.AuthorizationService {
	return client.NewAuthorizationService(
		authorizationServiceClient,
		userDescriptorSerializer,
	)
}
