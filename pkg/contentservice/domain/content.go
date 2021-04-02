package domain

import (
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type ContentID uuid.UUID

type ContentType int

const (
	ContentTypeSong ContentType = iota
	ContentTypePodcast
)

type Content struct {
	ID   ContentID
	Name string
	ContentType
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
