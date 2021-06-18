package integrationevent

import (
	"fmt"

	log "github.com/CuriosityMusicStreaming/ComponentsPool/pkg/app/logger"
)

func NewIntegrationEventHandler(logger log.Logger) Handler {
	return &integrationEventListener{logger: logger}
}

type integrationEventListener struct {
	logger log.Logger
}

func (handler *integrationEventListener) Handle(msgBody string) error {
	handler.logger.Info(fmt.Sprintf("Event received with body %s", msgBody))

	return nil
}
