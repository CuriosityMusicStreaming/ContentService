package storedevent

import (
	"contentservice/pkg/contentservice/domain"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
)

type EventSerializer interface {
	Serialize(event domain.Event) (string, error)
}

func NewEventSerializer() EventSerializer {
	return &eventSerializer{}
}

type eventSerializer struct {
}

type eventBody struct {
	Type    string
	Payload *json.RawMessage
}

func (serializer *eventSerializer) Serialize(event domain.Event) (string, error) {
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

	fmt.Println("BODY", body)

	return string(messageBody), err
}

func serializeAsJSON(event domain.Event) ([]byte, error) {
	return json.Marshal(serializeEvent(event))
}

func serializeEvent(event domain.Event) (eventPayload interface{}) {
	switch currEvent := event.(type) {
	case domain.ContentContentAvailabilityTypeChanged:
		eventPayload = struct {
			ContentID uuid.UUID `json:"content_id"`
		}{ContentID: uuid.UUID(currEvent.ContentID)}
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
