package domain

type ContentService interface {
	AddContent(name string, contentType ContentType) error
}

func NewContentService(repository ContentRepository) ContentService {
	return &contentService{
		repo: repository,
	}
}

type contentService struct {
	repo ContentRepository
}

func (c *contentService) AddContent(name string, contentType ContentType) error {
	id := c.repo.NewID()
	err := c.repo.Store(Content{
		ID:          id,
		Name:        name,
		ContentType: contentType,
	})
	if err != nil {
		return err
	}

	return nil
}
