package mock_test

import (
	"github.com/stretchr/testify/mock"
	"gohanectl/hanectl/model"
)

type DeviceServiceMock struct {
	mock.Mock
}

func (m *DeviceServiceMock) DeviceCommand(deviceKey string, params model.Dictionary) bool {
	args := m.Called(deviceKey, params)
	return args.Bool(0)
}
func (m *DeviceServiceMock) DeviceStates() model.Dictionary {
	args := m.Called()
	return args.Get(0).(model.Dictionary)

}
func (m *DeviceServiceMock) DeviceState(deviceKey string) interface{} {
	args := m.Called(deviceKey)
	return args.Get(0)
}
func (m *DeviceServiceMock) GetDevice(deviceKey string) (*model.Device, error) {
	args := m.Called(deviceKey)
	return args.Get(0).(*model.Device), args.Error(1)
}
func (m *DeviceServiceMock) GetDevices() (*model.Devices, error) {
	args := m.Called()
	return args.Get(0).(*model.Devices), args.Error(1)
}
func (m *DeviceServiceMock) GetDevicesDto(filter model.DeviceFilter) (*model.DevicesDto, error) {
	args := m.Called(filter)
	return args.Get(0).(*model.DevicesDto), args.Error(1)
}
func (m *DeviceServiceMock) DeviceUpdate(dev *model.Device, oldState string) {
	m.Called(dev, oldState)
}
func (m *DeviceServiceMock) ReloadDevices() error {
	args := m.Called()
	return args.Error(0)
}
