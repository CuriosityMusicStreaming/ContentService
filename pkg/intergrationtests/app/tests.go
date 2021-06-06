package app

import (
	contentserviceapi "contentservice/api/contentservice"
	"github.com/CuriosityMusicStreaming/ComponentsPool/pkg/app/auth"
	"golang.org/x/net/context"
)

type UserContainer interface {
	AddAuthor(descriptor auth.UserDescriptor)
	AddListener(descriptor auth.UserDescriptor)

	Clear()
}

func RunTests(contentServiceClient contentserviceapi.ContentServiceClient, container UserContainer) {
	ContentTests(&contentServiceApiFacade{
		client:     contentServiceClient,
		serializer: auth.NewUserDescriptorSerializer(),
	}, container)
}

type contentServiceApiFacade struct {
	client     contentserviceapi.ContentServiceClient
	serializer auth.UserDescriptorSerializer
}

func (facade *contentServiceApiFacade) AddContent(
	name string,
	contentType contentserviceapi.ContentType,
	availabilityType contentserviceapi.ContentAvailabilityType,
	userDescriptor auth.UserDescriptor,
) (*contentserviceapi.AddContentResponse, error) {
	userToken, err := facade.serializer.Serialize(userDescriptor)
	assertErr(err)

	return facade.client.AddContent(context.Background(), &contentserviceapi.AddContentRequest{
		Name:             name,
		Type:             contentType,
		AvailabilityType: availabilityType,
		UserToken:        userToken,
	})
}
