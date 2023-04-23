package service

import (
	"context"
	"errors"
	"fmt"
	"shortener/internal/entities"
	"shortener/pkg/hasher"
	pb "shortener/proto/generate"
)

// Service implements gRPC server
type Service struct {
	pb.UnimplementedShortenerServer
	repo entities.Repository
}

// New is a constructor for Service
func New(r entities.Repository) *Service {
	fmt.Println("starting service")
	return &Service{repo: r}
}

func (s *Service) Post(ctx context.Context, request *pb.PostRequest) (*pb.PostResponse, error) {
	hash := hasher.Apply(request.LinkToHash)
	err := s.repo.CheckIfHashedExists(ctx, hash)
	if err == nil {
		return nil, errors.New("link exists already")
	}
	err = s.repo.CreateLink(ctx, hash, request.LinkToHash)
	if err != nil {
		return nil, errors.New("link cannot be saved")
	}
	response := &pb.PostResponse{
		HashedLink: hash,
	}
	return response, nil
}

func (s *Service) Get(ctx context.Context, request *pb.GetRequest) (*pb.GetResponse, error) {
	hash := request.HashedLink
	fmt.Println("hash is", hash)
	err := s.repo.CheckIfHashedExists(ctx, hash)
	if err != nil && errors.Is(err, errors.New("link not found")) {
		return nil, errors.New("link does not exist")
	} else if err != nil {
		return nil, err
	}
	link, er := s.repo.ReturnLink(ctx, hash)
	if er != nil {
		return nil, errors.New("original link cannot be received")
	}
	response := &pb.GetResponse{
		OriginalLink: link,
	}
	return response, nil
}
