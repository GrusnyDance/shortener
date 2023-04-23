package main

import (
	"context"
	"flag"
	"github.com/joho/godotenv"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/grpclog"
	"log"
	"net"
	"net/http"
	"os"
	"shortener/internal/service"
	"shortener/internal/storage/postgres/repository"
	pb "shortener/proto/generate"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		grpclog.Fatal(err)
	}
	flag.Parse()

	// Create a listener on TCP port
	lis, err := net.Listen("tcp", os.Getenv("GRPC_PORT"))
	if err != nil {
		log.Fatalln("Failed to listen:", err)
	}
	defer lis.Close() // надо ли?

	// Create a gRPC service object
	s := grpc.NewServer()
	instance, err := repository.Init()
	if err != nil {
		grpclog.Fatal(err)
	}
	// Attach the Shortener service to the service
	pb.RegisterShortenerServer(s, service.New(instance))
	// Serve gRPC service
	log.Println("Serving gRPC on connection ")
	go func() {
		log.Fatalln(s.Serve(lis))
	}()

	// Client connection is used by the gRPC-Gateway to forward
	// incoming HTTP/REST requests to the gRPC service for processing
	conn, err := grpc.Dial(os.Getenv("GRPC_PORT"), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalln("Failed to dial service:", err)
	}
	defer conn.Close()

	mux := runtime.NewServeMux()
	// Register Shortener
	err = pb.RegisterShortenerHandler(context.Background(), mux, conn)
	if err != nil {
		log.Fatalln("Failed to register gateway:", err)
	}

	gwServer := &http.Server{
		Addr:    os.Getenv("HTTP_PROXY_PORT"),
		Handler: mux,
	}

	log.Println("Serving gRPC-Gateway on connection")
	log.Fatalln(gwServer.ListenAndServe())
}
