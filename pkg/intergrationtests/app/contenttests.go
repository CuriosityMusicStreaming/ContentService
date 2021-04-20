package app

import (
	contentserviceapi "contentservice/api/contentservice"
	userserviceapi "contentservice/api/userservice"
	"context"
	"github.com/CuriosityMusicStreaming/ComponentsPool/pkg/app/auth"
	"github.com/google/uuid"
)

func ContentTests(contentServiceClient contentserviceapi.ContentServiceClient, userServiceClient userserviceapi.UserServiceClient) {
	ctx := context.Background()
	userDescriptorSerializer := auth.NewUserDescriptorSerializer()

	addContent(ctx, contentServiceClient, userServiceClient, userDescriptorSerializer)
}

func addContent(
	ctx context.Context,
	contentServiceClient contentserviceapi.ContentServiceClient,
	userServiceClient userserviceapi.UserServiceClient,
	serializer auth.UserDescriptorSerializer,
) {
	{
		creator, err := newCreator(userServiceClient)
		assertErr(err)

		userToken, err := serializer.Serialize(creator)
		assertErr(err)

		_, err = contentServiceClient.AddContent(ctx, &contentserviceapi.AddContentRequest{
			Name:             "new song",
			Type:             contentserviceapi.ContentType_Song,
			AvailabilityType: contentserviceapi.ContentAvailabilityType_Public,
			UserToken:        userToken,
		})
		assertErr(err)
	}
}

func newCreator(userServiceClient userserviceapi.UserServiceClient) (auth.UserDescriptor, error) {
	ctx := context.Background()
	resp, err := userServiceClient.AddUser(ctx, &userserviceapi.AddUserRequest{
		Email:    "root@mail,com",
		Password: "12345Q",
		Role:     userserviceapi.UserRole_CREATOR,
	})
	if err != nil {
		return auth.UserDescriptor{}, err
	}

	creatorID, err := uuid.Parse(resp.UserId)
	if err != nil {
		return auth.UserDescriptor{}, err
	}

	return auth.UserDescriptor{UserID: creatorID}, err
}

func newListener(userServiceClient userserviceapi.UserServiceClient) (auth.UserDescriptor, error) {
	ctx := context.Background()
	resp, err := userServiceClient.AddUser(ctx, &userserviceapi.AddUserRequest{
		Email:    "root@mail,com",
		Password: "12345Q",
		Role:     userserviceapi.UserRole_LISTENER,
	})
	if err != nil {
		return auth.UserDescriptor{}, err
	}

	creatorID, err := uuid.Parse(resp.UserId)
	if err != nil {
		return auth.UserDescriptor{}, err
	}

	return auth.UserDescriptor{UserID: creatorID}, err
}
