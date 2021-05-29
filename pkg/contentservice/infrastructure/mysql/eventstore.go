package mysql

import (
	"contentservice/pkg/contentservice/app/storedevent"
	"fmt"
	"github.com/CuriosityMusicStreaming/ComponentsPool/pkg/infrastructure/mysql"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

func NewEventStore(client mysql.Client) storedevent.Store {
	return &eventStore{client: client}
}

type eventStore struct {
	client mysql.Client
}

func (store *eventStore) Append(event storedevent.StoredEvent) error {
	const insertSql = `INSERT INTO stored_event (stored_event_id, type, body, created_at) VALUES(?, ?, ?, now())`
	binaryID, err := uuid.UUID(event.ID).MarshalBinary()
	if err != nil {
		return errors.WithStack(err)
	}

	_, err = store.client.Exec(insertSql, binaryID, event.Type, event.Body)
	return errors.WithStack(err)
}

func (store *eventStore) GetAllAfter(id *storedevent.ID) ([]storedevent.StoredEvent, error) {
	selectSql := `SELECT * FROM stored_event WHERE %s ORDER BY created_at`
	var args []interface{}

	if id != nil {
		selectSql = fmt.Sprintf(selectSql, "created_at > (SELECT created_at FROM stored_event WHERE stored_event_id = ?)")
		binaryID, err := uuid.UUID(*id).MarshalBinary()
		if err != nil {
			return nil, errors.WithStack(err)
		}
		args = append(args, binaryID)
	}

	var storedEvents []sqlxStoredEvent

	err := store.client.Select(&storedEvents, selectSql, args...)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	result := make([]storedevent.StoredEvent, 0, len(storedEvents))
	for _, event := range storedEvents {
		result = append(result, storedevent.StoredEvent{
			ID:   storedevent.ID(event.ID),
			Type: event.Type,
			Body: event.Body,
		})
	}

	return result, nil
}

type sqlxStoredEvent struct {
	ID   uuid.UUID
	Type string
	Body string
}
