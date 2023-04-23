package service

import (
	"context"
	"errors"
	"github.com/sirupsen/logrus"
	"shortener/internal/entities"
	"shortener/pkg/hasher"
	pb "shortener/proto/generate"
)

// Service implements gRPC server
type Service struct {
	pb.UnimplementedShortenerServer
	repo   entities.Repository
	logger *logrus.Logger
}

// New is a constructor for Service
func New(r entities.Repository, log *logrus.Logger) *Service {
	return &Service{repo: r, logger: log}
}

func (s *Service) Post(ctx context.Context, request *pb.PostRequest) (*pb.PostResponse, error) {
	s.logger.Infoln("post, original link is", request.LinkToHash)

	hash := hasher.Apply(request.LinkToHash)
	err := s.repo.CheckIfHashedExists(ctx, hash)
	if err == nil {
		s.logger.Error(entities.ErrLinkExists)
		return nil, entities.ErrLinkExists
	}
	err = s.repo.CreateLink(ctx, hash, request.LinkToHash)
	if err != nil {
		s.logger.Error(entities.ErrStorage)
		return nil, entities.ErrStorage
	}
	response := &pb.PostResponse{
		HashedLink: hash,
	}
	return response, nil
}

func (s *Service) Get(ctx context.Context, request *pb.GetRequest) (*pb.GetResponse, error) {
	s.logger.Infoln("get, hashed link is", request.HashedLink)

	hash := request.HashedLink
	err := s.repo.CheckIfHashedExists(ctx, hash)
	if err != nil && errors.Is(err, entities.ErrNotFound) {
		s.logger.Error(entities.ErrNotFound)
		return nil, entities.ErrNotFound
	} else if err != nil {
		s.logger.Error(entities.ErrStorage)
		return nil, entities.ErrStorage
	}
	link, er := s.repo.ReturnLink(ctx, hash)
	if er != nil {
		s.logger.Error(entities.ErrStorage)
		return nil, entities.ErrStorage
	}
	response := &pb.GetResponse{
		OriginalLink: link,
	}
	return response, nil
}
