package transport

import (
	"contentservice/pkg/contentservice/domain"
	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func translateError(err error) error {
	switch errors.Cause(err) {
	case domain.ErrContentNotFound:
		return status.Error(codes.NotFound, err.Error())
	case domain.ErrOnlyAuthorCanDeleteContent:
	case domain.ErrOnlyAuthorCanManageContent:
		return status.Error(codes.PermissionDenied, err.Error())
	case ErrUnknownContentType:
	case ErrUnknownContentAvailabilityType:
		return status.Error(codes.InvalidArgument, err.Error())
	}

	return err
}
