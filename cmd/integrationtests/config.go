package main

import (
	"github.com/kelseyhightower/envconfig"
	"github.com/pkg/errors"
)

func parseConfig() (*config, error) {
	c := new(config)
	if err := envconfig.Process(appID, c); err != nil {
		return nil, errors.Wrap(err, "failed to parse env")
	}
	return c, nil
}

type config struct {
	ContentServiceHost     string `envconfig:"contentservice_host"`
	ContentServiceRESTPort string `envconfig:"contentservice_rest_port"`
	ContentServiceGRPCPort string `envconfig:"contentservice_grpc_port"`
	MaxWaitTimeSeconds     int    `envconfig:"max_wait_time_seconds"`

	ServeGRPCAddress string `envconfig:"serve_grpc_address" default:":8002"`
}
