package query

import (
	appservice "contentservice/pkg/contentservice/app/service"
	"github.com/CuriosityMusicStreaming/ComponentsPool/pkg/app/auth"
)

func NewAuthorizedContentQueryService(queryService ContentQueryService, userDescriptor auth.UserDescriptor) ContentQueryService {
	return &authorizedContentQueryService{queryService: queryService, userDescriptor: userDescriptor}
}

type authorizedContentQueryService struct {
	queryService   ContentQueryService
	userDescriptor auth.UserDescriptor
}

func (service *authorizedContentQueryService) ContentList(spec ContentSpecification) ([]ContentView, error) {
	views, err := service.queryService.ContentList(spec)
	if err != nil {
		return nil, err
	}

	var filteredContentViews []ContentView
	for _, view := range views {
		if view.AvailabilityType == appservice.ContentAvailabilityTypePrivate && view.AuthorID != service.userDescriptor.UserID {
			continue
		}
		filteredContentViews = append(filteredContentViews, view)
	}
	return filteredContentViews, nil
}
