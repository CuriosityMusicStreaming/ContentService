package query

import (
	"contentservice/pkg/contentservice/app/service"
	"github.com/google/uuid"
)

type ContentView struct {
	ID               uuid.UUID
	Name             string
	AuthorID         uuid.UUID
	ContentType      service.ContentType
	AvailabilityType service.ContentAvailabilityType
}

type ContentSpecification struct {
	ContentIDs []uuid.UUID
}

type ContentQueryService interface {
	ContentList(spec ContentSpecification) ([]ContentView, error)
}
