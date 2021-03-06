package repository

import (
	"github.com/stretchr/testify/assert"
	"gohanectl/hanectl/config"
	"gohanectl/hanectl/test/mock_test"
	"testing"
)

var devicesConfigDir = "../../config"
var devicesFile = "../../config/devices.yml"

func TestGetDevices(t *testing.T) {
	cfg := new(mock_test.ConfigurationRepoMock)
	cfg.On("GetStr", config.ConfigDirectory, config.DefConfigDirectory).Return(devicesConfigDir)
	cfg.On("GetStr", config.DeviceConfig, "").Return(devicesFile)
	cfg.On("GetStr", config.ScriptsTemplatesDirectory, config.DefScriptsTemplatesDirectory).Return(config.DefScriptsTemplatesDirectory)
	cfg.On("GetStr", config.ScriptsTemplatesDevicesDirectory, config.DefScriptsTemplatesDevicesDirectory).Return(config.DefScriptsTemplatesDevicesDirectory)
	cfg.On("GetStr", config.ScriptsDirectory, config.DefScriptsDirectory).Return(config.DefScriptsDirectory)
	cfg.On("GetInt", config.DeviceTimeout, config.DefDeviceTimeout).Return(config.DefDeviceTimeout)

	repo := NewDeviceRepo(cfg)
	devices, err := repo.GetDevices()
	assert.Nil(t, err)
	assert.True(t, len(devices.Devices) > 0)

	cfg.AssertExpectations(t)
}

func TestGetDevice(t *testing.T) {
	/*cfg := new(mock_test.ConfigurationRepoMock)
	cfg.On("GetStr", config.ConfigDirectory, config.DefConfigDirectory).Return(devicesConfigDir)
	cfg.On("GetStr", config.DeviceConfig, "").Return(devicesFile)
	cfg.On("GetStr", config.ScriptsTemplatesDirectory, config.DefScriptsTemplatesDirectory).Return(config.DefScriptsTemplatesDirectory)
	cfg.On("GetStr", config.ScriptsTemplatesDevicesDirectory, config.DefScriptsTemplatesDevicesDirectory).Return(config.DefScriptsTemplatesDevicesDirectory)
	cfg.On("GetStr", config.ScriptsDirectory, config.DefScriptsDirectory).Return(config.DefScriptsDirectory)
	cfg.On("GetInt", config.DeviceTimeout, config.DefDeviceTimeout).Return(config.DefDeviceTimeout)

	repo := NewDeviceRepo(cfg)
	devices, err := repo.GetDevices()
	assert.Nil(t, err)
	assert.True(t, len(devices.Devices) > 0)

	device, err2 := repo.GetDevice("licht-arbeitszimmer")
	assert.Nil(t, err2)
	assert.Equal(t, device.DeviceKey, "licht-arbeitszimmer")
	assert.Equal(t, device.Type, "light")
	assert.Equal(t, device.Caption, "Arbeitszimmer")
	assert.Equal(t, device.Optimistic, true)
	assert.Equal(t, device.Room, "Licht")
	assert.Equal(t, device.Template, "shelly_1")
	assert.Equal(t, device.Mqtt.DeviceID, "shelly1-xxxxxx")
	assert.Equal(t, device.Rest.Url, "http://192.168.9.52")

	device, err2 = repo.GetDevice("unknown")
	assert.Errorf(t, err2, "cant find device with key unknown")
	assert.Nil(t, device)

	cfg.AssertExpectations(t)
	*/
}
