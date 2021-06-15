package mock_test

import (
	"github.com/stretchr/testify/mock"
	"gohanectl/hanectl/model"
)

type ServiceFactoryMock struct {
	mock.Mock
}

func (m *ServiceFactoryMock) GetDeviceService() model.IDeviceService {
	args := m.Called()
	return args.Get(0).(model.IDeviceService)
}
func (m *ServiceFactoryMock) GetConfigService() model.IConfigService {
	args := m.Called()
	return args.Get(0).(model.IConfigService)
}
func (m *ServiceFactoryMock) GetMqttService() model.IMqttService {
	args := m.Called()
	return args.Get(0).(model.IMqttService)
}
func (m *ServiceFactoryMock) GetTelegramService() model.ITelegramService {
	args := m.Called()
	return args.Get(0).(model.ITelegramService)
}
func (m *ServiceFactoryMock) GetUserService() model.IUserService {
	args := m.Called()
	return args.Get(0).(model.IUserService)
}
func (m *ServiceFactoryMock) GetNotificationService() model.INotificationService {
	args := m.Called()
	return args.Get(0).(model.INotificationService)
}
func (m *ServiceFactoryMock) GetRestService() model.IRestService {
	args := m.Called()
	return args.Get(0).(model.IRestService)
}
func (m *ServiceFactoryMock) GetSharedMemory() model.ISharedMemory {
	args := m.Called()
	return args.Get(0).(model.ISharedMemory)
}
func (m *ServiceFactoryMock) Finalize() {
	m.Called()
}
