package repository_test

import (
	"github.com/golang/mock/gomock"
	"shortener/internal/storage/postgres/repository/mock"
	"testing"
)

func TestReturnLink(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create a mock instance of your repository
	mockRepo := mock.NewMockRepository(ctrl)
}
