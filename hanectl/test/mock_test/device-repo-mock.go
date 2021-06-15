package mock_test

import (
	"github.com/stretchr/testify/mock"
	"gohanectl/hanectl/model"
)

type DeviceRepoMock struct {
	mock.Mock
}

func (m *DeviceRepoMock) ReloadDevices() error {
	args := m.Called()
	return args.Error(0)
}

func (m *DeviceRepoMock) GetDevice(deviceKey string) (*model.Device, error) {
	args := m.Called(deviceKey)
	return args.Get(0).(*model.Device), args.Error(1)
}

func (m *DeviceRepoMock) GetDevices() (*model.Devices, error) {
	args := m.Called()
	return args.Get(0).(*model.Devices), args.Error(1)
}