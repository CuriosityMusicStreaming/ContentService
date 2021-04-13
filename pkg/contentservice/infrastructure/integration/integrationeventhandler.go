package integration

import (
	"contentservice/pkg/contentservice/domain"
	"fmt"
	"github.com/CuriosityMusicStreaming/ComponentsPool/pkg/app/logger"
)

func NewIntegrationEventHandler(logger logger.Logger) domain.EventHandler {
	return &integrationEventListener{logger: logger}
}

type integrationEventListener struct {
	logger logger.Logger
}

func (handler *integrationEventListener) Handle(event domain.Event) error {
	handler.logger.Info(fmt.Sprintf("event sended to queue %s", event.ID()))
	return nil
}
