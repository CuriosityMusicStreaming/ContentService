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

	DatabaseUser     string `envconfig:"db_user" default:"root"`
	DatabasePassword string `envconfig:"db_password" default:"1234"`
	DatabaseHost     string `envconfig:"db_host" default:"contentservice-db"`
	DatabaseName     string `envconfig:"db_name" default:"contentservice"`

	MaxDatabaseConnections int `envconfig:"max_connections" default:"10"`
}
