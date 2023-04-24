package service

import (
	"context"
	"errors"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

	if request.LinkToHash == "" {
		return nil, s.ProcessErr(entities.ErrEmptyLink, codes.InvalidArgument)
	}
	hash := hasher.Apply(request.LinkToHash)
	err := s.repo.CheckIfHashedExists(ctx, hash)
	if err == nil {
		return nil, s.ProcessErr(entities.ErrLinkExists, codes.AlreadyExists)
	}
	err = s.repo.CreateLink(ctx, hash, request.LinkToHash)
	if err != nil {
		return nil, s.ProcessErr(entities.ErrStorage, codes.Unavailable)
	}
	response := &pb.PostResponse{
		HashedLink: hash,
	}
	return response, nil
}

func (s *Service) Get(ctx context.Context, request *pb.GetRequest) (*pb.GetResponse, error) {
	s.logger.Infoln("get, hashed link is", request.HashedLink)

	if request.HashedLink == "" {
		return nil, s.ProcessErr(entities.ErrEmptyLink, codes.InvalidArgument)
	}
	hash := request.HashedLink
	err := s.repo.CheckIfHashedExists(ctx, hash)
	if err != nil && errors.Is(err, entities.ErrNotFound) {
		return nil, s.ProcessErr(entities.ErrNotFound, codes.NotFound)
	} else if err != nil {
		return nil, s.ProcessErr(entities.ErrStorage, codes.Unavailable)
	}
	link, er := s.repo.ReturnLink(ctx, hash)
	if er != nil {
		return nil, s.ProcessErr(entities.ErrStorage, codes.Unavailable)
	}
	response := &pb.GetResponse{
		OriginalLink: link,
	}
	return response, nil
}

func (s *Service) ProcessErr(er error, code codes.Code) error {
	st := status.New(code, er.Error())
	s.logger.Error(entities.ErrLinkExists)
	return st.Err()
}
