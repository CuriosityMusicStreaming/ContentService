package transport

import (
	"contentservice/pkg/common/infrastructure/logging/activity"
	"context"
	"fmt"
	"github.com/pkg/errors"
	logger "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"net/http"
	"strings"
	"time"
)

var (
	ErrFailedToReadFromIncomingContext = errors.New("failed to read metadata from context")
)

func NewLoggingMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		now := time.Now()
		h.ServeHTTP(writer, request)

		logger.WithFields(logger.Fields{
			"duration": time.Since(now),
			"method":   request.Method,
			"url":      request.RequestURI,
		}).Info("request finished")
	})
}

func NewLoggerServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			logger.WithFields(logger.Fields{
				"args":   req,
				"method": getGRPCMethodName(info),
			}).Error("failed to read metadata from request. Aborting request")
			return nil, ErrFailedToReadFromIncomingContext
		}

		activityID, err := activity.ParseActivityID(md.Get("activityID")[0])
		if err != nil {
			logger.WithFields(logger.Fields{
				"args":   req,
				"method": getGRPCMethodName(info),
			}).Error("failed to parse activity id from request. Aborting request")
			return nil, ErrFailedToReadFromIncomingContext
		}

		start := time.Now()

		resp, err = handler(ctx, req)

		fields := logger.Fields{
			"activityID": activityID.String(),
			"args":       req,
			"duration":   fmt.Sprintf("%v", time.Since(start)),
			"method":     getGRPCMethodName(info),
		}

		entry := logger.WithFields(fields)
		if err != nil {
			entry.Error("call failed")
		} else {
			entry.Info("call finished")
		}

		return resp, err
	}
}

func getGRPCMethodName(info *grpc.UnaryServerInfo) string {
	method := info.FullMethod
	return method[strings.LastIndex(method, "/")+1:]
}
