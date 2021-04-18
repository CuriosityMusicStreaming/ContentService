package domain

import (
	"github.com/pkg/errors"
)

var (
	ErrOnlyAuthorCanDeleteContent = errors.New("only author can delete content")
	ErrOnlyAuthorCanManageContent = errors.New("only author can manage content")
)

type ContentService interface {
	AddContent(name string, authorID AuthorID, contentType ContentType, availabilityType ContentAvailabilityType) (ContentID, error)
	DeleteContent(contentID ContentID, authorID AuthorID) error
	SetContentAvailabilityType(contentID ContentID, authorID AuthorID, availabilityType ContentAvailabilityType) error
}

func NewContentService(repository ContentRepository, dispatcher EventDispatcher) ContentService {
	return &contentService{
		repo:            repository,
		eventDispatcher: dispatcher,
	}
}

type contentService struct {
	repo            ContentRepository
	eventDispatcher EventDispatcher
}

func (service *contentService) AddContent(name string, authorID AuthorID, contentType ContentType, availabilityType ContentAvailabilityType) (ContentID, error) {
	id := service.repo.NewID()
	err := service.repo.Store(Content{
		ID:               id,
		Name:             name,
		AuthorID:         authorID,
		ContentType:      contentType,
		AvailabilityType: availabilityType,
	})
	if err != nil {
		return ContentID{}, err
	}

	return id, nil
}

func (service *contentService) DeleteContent(contentID ContentID, authorID AuthorID) error {
	content, err := service.repo.Find(contentID)
	if err != nil {
		return err
	}

	if content.AuthorID != authorID {
		return ErrOnlyAuthorCanDeleteContent
	}

	err = service.repo.Remove(content.ID)
	if err != nil {
		return err
	}

	return service.eventDispatcher.Dispatch(ContentDeleted{ContentID: contentID})
}

func (service *contentService) SetContentAvailabilityType(contentID ContentID, authorID AuthorID, availabilityType ContentAvailabilityType) error {
	content, err := service.repo.Find(contentID)
	if err != nil {
		return err
	}

	if content.AuthorID != authorID {
		return ErrOnlyAuthorCanManageContent
	}

	if content.AvailabilityType == availabilityType {
		return nil
	}

	content.AvailabilityType = availabilityType

	err = service.repo.Store(content)
	if err != nil {
		return err
	}

	return service.eventDispatcher.Dispatch(ContentContentAvailabilityTypeChanged{ContentID: contentID})
}
