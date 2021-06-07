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
	title string,
	contentType contentserviceapi.ContentType,
	availabilityType contentserviceapi.ContentAvailabilityType,
	userDescriptor auth.UserDescriptor,
) (*contentserviceapi.AddContentResponse, error) {
	userToken, err := facade.serializer.Serialize(userDescriptor)
	assertNoErr(err)

	return facade.client.AddContent(context.Background(), &contentserviceapi.AddContentRequest{
		Name:             title,
		Type:             contentType,
		AvailabilityType: availabilityType,
		UserToken:        userToken,
	})
}

func (facade *contentServiceApiFacade) GetAuthorContent(userDescriptor auth.UserDescriptor) (*contentserviceapi.GetAuthorContentResponse, error) {
	userToken, err := facade.serializer.Serialize(userDescriptor)
	assertNoErr(err)

	return facade.client.GetAuthorContent(context.Background(), &contentserviceapi.GetAuthorContentRequest{
		UserToken: userToken,
	})
}

func (facade *contentServiceApiFacade) GetContentList(contentIDs []string) (*contentserviceapi.GetContentListResponse, error) {
	return facade.client.GetContentList(context.Background(), &contentserviceapi.GetContentListRequest{
		ContentIDs: contentIDs,
	})
}

func (facade contentServiceApiFacade) DeleteContent(userDescriptor auth.UserDescriptor, contentID string) error {
	userToken, err := facade.serializer.Serialize(userDescriptor)
	assertNoErr(err)

	_, err = facade.client.DeleteContent(context.Background(), &contentserviceapi.DeleteContentRequest{
		ContentID: contentID,
		UserToken: userToken,
	})
	return err
}

func (facade contentServiceApiFacade) SetContentAvailabilityType(
	userDescriptor auth.UserDescriptor,
	contentID string,
	contentAvailabilityType contentserviceapi.ContentAvailabilityType,
) error {
	userToken, err := facade.serializer.Serialize(userDescriptor)
	assertNoErr(err)

	_, err = facade.client.SetContentAvailabilityType(context.Background(), &contentserviceapi.SetContentAvailabilityTypeRequest{
		ContentID:                  contentID,
		NewContentAvailabilityType: contentAvailabilityType,
		UserToken:                  userToken,
	})
	return err
}
