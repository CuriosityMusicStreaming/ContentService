package main

import (
	"github.com/kelseyhightower/envconfig"
	"github.com/pkg/errors"
)

func parseEnv() (*config, error) {
	c := new(config)
	if err := envconfig.Process(appID, c); err != nil {
		return nil, errors.Wrap(err, "failed to parse env")
	}
	return c, nil
}

type config struct {
	ServeRESTAddress string `envconfig:"serve_rest_address" default:":8001"`
	ServeGRPCAddress string `envconfig:"serve_grpc_address" default:":8002"`
	DatabaseDriver   string `envconfig:"db_driver" default:"mysql"`
	DSN              string `envconfig:"dsn" default:"root:1234@/orderservice"`
}
