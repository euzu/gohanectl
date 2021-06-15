package mock_test

import (
	"github.com/stretchr/testify/mock"
	"gohanectl/hanectl/model"
)

type NotificationServiceMock struct {
	mock.Mock
}

func (m *NotificationServiceMock) GetNotifications(deviceKey string, key string) ([]*model.Notification, error) {
	args := m.Called(deviceKey, key)
	return args.Get(0).([]*model.Notification), args.Error(1)

}
func (m *NotificationServiceMock) ReloadNotifications() error {
	args := m.Called()
	return args.Error(0)

}
