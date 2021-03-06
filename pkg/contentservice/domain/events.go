package domain

type Event interface {
	ID() string
}

type EventHandler interface {
	Handle(event Event) error
}

type HandlerFunc func(event Event) error

func (f HandlerFunc) Handle(event Event) error {
	return f(event)
}

type EventDispatcher interface {
	Dispatch(event Event) error
}

type EventSource interface {
	Subscribe(handler EventHandler)
}

type EventPublisher interface {
	EventDispatcher
	EventSource
}

func NewEventPublisher() EventPublisher {
	return &eventPublisher{}
}

type eventPublisher struct {
	subscribers []EventHandler
}

func (e *eventPublisher) Dispatch(event Event) error {
	for _, subscriber := range e.subscribers {
		err := subscriber.Handle(event)
		if err != nil {
			return err
		}
	}

	return nil
}

func (e *eventPublisher) Subscribe(handler EventHandler) {
	e.subscribers = append(e.subscribers, handler)
}

type ContentAdded struct {
	ContentID ContentID
	AuthorID  AuthorID
}

func (e ContentAdded) ID() string {
	return "content_added"
}

type ContentDeleted struct {
	ContentID ContentID
}

func (e ContentDeleted) ID() string {
	return "content_deleted"
}

type ContentContentAvailabilityTypeChanged struct {
	ContentID                  ContentID
	NewContentAvailabilityType ContentAvailabilityType
}

func (e ContentContentAvailabilityTypeChanged) ID() string {
	return "content_availability_type_changed"
}
