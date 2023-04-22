package service

import (
	"context"
	"errors"
	pb "shortener/proto/generate"
)

// Service implements gRPC server
type Service struct {
	pb.UnimplementedShortenerServer
	//repo Rep
}

// New is a constructor for Service
func New() *Service {
	return &Service{}
}

func (s *Service) Post(ctx context.Context, request *pb.PostRequest) (*pb.PostResponse, error) {
	//l := request.LinkToHash
	response := &pb.PostResponse{
		HashedLink: "Hello baby",
	}
	return response, nil
}

func (s *Service) Get(ctx context.Context, request *pb.GetRequest) (*pb.GetResponse, error) {
	//l := request.LinkToHash
	response := &pb.GetResponse{
		OriginalLink: "Byebye baby",
	}
	return response, errors.New("not found")
}
