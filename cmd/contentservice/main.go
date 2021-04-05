package main

import (
	"contentservice/api/contentservice"
	"contentservice/pkg/common/infrastructure/mysql"
	"contentservice/pkg/common/infrastructure/server"
	"contentservice/pkg/contentservice/infrastructure"
	"contentservice/pkg/contentservice/infrastructure/transport"
	"context"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"
	"io"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	logger "github.com/sirupsen/logrus"
)

var appID = "UNKNOWN"

func main() {
	logger.SetFormatter(&logger.JSONFormatter{})

	config, err := parseEnv()
	if err != nil {
		logger.Fatal(err)
	}

	err = runService(config)
	if err == server.ErrStopped {
		logger.Info("service is successfully stopped")
	} else if err != nil {
		logger.Fatal(err)
	}
}

func runService(config *config) error {
	dsn := mysql.DSN{
		User:     config.DatabaseUser,
		Password: config.DatabasePassword,
		Host:     config.DatabaseHost,
		Database: config.DatabaseName,
	}
	connector := mysql.NewConnector()

	err := connector.Open(dsn, config.MaxDatabaseConnections)
	if err != nil {
		return err
	}

	defer connector.Close()

	stopChan := make(chan struct{})
	listenForKillSignal(stopChan)

	container := infrastructure.NewDependencyContainer()

	serviceApi := transport.NewContentServiceServer(container)
	serverHub := server.NewHub(stopChan)

	baseServer := grpc.NewServer(grpc.UnaryInterceptor(makeGRPCUnaryInterceptor()))
	contentservice.RegisterContentServiceServer(baseServer, serviceApi)

	serverHub.AddServer(server.NewGrpcServer(
		baseServer,
		server.GrpcServerConfig{
			ServeAddress: config.ServeGRPCAddress,
		}),
	)

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	var httpServer *http.Server

	serverHub.AddServer(&server.FuncServer{
		ServeImpl: func() error {
			grpcGatewayMux := runtime.NewServeMux()
			opts := []grpc.DialOption{grpc.WithInsecure()}
			err := contentservice.RegisterContentServiceHandlerFromEndpoint(ctx, grpcGatewayMux, config.ServeGRPCAddress, opts)
			if err != nil {
				return err
			}

			router := mux.NewRouter()
			router.PathPrefix("/api/").Handler(grpcGatewayMux)

			// Implement healthcheck for Kubernetes
			router.HandleFunc("/resilience/ready", func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(http.StatusOK)
				_, _ = io.WriteString(w, http.StatusText(http.StatusOK))
			}).Methods(http.MethodGet)

			httpServer = &http.Server{
				Handler:      transport.NewLoggingMiddleware(router),
				Addr:         config.ServeRESTAddress,
				WriteTimeout: 15 * time.Second,
				ReadTimeout:  15 * time.Second,
			}

			logger.Info("REST server started")
			return httpServer.ListenAndServe()
		},
		StopImpl: func() error {
			cancel()
			return httpServer.Shutdown(context.Background())
		},
	})

	return serverHub.Run()
}

func listenForKillSignal(stopChan chan<- struct{}) {
	go func() {
		ch := make(chan os.Signal, 1)
		signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT)
		<-ch
		stopChan <- struct{}{}
	}()
}

func makeGRPCUnaryInterceptor() grpc.UnaryServerInterceptor {
	loggerInterceptor := transport.NewLoggerServerInterceptor()
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		resp, err = loggerInterceptor(ctx, req, info, handler)
		return resp, err
	}
}
