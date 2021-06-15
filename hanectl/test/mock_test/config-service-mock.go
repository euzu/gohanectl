package mock_test

import (
	"github.com/stretchr/testify/mock"
	"gohanectl/hanectl/model"
)

type ConfigServiceMock struct {
	mock.Mock
}

func (m *ConfigServiceMock) GetServerStatus() (*model.ServerStatus, error) {
	args := m.Called()
	return args.Get(0).(*model.ServerStatus), args.Error(1)
}
func (m *ConfigServiceMock) GetRoomsConfig() (model.Dictionary, error) {
	args := m.Called()
	return args.Get(0).(model.Dictionary), args.Error(1)
}
func (m *ConfigServiceMock) SetWebsocketStatus(connected bool) {
	m.Called(connected)
}
func (m *ConfigServiceMock) SetMqttStatus(connectionStatus int) {
	m.Called(connectionStatus)
}