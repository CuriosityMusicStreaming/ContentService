package app

import (
	"github.com/CuriosityMusicStreaming/ComponentsPool/pkg/app/auth"

	contentserviceapi "contentservice/api/contentservice"
)

type UserContainer interface {
	AddAuthor(descriptor auth.UserDescriptor)
	AddListener(descriptor auth.UserDescriptor)

	Clear()
}

func RunTests(contentServiceClient contentserviceapi.ContentServiceClient, container UserContainer) {
	contentTests(&contentServiceAPIFacade{
		client:     contentServiceClient,
		serializer: auth.NewUserDescriptorSerializer(),
	}, container)
}
