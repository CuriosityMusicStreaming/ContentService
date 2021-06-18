package app

import (
	"github.com/CuriosityMusicStreaming/ComponentsPool/pkg/app/auth"
	"github.com/google/uuid"

	contentserviceapi "contentservice/api/contentservice"
)

func contentTests(serviceAPIFacade *contentServiceAPIFacade, container UserContainer) {
	addContent(serviceAPIFacade, container)
	manageContent(serviceAPIFacade, container)
}

func addContent(serviceAPIFacade *contentServiceAPIFacade, container UserContainer) {
	author := auth.UserDescriptor{UserID: uuid.New()}
	listener := auth.UserDescriptor{UserID: uuid.New()}

	container.AddAuthor(author)
	container.AddListener(listener)

	const contentTitle = "new song"

	{
		contentType := contentserviceapi.ContentType_Song
		contentAvailabilityType := contentserviceapi.ContentAvailabilityType_Public
		_, err := serviceAPIFacade.AddContent(
			contentTitle,
			contentType,
			contentAvailabilityType,
			listener,
		)

		assertEqual(err, ErrOnlyAuthorCanCreateContent)
	}

	{
		firstContentTitle := contentTitle
		firstContentType := contentserviceapi.ContentType_Song
		firstContentAvailabilityType := contentserviceapi.ContentAvailabilityType_Public

		secondContentTitle := "new podcast"
		secondContentType := contentserviceapi.ContentType_Podcast
		secondContentAvailabilityType := contentserviceapi.ContentAvailabilityType_Public

		addContentResp, err := serviceAPIFacade.AddContent(
			firstContentTitle,
			firstContentType,
			firstContentAvailabilityType,
			author,
		)
		assertNoErr(err)

		firstContentID := addContentResp.ContentID

		contentsResp, err := serviceAPIFacade.GetAuthorContent(author)
		assertNoErr(err)

		assertEqual(1, len(contentsResp.Contents))
		content := contentsResp.Contents[0]
		assertEqual(firstContentID, content.ContentID)
		assertEqual(firstContentTitle, content.Name)
		assertEqual(firstContentType, content.Type)
		assertEqual(firstContentAvailabilityType, content.AvailabilityType)

		addContentResp, err = serviceAPIFacade.AddContent(
			secondContentTitle,
			secondContentType,
			secondContentAvailabilityType,
			author,
		)
		assertNoErr(err)

		secondContentID := addContentResp.ContentID

		contentsResp, err = serviceAPIFacade.GetAuthorContent(author)
		assertNoErr(err)

		assertEqual(2, len(contentsResp.Contents))

		assertNoErr(serviceAPIFacade.DeleteContent(author, firstContentID))
		assertNoErr(serviceAPIFacade.DeleteContent(author, secondContentID))
	}

	container.Clear()
}

func manageContent(serviceAPIFacade *contentServiceAPIFacade, container UserContainer) {
	author := auth.UserDescriptor{UserID: uuid.New()}
	anotherAuthor := auth.UserDescriptor{UserID: uuid.New()}

	container.AddAuthor(author)
	container.AddAuthor(anotherAuthor)

	{
		firstContentTitle := "new song"
		firstContentType := contentserviceapi.ContentType_Song
		firstContentAvailabilityType := contentserviceapi.ContentAvailabilityType_Public

		addContentResp, err := serviceAPIFacade.AddContent(
			firstContentTitle,
			firstContentType,
			firstContentAvailabilityType,
			author,
		)
		assertNoErr(err)

		firstContentID := addContentResp.ContentID

		// Error cause author owns first content, not anotherAuthor
		assertEqual(serviceAPIFacade.DeleteContent(anotherAuthor, firstContentID), ErrOnlyAuthorCanManageContent)

		assertNoErr(serviceAPIFacade.SetContentAvailabilityType(author, firstContentID, contentserviceapi.ContentAvailabilityType_Private))

		// Error cause anotherAuthor cannot manage firstContent
		assertEqual(serviceAPIFacade.SetContentAvailabilityType(anotherAuthor, firstContentID, contentserviceapi.ContentAvailabilityType_Private), ErrOnlyAuthorCanManageContent)

		assertNoErr(serviceAPIFacade.DeleteContent(author, firstContentID))
	}

	{
		resp, err := serviceAPIFacade.AddContent(
			"new song",
			contentserviceapi.ContentType_Song,
			contentserviceapi.ContentAvailabilityType_Public,
			author,
		)
		assertNoErr(err)

		contentID := resp.ContentID

		assertNoErr(serviceAPIFacade.SetContentAvailabilityType(author, contentID, contentserviceapi.ContentAvailabilityType_Private))

		contentsResp, err := serviceAPIFacade.GetContentList([]string{contentID})
		assertNoErr(err)

		assertEqual(0, len(contentsResp.Contents))
	}
}
