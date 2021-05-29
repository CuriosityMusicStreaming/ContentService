package transport

import "github.com/CuriosityMusicStreaming/ComponentsPool/pkg/app/logger"

func NewMockIntegrationTransport(logger logger.Logger) *mockIntegrationTransport {
	return &mockIntegrationTransport{logger: logger}
}

type mockIntegrationTransport struct {
	logger logger.Logger
}

func (m *mockIntegrationTransport) Name() string {
	return "mock_transport"
}

func (m *mockIntegrationTransport) Send(eventType string, msgBody string) error {
	m.logger.WithFields(map[string]interface{}{
		"event_type": eventType,
		"body":       msgBody,
	}).Info("event sent")

	return nil
}
