package domain

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

//nolint
func TestContentService_AddContent(t *testing.T) {
	mockRepo := newMockRepo()
	mockEventDispatcher := newMockEventDispatcher()

	contentService := NewContentService(mockRepo, mockEventDispatcher)

	{
		contentName := "song"
		authorID := uuid.New()
		contentType := ContentTypeSong
		availabilityType := ContentAvailabilityTypePublic

		contentID, err := contentService.AddContent(contentName, AuthorID(authorID), contentType, availabilityType)
		assert.NoError(t, err)

		assert.Equal(t, len(mockRepo.contents), 1, "content must be added to repo")

		content := mockRepo.contents[contentID]

		assert.Equal(t, content.Title, contentName)
		assert.Equal(t, content.AuthorID, AuthorID(authorID))
		assert.Equal(t, content.ContentType, contentType)
		assert.Equal(t, content.AvailabilityType, availabilityType)
	}
}

//nolint
func TestContentService_SetContentAvailabilityType(t *testing.T) {
	mockRepo := newMockRepo()
	mockEventDispatcher := newMockEventDispatcher()

	contentService := NewContentService(mockRepo, mockEventDispatcher)

	contentName := "song"
	authorID := AuthorID(uuid.New())
	anotherUserID := AuthorID(uuid.New())
	contentType := ContentTypeSong
	availabilityType := ContentAvailabilityTypePublic

	contentID, err := contentService.AddContent(contentName, authorID, contentType, availabilityType)
	assert.NoError(t, err)

	{
		newAvailabilityType := ContentAvailabilityTypePrivate
		err := contentService.SetContentAvailabilityType(contentID, anotherUserID, newAvailabilityType)
		assert.Error(t, err, ErrOnlyAuthorCanManageContent.Error())
	}

	mockEventDispatcher.clear()

	{
		newAvailabilityType := ContentAvailabilityTypePrivate
		err := contentService.SetContentAvailabilityType(contentID, authorID, newAvailabilityType)
		assert.NoError(t, err)

		content := mockRepo.contents[contentID]
		assert.Equal(t, content.AvailabilityType, newAvailabilityType)

		assert.Equal(t, 1, len(mockEventDispatcher.events), "changing availability type dispatches event")

		assert.Equal(t, "content_availability_type_changed", mockEventDispatcher.events[0].ID())
	}
}

//nolint
func TestContentService_DeleteContent(t *testing.T) {
	mockRepo := newMockRepo()
	mockEventDispatcher := newMockEventDispatcher()

	contentService := NewContentService(mockRepo, mockEventDispatcher)

	contentName := "song"
	authorID := AuthorID(uuid.New())
	anotherUserID := AuthorID(uuid.New())
	contentType := ContentTypeSong
	availabilityType := ContentAvailabilityTypePublic

	contentID, err := contentService.AddContent(contentName, authorID, contentType, availabilityType)
	assert.NoError(t, err)

	{
		err := contentService.DeleteContent(contentID, anotherUserID)
		assert.Error(t, err, ErrOnlyAuthorCanManageContent.Error())
	}

	{
		err := contentService.DeleteContent(contentID, authorID)
		assert.NoError(t, err)

		assert.Equal(t, len(mockRepo.contents), 0)
	}

	{
		err := contentService.DeleteContent(ContentID(uuid.New()), authorID)
		assert.Error(t, err, ErrContentNotFound.Error())
	}
}

func newMockRepo() *mockContentRepository {
	return &mockContentRepository{contents: make(map[ContentID]Content)}
}

type mockContentRepository struct {
	contents map[ContentID]Content
}

func (repo *mockContentRepository) NewID() ContentID {
	return ContentID(uuid.New())
}

func (repo *mockContentRepository) Find(contentID ContentID) (Content, error) {
	content, ok := repo.contents[contentID]
	if !ok {
		return Content{}, ErrContentNotFound
	}
	return content, nil
}

func (repo *mockContentRepository) Store(content Content) error {
	repo.contents[content.ID] = content

	return nil
}

func (repo *mockContentRepository) Remove(contentID ContentID) error {
	delete(repo.contents, contentID)
	return nil
}

func newMockEventDispatcher() *mockEventDispatcher {
	return &mockEventDispatcher{}
}

type mockEventDispatcher struct {
	events []Event
}

func (eventDispatcher *mockEventDispatcher) Dispatch(event Event) error {
	eventDispatcher.events = append(eventDispatcher.events, event)

	return nil
}

func (eventDispatcher *mockEventDispatcher) clear() {
	eventDispatcher.events = nil
}
