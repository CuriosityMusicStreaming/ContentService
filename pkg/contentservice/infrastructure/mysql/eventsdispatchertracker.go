package mysql

import (
	"contentservice/pkg/contentservice/app/storedevent"
	"database/sql"
	"github.com/CuriosityMusicStreaming/ComponentsPool/pkg/infrastructure/mysql"
	"github.com/google/uuid"
	"sync"
)

func NewEventsDispatchTracker(client mysql.Client) *eventsDispatchTracker {
	return &eventsDispatchTracker{client: client}
}

type eventsDispatchTracker struct {
	client mysql.Client
	mutex  sync.Mutex
}

func (tracker *eventsDispatchTracker) TrackLastID(transportName string, id storedevent.ID) error {
	const insertQuery = `
		INSERT INTO tracked_stored_event (transport_name, last_stored_event_id, created_at) VALUES (?, ?, now())
		ON DUPLICATE KEY UPDATE last_stored_event_id=VALUES(last_stored_event_id)
	`
	binaryID, err := uuid.UUID(id).MarshalBinary()
	if err != nil {
		return err
	}

	_, err = tracker.client.Exec(insertQuery, transportName, binaryID)
	return err
}

func (tracker *eventsDispatchTracker) LastId(transportName string) (*storedevent.ID, error) {
	const selectQuery = `SELECT last_stored_event_id FROM tracked_stored_event WHERE transport_name = ?`

	var id uuid.UUID
	err := tracker.client.Get(&id, selectQuery, transportName)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		} else {
			return nil, err
		}
	}

	ID := storedevent.ID(id)
	return &ID, nil
}

func (tracker *eventsDispatchTracker) Lock() error {
	tracker.mutex.Lock()

	return nil
}

func (tracker *eventsDispatchTracker) Unlock() error {
	tracker.mutex.Unlock()

	return nil
}
