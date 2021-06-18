package query

import (
	"github.com/CuriosityMusicStreaming/ComponentsPool/pkg/app/auth"

	appservice "contentservice/pkg/contentservice/app/service"
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

	//nolint:prealloc
	var filteredContentViews []ContentView
	for _, view := range views {
		if view.AvailabilityType == appservice.ContentAvailabilityTypePrivate && view.AuthorID != service.userDescriptor.UserID {
			continue
		}
		filteredContentViews = append(filteredContentViews, view)
	}
	return filteredContentViews, nil
}
