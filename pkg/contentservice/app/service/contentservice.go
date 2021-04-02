package service

import "contentservice/pkg/contentservice/domain"

type ContentType int

const (
	ContentTypeSong    = ContentType(domain.ContentTypeSong)
	ContentTypePodcast = ContentType(domain.ContentTypePodcast)
)

type ContentService interface {
	AddContent(name string, contentType ContentType) error
}

func NewContentService(domainService domain.ContentService) ContentService {
	return &contentService{
		domainService: domainService,
	}
}

type contentService struct {
	domainService domain.ContentService
}

func (service *contentService) AddContent(name string, contentType ContentType) error {
	return service.domainService.AddContent(name, domain.ContentType(contentType))
}
