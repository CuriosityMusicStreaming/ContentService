package app

import (
	"context"

	"github.com/CuriosityMusicStreaming/ComponentsPool/pkg/app/auth"
	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	contentserviceapi "contentservice/api/contentservice"
)

type contentServiceAPIFacade struct {
	client     contentserviceapi.ContentServiceClient
	serializer auth.UserDescriptorSerializer
}

func (facade *contentServiceAPIFacade) AddContent(
	title string,
	contentType contentserviceapi.ContentType,
	availabilityType contentserviceapi.ContentAvailabilityType,
	userDescriptor auth.UserDescriptor,
) (*contentserviceapi.AddContentResponse, error) {
	userToken, err := facade.serializer.Serialize(userDescriptor)
	assertNoErr(err)

	resp, err := facade.client.AddContent(context.Background(), &contentserviceapi.AddContentRequest{
		Name:             title,
		Type:             contentType,
		AvailabilityType: availabilityType,
		UserToken:        userToken,
	})
	return resp, transformError(err)
}

func (facade *contentServiceAPIFacade) GetAuthorContent(userDescriptor auth.UserDescriptor) (*contentserviceapi.GetAuthorContentResponse, error) {
	userToken, err := facade.serializer.Serialize(userDescriptor)
	assertNoErr(err)

	resp, err := facade.client.GetAuthorContent(context.Background(), &contentserviceapi.GetAuthorContentRequest{
		UserToken: userToken,
	})
	return resp, transformError(err)
}

func (facade *contentServiceAPIFacade) GetContentList(contentIDs []string) (*contentserviceapi.GetContentListResponse, error) {
	return facade.client.GetContentList(context.Background(), &contentserviceapi.GetContentListRequest{
		ContentIDs: contentIDs,
	})
}

func (facade contentServiceAPIFacade) DeleteContent(userDescriptor auth.UserDescriptor, contentID string) error {
	userToken, err := facade.serializer.Serialize(userDescriptor)
	assertNoErr(err)

	_, err = facade.client.DeleteContent(context.Background(), &contentserviceapi.DeleteContentRequest{
		ContentID: contentID,
		UserToken: userToken,
	})
	return transformError(err)
}

func (facade contentServiceAPIFacade) SetContentAvailabilityType(
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
	return transformError(err)
}

func transformError(err error) error {
	s, ok := status.FromError(err)
	if ok {
		switch s.Code() {
		case codes.InvalidArgument:
			return ErrOnlyAuthorCanCreateContent
		case codes.PermissionDenied:
			return ErrOnlyAuthorCanManageContent
		}
	}
	return err
}

var (
	ErrOnlyAuthorCanCreateContent = errors.New("only author can create content")
	ErrOnlyAuthorCanManageContent = errors.New("only author can manage content")
)
