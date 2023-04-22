package main

import (
	"context"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"net"
	"net/http"
	pb "shortener/proto/generate"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedGreeterServer
}

func (s *server) SayHello(ctx context.Context, request *pb.HelloRequest) (*pb.HelloResponse, error) {
	name := request.Name
	response := &pb.HelloResponse{
		Message: "Hello " + name,
	}
	return response, nil
}

func main() {
	// Create a listener on TCP port
	lis, err := net.Listen("tcp", "localhost:8080")
	if err != nil {
		log.Fatalln("Failed to listen:", err)
	}

	// Create a gRPC server object
	s := grpc.NewServer()
	// Attach the Greeter service to the server
	pb.RegisterGreeterServer(s, &server{})
	// Serve gRPC server
	log.Println("Serving gRPC on connection ")
	go func() {
		log.Fatalln(s.Serve(lis))
	}()

	// Client connection is used by the gRPC-Gateway to forward
	// incoming HTTP/REST requests to the gRPC server for processing
	conn, err := grpc.Dial("localhost:8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalln("Failed to dial server:", err)
	}
	defer conn.Close()

	mux := runtime.NewServeMux()
	// Register Greeter
	err = pb.RegisterGreeterHandler(context.Background(), mux, conn)
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
