package transport

import (
	"context"
	"fmt"
	logger "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"net/http"
	"strings"
	"time"
)

func NewLoggingMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		now := time.Now()
		h.ServeHTTP(writer, request)

		logger.WithFields(logger.Fields{
			"duration":  time.Since(now),
			"method":    request.Method,
			"url":       request.URL,
			"userAgent": request.UserAgent(),
		}).Info("request finished")
	})
}

func NewLoggerServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		start := time.Now()

		resp, err = handler(ctx, req)

		duration := fmt.Sprintf("%v", time.Since(start))

		fields := logger.Fields{
			"args":     req,
			"duration": duration,
			"method":   getGRPCMethodName(info),
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
