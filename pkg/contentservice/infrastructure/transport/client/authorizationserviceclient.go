package client

import (
	"context"

	"contentservice/api/authorizationservice"
	"contentservice/pkg/contentservice/app/auth"

	commonauth "github.com/CuriosityMusicStreaming/ComponentsPool/pkg/app/auth"
)

func NewAuthorizationService(
	authorizationServiceClient authorizationservice.AuthorizationServiceClient,
	userDescriptorSerializer commonauth.UserDescriptorSerializer,
) auth.AuthorizationService {
	return &userServiceAuthorizationService{
		authorizationServiceClient: authorizationServiceClient,
		userDescriptorSerializer:   userDescriptorSerializer,
	}
}

type userServiceAuthorizationService struct {
	authorizationServiceClient authorizationservice.AuthorizationServiceClient
	userDescriptorSerializer   commonauth.UserDescriptorSerializer
}

func (service *userServiceAuthorizationService) CanAddContent(descriptor commonauth.UserDescriptor) (bool, error) {
	userToken, err := service.userDescriptorSerializer.Serialize(descriptor)
	if err != nil {
		return false, err
	}

	ctx := context.Background()
	resp, err := service.authorizationServiceClient.CanAddContent(ctx, &authorizationservice.CanAddContentRequest{UserToken: userToken})
	if err != nil {
		return false, err
	}

	return resp.CanAdd, err
}
