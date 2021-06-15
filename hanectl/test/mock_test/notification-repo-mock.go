package mock_test

import (
	"github.com/stretchr/testify/mock"
	"gohanectl/hanectl/model"
)

type NotificationRepoMock struct {
	mock.Mock
}

func (m *NotificationRepoMock) GetAllNotifications() (*model.Notifications, error) {
	args := m.Called()
	return args.Get(0).(*model.Notifications), args.Error(1)
}

func (m *NotificationRepoMock) GetNotifications(deviceKey string, key string) ([]*model.Notification, error) {
	args := m.Called(deviceKey, key)
	return args.Get(0).([]*model.Notification), args.Error(1)
}

func (m *NotificationRepoMock) ReloadNotifications() error {
	args := m.Called()
	return args.Error(0)
}
