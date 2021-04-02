package repository

import (
	"contentservice/pkg/contentservice/domain"
	"github.com/google/uuid"
)

func NewContentRepository() domain.ContentRepository {
	return &contentRepository{
		contents: map[domain.ContentID]domain.Content{},
	}
}

type contentRepository struct {
	contents map[domain.ContentID]domain.Content
}

func (repo *contentRepository) NewID() domain.ContentID {
	return domain.ContentID(uuid.New())
}

func (repo *contentRepository) Find(contentID domain.ContentID) (domain.Content, error) {
	content, ok := repo.contents[contentID]
	if !ok {
		return domain.Content{}, domain.ErrContentNotFound
	}

	return content, nil
}

func (repo *contentRepository) Store(content domain.Content) error {
	_, ok := repo.contents[content.ID]
	if ok {
		return nil
	}

	repo.contents[content.ID] = content

	return nil
}

func (repo *contentRepository) Remove(contentID domain.ContentID) error {
	return nil
}
