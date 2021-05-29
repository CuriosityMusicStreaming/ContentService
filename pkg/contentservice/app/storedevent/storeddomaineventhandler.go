package storedevent

import (
	"contentservice/pkg/contentservice/domain"
	"encoding/json"
	"github.com/google/uuid"
)

type StoredDomainEventHandler struct {
	transport Transport
}

func (handler *StoredDomainEventHandler) Handle(event domain.Event) error {
	payload, err := serializeAsJSON(event)
	if err != nil {
		return err
	}

	storedEvent := Event{
		Type: event.ID(),
		Body: payload,
	}

	//handler.transport.Send()
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
	}
	return
}
