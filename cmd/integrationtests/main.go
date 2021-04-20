package main

import (
	contentserviceapi "contentservice/api/contentservice"
	userserviceapi "contentservice/api/userservice"
	"contentservice/pkg/intergrationtests/app"
	"fmt"
	log "github.com/CuriosityMusicStreaming/ComponentsPool/pkg/app/logger"
	jsonlog "github.com/CuriosityMusicStreaming/ComponentsPool/pkg/infrastructure/logger"
	"google.golang.org/grpc"
)

var appID = "UNKNOWN"

func main() {
	logger := initLogger()

	config, err := parseConfig()
	if err != nil {
		logger.FatalError(err)
	}

	err = runService(config)
	if err != nil {
		logger.FatalError(err)
	}
}

func runService(config *config) error {
	opts := []grpc.DialOption{
		grpc.WithInsecure(),
	}

	contentServiceClient, err := initContentServiceClient(opts, config)
	if err != nil {
		return err
	}
	userServiceClient, err := initUserServiceClient(opts, config)
	if err != nil {
		return err
	}

	app.RunTests(
		contentServiceClient,
		userServiceClient,
	)

	return nil
}

func initLogger() log.MainLogger {
	return jsonlog.NewLogger(&jsonlog.Config{AppName: appID})
}

func initContentServiceClient(commonOpts []grpc.DialOption, config *config) (contentserviceapi.ContentServiceClient, error) {
	conn, err := grpc.Dial(fmt.Sprintf("%s:%s", config.ContentServiceHost, config.ContentServiceGRPCPort), commonOpts...)
	if err != nil {
		return nil, err
	}

	return contentserviceapi.NewContentServiceClient(conn), nil
}

func initUserServiceClient(commonOpts []grpc.DialOption, config *config) (userserviceapi.UserServiceClient, error) {
	conn, err := grpc.Dial(fmt.Sprintf("%s:%s", config.UserServiceHost, config.UserServiceGRPCPort), commonOpts...)
	if err != nil {
		return nil, err
	}

	return userserviceapi.NewUserServiceClient(conn), nil
}
