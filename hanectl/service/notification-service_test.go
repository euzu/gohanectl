package service

import (
	"github.com/stretchr/testify/assert"
	"gohanectl/hanectl/model"
	"gohanectl/hanectl/test/mock_test"
	"testing"
)

func TestGetNotifications(t *testing.T) {
	repoMock := new(mock_test.NotificationRepoMock)
	srv := NotificationService{
		notificationRepo: repoMock,
	}
	repoMock.On("GetNotifications", "deviceKey", "key").Return([]*model.Notification{}, nil)

	_, err := srv.GetNotifications("deviceKey", "key")
	assert.Nil(t, err)

	repoMock.AssertExpectations(t)
}
