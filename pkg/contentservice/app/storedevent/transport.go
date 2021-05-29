package storedevent

type Transport interface {
	Send(msgBody string, eventType string) error
}
