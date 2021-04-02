package transport

import (
	"contentservice/pkg/contentservice/app/service"
	"contentservice/pkg/contentservice/infrastructure"
	"context"
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
	contentType, ok := apiToContentTypeMap[req.Type]
	if !ok {
		return nil, ErrUnknownContentType
	}

	err := server.container.ContentService().AddContent(req.Name, contentType)
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, err
}

var apiToContentTypeMap = map[api.ContentType]service.ContentType{
	api.ContentType_Song:    service.ContentTypeSong,
	api.ContentType_Podcast: service.ContentTypePodcast,
}

var (
	ErrUnknownContentType = errors.New("unknown content type")
)
