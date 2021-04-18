package auth

import "github.com/CuriosityMusicStreaming/ComponentsPool/pkg/app/auth"

type AuthorizationService interface {
	CanAddContent(descriptor auth.UserDescriptor) (bool, error)
}
