package storedevent

import (
	"contentservice/pkg/contentservice/domain"
	"encoding/json"
	"github.com/CuriosityMusicStreaming/ComponentsPool/pkg/app/storedevent"
	commondomain "github.com/CuriosityMusicStreaming/ComponentsPool/pkg/domain"
	"github.com/google/uuid"
)

func NewEventSerializer() storedevent.EventSerializer {
	return &eventSerializer{}
}

type eventSerializer struct {
}

type eventBody struct {
	Type    string
	Payload *json.RawMessage
}

func (serializer *eventSerializer) Serialize(event commondomain.Event) (string, error) {
	payload, err := serializeAsJSON(event)
	if err != nil {
		return "", err
	}

	payloadRawMessage := json.RawMessage(payload)
	body := eventBody{
		Type:    event.ID(),
		Payload: &payloadRawMessage,
	}

	messageBody, err := json.Marshal(body)

	return string(messageBody), err
}

func serializeAsJSON(event commondomain.Event) ([]byte, error) {
	return json.Marshal(serializeEvent(event))
}

func serializeEvent(event commondomain.Event) (eventPayload interface{}) {
	switch currEvent := event.(type) {
	case domain.ContentContentAvailabilityTypeChanged:
		eventPayload = struct {
			ContentID                  uuid.UUID `json:"content_id"`
			NewContentAvailabilityType int       `json:"new_content_availability_type"`
		}{
			ContentID:                  uuid.UUID(currEvent.ContentID),
			NewContentAvailabilityType: int(currEvent.NewContentAvailabilityType),
		}
	case domain.ContentDeleted:
		eventPayload = struct {
			ContentID uuid.UUID `json:"content_id"`
		}{ContentID: uuid.UUID(currEvent.ContentID)}
	case domain.ContentAdded:
		eventPayload = struct {
			ContentID uuid.UUID `json:"content_id"`
			AuthorID  uuid.UUID `json:"author_id"`
		}{
			ContentID: uuid.UUID(currEvent.ContentID),
			AuthorID:  uuid.UUID(currEvent.AuthorID),
		}
	}
	return
}
