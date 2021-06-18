package query

import (
	"github.com/google/uuid"

	"contentservice/pkg/contentservice/app/service"
)

type ContentView struct {
	ID               uuid.UUID
	Title            string
	AuthorID         uuid.UUID
	ContentType      service.ContentType
	AvailabilityType service.ContentAvailabilityType
}

type ContentSpecification struct {
	ContentIDs []uuid.UUID
	AuthorIDs  []uuid.UUID
}

type ContentQueryService interface {
	ContentList(spec ContentSpecification) ([]ContentView, error)
}
