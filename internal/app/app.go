package app

import (
	"context"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	l "log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"shortener/internal/entities"
	"shortener/internal/service"
	"shortener/internal/storage/cache"
	"shortener/internal/storage/postgres/repository"
	pb "shortener/proto/generate"
	"syscall"
)

func Start() {
	// Setup logging
	f, err := os.OpenFile("log.txt", os.O_WRONLY|os.O_CREATE, 0755)
	if err != nil {
		l.Fatal(err)
	}
	defer f.Close()
	logger := InitLogger(f)

	// Create a listener on TCP port
	lis, err := net.Listen("tcp", os.Getenv("GRPC_PORT"))
	if err != nil {
		logger.Fatalln("Failed to listen:", err)
	}

	// Create a gRPC service object
	s := grpc.NewServer()

	// Init storage
	instance, er := InitStorage()
	if er != nil {
		logger.Fatalln(er)
	}
	defer instance.Close()

	// Attach the Shortener service to the service
	pb.RegisterShortenerServer(s, service.New(instance, logger))
	logger.Infoln("Serving gRPC on connection ")
	go func() {
		logger.Fatalln(s.Serve(lis))
	}()

	// Client connection is used by the gRPC-Gateway to forward
	// incoming HTTP/REST requests to the gRPC service
	conn, err := grpc.Dial(os.Getenv("GRPC_PORT"), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Fatalln("Failed to dial service:", err)
	}
	defer conn.Close()

	mux := runtime.NewServeMux()
	// Register Shortener
	ctx, cancel := context.WithCancel(context.Background())
	err = pb.RegisterShortenerHandler(ctx, mux, conn)
	if err != nil {
		logger.Fatalln("Failed to register gateway:", err)
	}

	gwServer := &http.Server{
		Addr:    os.Getenv("HTTP_PROXY_PORT"),
		Handler: mux,
	}

	// Graceful
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	go func() {
		select {
		case <-ctx.Done():
			return
		case <-sigCh:
			logger.Infoln("Received shutdown signal, initiating graceful shutdown")
			s.GracefulStop()
			if err = gwServer.Shutdown(ctx); err != nil {
				logger.Fatalln("error finishing http", err)
			}
			cancel()
		}
	}()

	logger.Infoln("Serving gRPC-Gateway on connection")
	logger.Fatalln(gwServer.ListenAndServe())
}

func InitLogger(f *os.File) *logrus.Logger {
	logger := &logrus.Logger{
		Out:       f,
		Formatter: new(logrus.TextFormatter),
		Hooks:     make(logrus.LevelHooks),
		Level:     logrus.InfoLevel,
	}
	return logger
}

func InitStorage() (entities.Repository, error) {
	var instance entities.Repository
	var err error
	if os.Getenv("ENABLE_DB") == "true" {
		instance, err = repository.Init()
		if err != nil {
			return nil, err
		}
	} else {
		instance = cache.Init()
	}
	return instance, err
}
