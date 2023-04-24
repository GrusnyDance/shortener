package service_test

import (
	"context"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"shortener/internal/service"
	"shortener/internal/storage/postgres/repository/mock"
	pb "shortener/proto/generate"
	"testing"
)

func BasicPost(mockRepo *mock.MockRepository, service *service.Service) (*pb.PostResponse, string, error) {
	ctx := context.Background()
	request := &pb.PostRequest{
		LinkToHash: "https://example.com",
	}

	expectedHash := "b1Nx5zdaUy"
	mockRepo.EXPECT().CheckIfHashedExists(ctx, expectedHash).Return(errors.New("link not found"))
	mockRepo.EXPECT().CreateLink(ctx, expectedHash, request.LinkToHash).Return(nil)
	response, err := service.Post(ctx, request)
	return response, expectedHash, err
}

func TestPostSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock.NewMockRepository(ctrl)
	logger := logrus.New()
	service := service.New(mockRepo, logger)

	response, expectedHash, err := BasicPost(mockRepo, service)

	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, expectedHash, response.HashedLink)
}

func TestPostEmptyLink(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock.NewMockRepository(ctrl)
	logger := logrus.New()
	service := service.New(mockRepo, logger)

	ctx := context.Background()
	request := &pb.PostRequest{
		LinkToHash: "",
	}
	response, err := service.Post(ctx, request)

	assert.Error(t, err)
	assert.Nil(t, response)
	assert.Equal(t, codes.InvalidArgument, status.Code(err))
}

func TestPostLinkExists(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock.NewMockRepository(ctrl)
	logger := logrus.New()
	service := service.New(mockRepo, logger)

	ctx := context.Background()
	request := &pb.PostRequest{
		LinkToHash: "https://example.com",
	}

	response, expectedHash, err := BasicPost(mockRepo, service)
	mockRepo.EXPECT().CheckIfHashedExists(ctx, expectedHash).Return(nil)
	response, err = service.Post(ctx, request)

	assert.Error(t, err)
	assert.Nil(t, response)
	assert.Equal(t, codes.AlreadyExists, status.Code(err))
}

func TestGetSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock.NewMockRepository(ctrl)
	logger := logrus.New()
	service := service.New(mockRepo, logger)

	ctx := context.Background()
	_, expectedHash, _ := BasicPost(mockRepo, service)

	requestGet := &pb.GetRequest{
		HashedLink: "b1Nx5zdaUy",
	}

	expectedLink := "https://example.com"
	mockRepo.EXPECT().CheckIfHashedExists(ctx, expectedHash).Return(nil)
	mockRepo.EXPECT().ReturnLink(ctx, requestGet.HashedLink).Return(expectedLink, nil)

	response, err := service.Get(ctx, requestGet)

	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, expectedLink, response.OriginalLink)
}

func TestGetLinkNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock.NewMockRepository(ctrl)
	logger := logrus.New()
	service := service.New(mockRepo, logger)

	ctx := context.Background()
	request := &pb.GetRequest{
		HashedLink: "sfgsdfg",
	}

	mockRepo.EXPECT().CheckIfHashedExists(ctx, request.HashedLink).Return(errors.New("link not found"))
	response, err := service.Get(ctx, request)

	assert.Error(t, err)
	assert.Nil(t, response)
}

func TestGetEmptyLink(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock.NewMockRepository(ctrl)
	logger := logrus.New()
	service := service.New(mockRepo, logger)

	ctx := context.Background()
	request := &pb.GetRequest{
		HashedLink: "",
	}

	response, err := service.Get(ctx, request)

	assert.Error(t, err)
	assert.Nil(t, response)
	assert.Equal(t, codes.InvalidArgument, status.Code(err))
}
