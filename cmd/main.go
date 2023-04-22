package main

import (
	"context"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"net"
	"net/http"
	"shortener/internal/service"
	pb "shortener/proto/generate"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

func main() {
	// Create a listener on TCP port
	lis, err := net.Listen("tcp", "localhost:8080")
	if err != nil {
		log.Fatalln("Failed to listen:", err)
	}
	defer lis.Close() // надо ли?

	// Create a gRPC service object
	s := grpc.NewServer()
	// Attach the Shortener service to the service
	pb.RegisterShortenerServer(s, service.New())
	// Serve gRPC service
	log.Println("Serving gRPC on connection ")
	go func() {
		log.Fatalln(s.Serve(lis))
	}()

	// Client connection is used by the gRPC-Gateway to forward
	// incoming HTTP/REST requests to the gRPC service for processing
	conn, err := grpc.Dial("localhost:8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
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
		Addr:    "localhost:8085",
		Handler: mux,
	}

	log.Println("Serving gRPC-Gateway on connection")
	log.Fatalln(gwServer.ListenAndServe())
}
