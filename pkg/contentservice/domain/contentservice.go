package domain

import (
	"github.com/pkg/errors"
)

var (
	ErrOnlyAuthorCanDeleteContent = errors.New("only author can delete content")
	ErrOnlyAuthorCanMaangeContent = errors.New("only author can manage content")
)

type ContentService interface {
	AddContent(name string, authorID AuthorID, contentType ContentType, availabilityType ContentAvailabilityType) error
	DeleteContent(contentID ContentID, authorID AuthorID) error
	SetContentAvailabilityType(contentID ContentID, authorID AuthorID, availabilityType ContentAvailabilityType) error
}

func NewContentService(repository ContentRepository) ContentService {
	return &contentService{
		repo: repository,
	}
}

type contentService struct {
	repo ContentRepository
}

func (service *contentService) AddContent(name string, authorID AuthorID, contentType ContentType, availabilityType ContentAvailabilityType) error {
	id := service.repo.NewID()
	err := service.repo.Store(Content{
		ID:               id,
		Name:             name,
		AuthorID:         authorID,
		ContentType:      contentType,
		AvailabilityType: availabilityType,
	})
	if err != nil {
		return err
	}

	return nil
}

func (service *contentService) DeleteContent(contentID ContentID, authorID AuthorID) error {
	content, err := service.repo.Find(contentID)
	if err != nil {
		return err
	}

	if content.AuthorID != authorID {
		return ErrOnlyAuthorCanDeleteContent
	}

	return service.repo.Remove(content.ID)
}

func (service *contentService) SetContentAvailabilityType(contentID ContentID, authorID AuthorID, availabilityType ContentAvailabilityType) error {
	content, err := service.repo.Find(contentID)
	if err != nil {
		return err
	}

	if content.AuthorID != authorID {
		return ErrOnlyAuthorCanMaangeContent
	}

	if content.AvailabilityType == availabilityType {
		return nil
	}

	content.AvailabilityType = availabilityType

	err = service.repo.Store(content)

	return err
}
