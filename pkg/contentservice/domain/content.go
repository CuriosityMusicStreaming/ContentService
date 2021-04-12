package domain

import (
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type ContentID uuid.UUID

type AuthorID uuid.UUID

type ContentType int

const (
	ContentTypeSong ContentType = iota
	ContentTypePodcast
)

type ContentAvailabilityType int

const (
	ContentAvailabilityTypePublic ContentAvailabilityType = iota
	ContentAvailabilityTypePrivate
)

type Content struct {
	ID       ContentID
	Name     string
	AuthorID AuthorID
	ContentType
	AvailabilityType ContentAvailabilityType
}

type ContentRepository interface {
	NewID() ContentID
	Find(contentID ContentID) (Content, error)
	Store(content Content) error
	Remove(contentID ContentID) error
}

var (
	ErrContentNotFound = errors.New("content not found")
)
