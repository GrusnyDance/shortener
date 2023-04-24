package service_test

import (
	"context"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"shortener/internal/entities"
	"shortener/internal/service"
	"shortener/internal/storage/postgres/repository/mock"
	pb "shortener/proto/generate"
	"testing"
)

func TestPost_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock.NewMockRepository(ctrl)
	logger := logrus.New()
	service := service.New(mockRepo, logger)

	ctx := context.Background()
	request := &pb.PostRequest{
		LinkToHash: "https://example.com",
	}

	expectedHash := "b1Nx5zdaUy"
	mockRepo.EXPECT().CheckIfHashedExists(ctx, expectedHash).Return(errors.New("link not found"))
	mockRepo.EXPECT().CreateLink(ctx, expectedHash, request.LinkToHash).Return(nil)

	response, err := service.Post(ctx, request)

	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, expectedHash, response.HashedLink)
}

func TestPost_EmptyLink(t *testing.T) {
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

func TestPost_LinkExists(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock.NewMockRepository(ctrl)
	logger := logrus.New()
	service := service.New(mockRepo, logger)

	ctx := context.Background()
	request := &pb.PostRequest{
		LinkToHash: "https://example.com",
	}

	expectedHash := "abc123"
	mockRepo.EXPECT().CheckIfHashedExists(ctx, expectedHash).Return(nil)
	mockRepo.EXPECT().CreateLink(ctx, expectedHash, request.LinkToHash).Return(entities.ErrLinkExists)

	response, err := service.Post(ctx, request)

	assert.Error(t, err)
	assert.Nil(t, response)
	assert.True(t, errors.Is(err, entities.ErrLinkExists))
	assert.Equal(t, codes.AlreadyExists, status.Code(err))
}

func TestPost_StorageError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock.NewMockRepository(ctrl)
	logger := logrus.New()
	service := service.New(mockRepo, logger)

	ctx := context.Background()
	request := &pb.PostRequest{
		LinkToHash: "https://example.com",
	}

	expectedHash := "abc123"
	mockRepo.EXPECT().CheckIfHashedExists(ctx, expectedHash).Return(nil)
	mockRepo.EXPECT().CreateLink(ctx, expectedHash, request.LinkToHash).Return(entities.ErrStorage)

	response, err := service.Post(ctx, request)

	assert.Error(t, err)
	assert.Nil(t, response)
	assert.True(t, errors.Is(err, entities.ErrStorage))
	assert.Equal(t, codes.Unavailable, status.Code(err))
}
