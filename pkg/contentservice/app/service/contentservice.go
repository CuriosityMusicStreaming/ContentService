package service

import (
	"contentservice/pkg/contentservice/app/auth"
	"contentservice/pkg/contentservice/domain"
	commonauth "github.com/CuriosityMusicStreaming/ComponentsPool/pkg/app/auth"
	"github.com/google/uuid"
)

type ContentType int

const (
	ContentTypeSong    = ContentType(domain.ContentTypeSong)
	ContentTypePodcast = ContentType(domain.ContentTypePodcast)
)

type ContentAvailabilityType int

const (
	ContentAvailabilityTypePublic  = ContentAvailabilityType(domain.ContentAvailabilityTypePublic)
	ContentAvailabilityTypePrivate = ContentAvailabilityType(domain.ContentAvailabilityTypePrivate)
)

type ContentService interface {
	AddContent(name string, userDescriptor commonauth.UserDescriptor, contentType ContentType, availabilityType ContentAvailabilityType) error
	DeleteContent(contentID uuid.UUID, userDescriptor commonauth.UserDescriptor) error
	SetContentAvailabilityType(contentID uuid.UUID, descriptor commonauth.UserDescriptor, availabilityType ContentAvailabilityType) error
}

func NewContentService(factory UnitOfWorkFactory, eventDispatcher domain.EventDispatcher, authorizationService auth.AuthorizationService) ContentService {
	return &contentService{
		unitOfWorkFactory:    factory,
		eventDispatcher:      eventDispatcher,
		authorizationService: authorizationService,
	}
}

type contentService struct {
	unitOfWorkFactory    UnitOfWorkFactory
	eventDispatcher      domain.EventDispatcher
	authorizationService auth.AuthorizationService
}

func (service *contentService) AddContent(name string, userDescriptor commonauth.UserDescriptor, contentType ContentType, availabilityType ContentAvailabilityType) error {
	if canAdd, err := service.authorizationService.CanAddContent(userDescriptor); !canAdd || err != nil {
		return err
	}

	return service.executeInUnitOfWork(func(provider RepositoryProvider) error {
		domainService := domain.NewContentService(provider.ContentRepository(), service.eventDispatcher)

		_, err := domainService.AddContent(name, domain.AuthorID(userDescriptor.UserID), domain.ContentType(contentType), domain.ContentAvailabilityType(availabilityType))
		return err
	})
}

func (service *contentService) DeleteContent(contentID uuid.UUID, userDescriptor commonauth.UserDescriptor) error {
	return service.executeInUnitOfWork(func(provider RepositoryProvider) error {
		domainService := domain.NewContentService(provider.ContentRepository(), service.eventDispatcher)

		return domainService.DeleteContent(domain.ContentID(contentID), domain.AuthorID(userDescriptor.UserID))
	})
}

func (service *contentService) SetContentAvailabilityType(contentID uuid.UUID, userDescriptor commonauth.UserDescriptor, availabilityType ContentAvailabilityType) error {
	return service.executeInUnitOfWork(func(provider RepositoryProvider) error {
		domainService := domain.NewContentService(provider.ContentRepository(), service.eventDispatcher)

		return domainService.SetContentAvailabilityType(
			domain.ContentID(contentID),
			domain.AuthorID(userDescriptor.UserID),
			domain.ContentAvailabilityType(availabilityType),
		)
	})
}

func (service *contentService) executeInUnitOfWork(f func(provider RepositoryProvider) error) error {
	unitOfWork, err := service.unitOfWorkFactory.NewUnitOfWork("")
	if err != nil {
		return err
	}
	defer func() {
		err = unitOfWork.Complete(err)
	}()
	err = f(unitOfWork)
	return err
}
