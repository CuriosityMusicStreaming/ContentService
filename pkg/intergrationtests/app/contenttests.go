package app

import (
	contentserviceapi "contentservice/api/contentservice"
	"github.com/CuriosityMusicStreaming/ComponentsPool/pkg/app/auth"
	"github.com/google/uuid"
)

func ContentTests(serviceApiFacade *contentServiceApiFacade, container UserContainer) {
	addContent(serviceApiFacade, container)
}

func addContent(serviceApiFacade *contentServiceApiFacade, container UserContainer) {
	author := auth.UserDescriptor{UserID: uuid.New()}
	listener := auth.UserDescriptor{UserID: uuid.New()}

	container.AddAuthor(author)
	container.AddListener(listener)

	{
		_, err := serviceApiFacade.AddContent(
			"new song",
			contentserviceapi.ContentType_Song,
			contentserviceapi.ContentAvailabilityType_Public,
			author,
		)
		assertErr(err)
	}

	container.Clear()
}
