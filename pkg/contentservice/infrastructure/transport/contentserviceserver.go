package transport

import (
	"contentservice/pkg/contentservice/app/query"
	"contentservice/pkg/contentservice/app/service"
	"contentservice/pkg/contentservice/infrastructure"
	"context"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"google.golang.org/protobuf/types/known/emptypb"

	api "contentservice/api/contentservice"
)

func NewContentServiceServer(container infrastructure.DependencyContainer) api.ContentServiceServer {
	return &contentServiceServer{
		container: container,
	}
}

type contentServiceServer struct {
	container infrastructure.DependencyContainer
}

func (server *contentServiceServer) AddContent(_ context.Context, req *api.AddContentRequest) (*emptypb.Empty, error) {
	userDesc, err := server.container.UserDescriptorSerializer().Deserialize(req.UserToken)
	if err != nil {
		return nil, err
	}

	contentType, ok := apiToContentTypeMap[req.Type]
	if !ok {
		return nil, ErrUnknownContentType
	}

	availabilityType, ok := apiToContentAvailabilityTypeMap[req.AvailabilityType]
	if !ok {
		return nil, ErrUnknownContentAvailabilityType
	}

	err = server.container.ContentService().AddContent(req.Name, userDesc, contentType, availabilityType)
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, err
}

func (server *contentServiceServer) DeleteContent(_ context.Context, req *api.DeleteContentRequest) (*emptypb.Empty, error) {
	userDesc, err := server.container.UserDescriptorSerializer().Deserialize(req.UserToken)
	if err != nil {
		return nil, err
	}

	contentID, err := uuid.Parse(req.ContentID)
	if err != nil {
		return nil, err
	}

	err = server.container.ContentService().DeleteContent(contentID, userDesc)
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, err
}

func (server *contentServiceServer) SetContentAvailabilityType(_ context.Context, req *api.SetContentAvailabilityTypeRequest) (*emptypb.Empty, error) {
	userDesc, err := server.container.UserDescriptorSerializer().Deserialize(req.UserToken)
	if err != nil {
		return nil, err
	}

	contentID, err := uuid.Parse(req.ContentID)
	if err != nil {
		return nil, err
	}

	availabilityType, ok := apiToContentAvailabilityTypeMap[req.NewContentAvailabilityType]
	if !ok {
		return nil, ErrUnknownContentAvailabilityType
	}

	err = server.container.ContentService().SetContentAvailabilityType(contentID, userDesc, availabilityType)
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, err
}

func (server *contentServiceServer) GetContentList(_ context.Context, req *api.GetContentListRequest) (*api.GetContentListResponse, error) {
	queryService := server.container.ContentQueryService()

	contentIDs := make([]uuid.UUID, len(req.ContentIDs))
	for _, contentID := range req.ContentIDs {
		contentUUID, err := uuid.Parse(contentID)
		if err != nil {
			return nil, err
		}

		contentIDs = append(contentIDs, contentUUID)
	}

	list, err := queryService.ContentList(query.ContentSpecification{ContentIDs: contentIDs})
	if err != nil {
		return nil, err
	}

	res := make([]*api.Content, len(list))
	for _, contentView := range list {
		res = append(res, &api.Content{
			ContentID:        contentView.ID.String(),
			Name:             contentView.Name,
			AuthorID:         contentView.AuthorID.String(),
			Type:             contentTypeToApiMap[contentView.ContentType],
			AvailabilityType: contentAvailabilityTypeToApiMap[contentView.AvailabilityType],
		})
	}

	return &api.GetContentListResponse{Contents: res}, nil
}

var apiToContentTypeMap = map[api.ContentType]service.ContentType{
	api.ContentType_Song:    service.ContentTypeSong,
	api.ContentType_Podcast: service.ContentTypePodcast,
}

var apiToContentAvailabilityTypeMap = map[api.ContentAvailabilityType]service.ContentAvailabilityType{
	api.ContentAvailabilityType_Public:  service.ContentAvailabilityTypePublic,
	api.ContentAvailabilityType_Private: service.ContentAvailabilityTypePrivate,
}

var contentTypeToApiMap = map[service.ContentType]api.ContentType{
	service.ContentTypeSong:    api.ContentType_Song,
	service.ContentTypePodcast: api.ContentType_Podcast,
}

var contentAvailabilityTypeToApiMap = map[service.ContentAvailabilityType]api.ContentAvailabilityType{
	service.ContentAvailabilityTypePublic:  api.ContentAvailabilityType_Public,
	service.ContentAvailabilityTypePrivate: api.ContentAvailabilityType_Private,
}

var (
	ErrUnknownContentType             = errors.New("unknown content type")
	ErrUnknownContentAvailabilityType = errors.New("unknown content availability type")
)
