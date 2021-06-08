package app

import (
	contentserviceapi "contentservice/api/contentservice"
	"github.com/CuriosityMusicStreaming/ComponentsPool/pkg/app/auth"
	"github.com/google/uuid"
)

func contentTests(serviceApiFacade *contentServiceApiFacade, container UserContainer) {
	addContent(serviceApiFacade, container)
	manageContent(serviceApiFacade, container)
}

func addContent(serviceApiFacade *contentServiceApiFacade, container UserContainer) {
	author := auth.UserDescriptor{UserID: uuid.New()}
	listener := auth.UserDescriptor{UserID: uuid.New()}

	container.AddAuthor(author)
	container.AddListener(listener)

	{
		contentTitle := "new song"
		contentType := contentserviceapi.ContentType_Song
		contentAvailabilityType := contentserviceapi.ContentAvailabilityType_Public
		_, err := serviceApiFacade.AddContent(
			contentTitle,
			contentType,
			contentAvailabilityType,
			listener,
		)

		assertEqual(err, ErrOnlyAuthorCanCreateContent)
	}

	{
		firstContentTitle := "new song"
		firstContentType := contentserviceapi.ContentType_Song
		firstContentAvailabilityType := contentserviceapi.ContentAvailabilityType_Public

		secondContentTitle := "new podcast"
		secondContentType := contentserviceapi.ContentType_Podcast
		secondContentAvailabilityType := contentserviceapi.ContentAvailabilityType_Public

		addContentResp, err := serviceApiFacade.AddContent(
			firstContentTitle,
			firstContentType,
			firstContentAvailabilityType,
			author,
		)
		assertNoErr(err)

		firstContentID := addContentResp.ContentID

		contentsResp, err := serviceApiFacade.GetAuthorContent(author)
		assertNoErr(err)

		assertEqual(1, len(contentsResp.Contents))
		content := contentsResp.Contents[0]
		assertEqual(firstContentID, content.ContentID)
		assertEqual(firstContentTitle, content.Name)
		assertEqual(firstContentType, content.Type)
		assertEqual(firstContentAvailabilityType, content.AvailabilityType)

		addContentResp, err = serviceApiFacade.AddContent(
			secondContentTitle,
			secondContentType,
			secondContentAvailabilityType,
			author,
		)
		assertNoErr(err)

		secondContentID := addContentResp.ContentID

		contentsResp, err = serviceApiFacade.GetAuthorContent(author)
		assertNoErr(err)

		assertEqual(2, len(contentsResp.Contents))

		assertNoErr(serviceApiFacade.DeleteContent(author, firstContentID))
		assertNoErr(serviceApiFacade.DeleteContent(author, secondContentID))
	}

	container.Clear()
}

func manageContent(serviceApiFacade *contentServiceApiFacade, container UserContainer) {
	author := auth.UserDescriptor{UserID: uuid.New()}
	anotherAuthor := auth.UserDescriptor{UserID: uuid.New()}

	container.AddAuthor(author)
	container.AddAuthor(anotherAuthor)

	{
		firstContentTitle := "new song"
		firstContentType := contentserviceapi.ContentType_Song
		firstContentAvailabilityType := contentserviceapi.ContentAvailabilityType_Public

		addContentResp, err := serviceApiFacade.AddContent(
			firstContentTitle,
			firstContentType,
			firstContentAvailabilityType,
			author,
		)
		assertNoErr(err)

		firstContentID := addContentResp.ContentID

		//Error cause author owns first content, not anotherAuthor
		assertEqual(serviceApiFacade.DeleteContent(anotherAuthor, firstContentID), ErrOnlyAuthorCanManageContent)

		assertNoErr(serviceApiFacade.SetContentAvailabilityType(author, firstContentID, contentserviceapi.ContentAvailabilityType_Private))

		//Error cause anotherAuthor cannot manage firstContent
		assertEqual(serviceApiFacade.SetContentAvailabilityType(anotherAuthor, firstContentID, contentserviceapi.ContentAvailabilityType_Private), ErrOnlyAuthorCanManageContent)

		assertNoErr(serviceApiFacade.DeleteContent(author, firstContentID))
	}

	{
		resp, err := serviceApiFacade.AddContent(
			"new song",
			contentserviceapi.ContentType_Song,
			contentserviceapi.ContentAvailabilityType_Public,
			author,
		)
		assertNoErr(err)

		contentID := resp.ContentID

		assertNoErr(serviceApiFacade.SetContentAvailabilityType(author, contentID, contentserviceapi.ContentAvailabilityType_Private))

		contentsResp, err := serviceApiFacade.GetContentList([]string{contentID})
		assertNoErr(err)

		assertEqual(0, len(contentsResp.Contents))
	}
}
