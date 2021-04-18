package userserviceadapter

import (
	userserviceapi "contentservice/api/userservice"
	"contentservice/pkg/contentservice/app/auth"
	"context"
	commonauth "github.com/CuriosityMusicStreaming/ComponentsPool/pkg/app/auth"
)

func NewAuthorizationService(
	userServiceApi userserviceapi.UserServiceClient,
	userDescriptorSerializer commonauth.UserDescriptorSerializer,
) auth.AuthorizationService {
	return &userServiceAuthorizationService{
		userServiceApi:           userServiceApi,
		userDescriptorSerializer: userDescriptorSerializer,
	}
}

type userServiceAuthorizationService struct {
	userServiceApi           userserviceapi.UserServiceClient
	userDescriptorSerializer commonauth.UserDescriptorSerializer
}

func (service *userServiceAuthorizationService) CanAddContent(descriptor commonauth.UserDescriptor) (bool, error) {
	userToken, err := service.userDescriptorSerializer.Serialize(descriptor)
	if err != nil {
		return false, err
	}

	ctx := context.Background()
	resp, err := service.userServiceApi.CanAddContent(ctx, &userserviceapi.CanAddContentRequest{UserToken: userToken})
	if err != nil {
		return false, err
	}

	return resp.CanAdd, err
}
