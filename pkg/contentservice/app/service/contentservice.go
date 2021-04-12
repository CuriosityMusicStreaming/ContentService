package service

import (
	"contentservice/pkg/contentservice/domain"
	"github.com/CuriosityMusicStreaming/ComponentsPool/pkg/app/auth"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

var (
	ErrOnlyCreatorCanAddContent    = errors.New("only creator can add content")
	ErrOnlyCreatorCanManageContent = errors.New("only creator can manage content")
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
	AddContent(name string, userDescriptor auth.UserDescriptor, contentType ContentType, availabilityType ContentAvailabilityType) error
	DeleteContent(contentID uuid.UUID, userDescriptor auth.UserDescriptor) error
	SetContentAvailabilityType(contentID uuid.UUID, descriptor auth.UserDescriptor, availabilityType ContentAvailabilityType) error
}

func NewContentService(domainService domain.ContentService) ContentService {
	return &contentService{
		domainService: domainService,
	}
}

type contentService struct {
	domainService domain.ContentService
}

func (service *contentService) AddContent(name string, userDescriptor auth.UserDescriptor, contentType ContentType, availabilityType ContentAvailabilityType) error {
	if userDescriptor.Role != auth.Creator {
		return ErrOnlyCreatorCanAddContent
	}

	return service.domainService.AddContent(name, domain.AuthorID(userDescriptor.UserID), domain.ContentType(contentType), domain.ContentAvailabilityType(availabilityType))
}

func (service *contentService) DeleteContent(contentID uuid.UUID, userDescriptor auth.UserDescriptor) error {
	if userDescriptor.Role != auth.Creator {
		return ErrOnlyCreatorCanManageContent
	}
	return service.domainService.DeleteContent(domain.ContentID(contentID), domain.AuthorID(userDescriptor.UserID))
}

func (service *contentService) SetContentAvailabilityType(contentID uuid.UUID, userDescriptor auth.UserDescriptor, availabilityType ContentAvailabilityType) error {
	if userDescriptor.Role != auth.Creator {
		return ErrOnlyCreatorCanManageContent
	}
	return service.domainService.SetContentAvailabilityType(
		domain.ContentID(contentID),
		domain.AuthorID(userDescriptor.UserID),
		domain.ContentAvailabilityType(availabilityType),
	)
}
