package service

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"gohanectl/hanectl/model"
	"gohanectl/hanectl/test/mock_test"
	"testing"
)

func TestDeviceCommandWithoutParamCommand(t *testing.T) {
	deviceKey := "test-device"
	repoMock := new(mock_test.DeviceRepoMock)
	srv := DeviceService{
		deviceRepo: repoMock,
	}
	repoMock.On("GetDevice", deviceKey).Return(&model.Device{DeviceKey: deviceKey}, nil)

	var params = model.Dictionary{}
	success := srv.DeviceCommand(deviceKey, params)
	assert.False(t, success)

	repoMock.AssertExpectations(t)
}

func TestDeviceCommandWithoutCommandTopic(t *testing.T) {
	deviceKey := "test-device"
	device := model.Device{
		DeviceKey: deviceKey,
		Mqtt: model.DeviceMqtt{
			CommandTopics: model.Dictionary{},
		},
	}
	repoMock := new(mock_test.DeviceRepoMock)
	srv := DeviceService{
		deviceRepo: repoMock,
	}
	repoMock.On("GetDevice", deviceKey).Return(&device, nil)

	var params = model.Dictionary{
		"command": "power",
	}
	success := srv.DeviceCommand(deviceKey, params)
	assert.False(t, success)

	repoMock.AssertExpectations(t)
}

func TestDeviceCommandMqttPowerOnOff(t *testing.T) {
	deviceKey := "test-device"
	topic := "cmnd/power"
	device := model.Device{
		DeviceKey: deviceKey,
		Mqtt: model.DeviceMqtt{
			CommandTopics: model.Dictionary{
				"power": model.Dictionary{
					"topic":       topic,
					"payload_on":  "on",
					"payload_off": "off",
				},
			},
		},
	}
	repoMock := new(mock_test.DeviceRepoMock)
	mqttServiceMock := new(mock_test.MqttServiceMock)
	srv := DeviceService{
		deviceRepo: repoMock,
		mqttService: mqttServiceMock,
	}
	repoMock.On("GetDevice", deviceKey).Return(&device, nil)
	mqttServiceMock.On("Publish", topic, "on").Times(1).Return(true)
	mqttServiceMock.On("Publish", topic, "off").Times(1).Return(true)

	success := srv.DeviceCommand(deviceKey, model.Dictionary{
		"command": "power",
		"payload": true,
	})
	assert.True(t, success)

	success = srv.DeviceCommand(deviceKey, model.Dictionary{
		"command": "power",
		"payload": false,
	})
	assert.True(t, success)

	//success = srv.DeviceCommand(deviceKey, model.Dictionary{
	//	"command": "power",
	//	"payload": 1,
	//})
	//assert.True(t, success)

	repoMock.AssertExpectations(t)
	mqttServiceMock.AssertExpectations(t)
}

func TestDeviceCommandMqttProperty(t *testing.T) {
	deviceKey := "test-device"
	topic := "cmnd/status"
	device := model.Device{
		DeviceKey: deviceKey,
		Mqtt: model.DeviceMqtt{
			CommandTopics: model.Dictionary{
				"status": model.Dictionary{
					"topic":   topic,
					"payload": "",
				},
			},
		},
	}
	repoMock := new(mock_test.DeviceRepoMock)
	mqttServiceMock := new(mock_test.MqttServiceMock)
	srv := DeviceService{
		deviceRepo: repoMock,
		mqttService: mqttServiceMock,
	}
	repoMock.On("GetDevice", deviceKey).Return(&device, nil)
	mqttServiceMock.On("Publish", topic, "").Times(1).Return(true)

	success := srv.DeviceCommand(deviceKey, model.Dictionary{
		"command": "status",
		"payload": "",
	})
	assert.True(t, success)

	success = srv.DeviceCommand(deviceKey, model.Dictionary{
		"command": "unknown",
		"payload": "",
	})
	assert.False(t, success)

	repoMock.AssertExpectations(t)
	mqttServiceMock.AssertExpectations(t)
}

func TestDeviceCommandRestPowerOnOff(t *testing.T) {
	deviceKey := "test-device"
	url := "http://192.168.10.10"
	powerOn := "/cm?cmnd=Power On"
	powerOff := "/cm?cmnd=Power Off"
	device := model.Device{
		DeviceKey: deviceKey,
		Rest: model.DeviceRest{
			Url: url,
			CommandPaths: model.Dictionary{
				"power": model.Dictionary{
					"power_on":  powerOn,
					"power_off": powerOff,
				},
			},
		},
	}
	repoMock := new(mock_test.DeviceRepoMock)
	restServiceMock := new(mock_test.RestServiceMock)
	srv := DeviceService{
		deviceRepo: repoMock,
		restService: restServiceMock,
	}

	repoMock.On("GetDevice", deviceKey).Return(&device, nil)
	restServiceMock.On("GetRequest", url+powerOn, &device).Times(1).Return(true)
	restServiceMock.On("GetRequest", url+powerOff, &device).Times(1).Return(true)

	success := srv.DeviceCommand(deviceKey, model.Dictionary{
		"command": "power",
		"payload": true,
	})
	assert.True(t, success)

	success = srv.DeviceCommand(deviceKey, model.Dictionary{
		"command": "power",
		"payload": false,
	})
	assert.True(t, success)

	//success = srv.DeviceCommand(deviceKey, model.Dictionary{
	//	"command": "power",
	//	"payload": 1,
	//})
	//assert.True(t, success)

	repoMock.AssertExpectations(t)
	restServiceMock.AssertExpectations(t)
}

func TestDeviceCommandRestProperty(t *testing.T) {
	deviceKey := "test-device"
	url := "http://192.168.10.10"
	path := "/cm?cmnd=Status"
	device := model.Device{
		DeviceKey: deviceKey,
		Rest: model.DeviceRest{
			Url: url,
			CommandPaths: model.Dictionary{
				"status": model.Dictionary{
					"path": path,
				},
			},
		},
	}
	repoMock := new(mock_test.DeviceRepoMock)
	restServiceMock := new(mock_test.RestServiceMock)
	srv := DeviceService{
		deviceRepo: repoMock,
		restService: restServiceMock,
	}
	repoMock.On("GetDevice", deviceKey).Return(&device, nil)
	restServiceMock.On("GetRequest", url+path, &device).Times(1).Return(true)

	success := srv.DeviceCommand(deviceKey, model.Dictionary{
		"command": "status",
	})
	assert.True(t, success)

	success = srv.DeviceCommand(deviceKey, model.Dictionary{
		"command": "unknown",
	})
	assert.False(t, success)

	repoMock.AssertExpectations(t)
	restServiceMock.AssertExpectations(t)
}

func TestDeviceStates(t *testing.T) {
	sharedMemoryMock := new(mock_test.SharedMemoryMock)
	srv := DeviceService{
		sharedMemory: sharedMemoryMock,
	}
	sharedMemoryMock.On("GetMemory").Times(1).Return(model.Dictionary{})
	srv.DeviceStates()

	sharedMemoryMock.AssertExpectations(t)
}

func TestDeviceState(t *testing.T)   {
	deviceKey := "test-device"
	sharedMemoryMock := new(mock_test.SharedMemoryMock)
	srv := DeviceService{
		sharedMemory: sharedMemoryMock,
	}
	sharedMemoryMock.On("GetDeviceMem", deviceKey).Times(1).Return(nil)
	srv.DeviceState(deviceKey)

	sharedMemoryMock.AssertExpectations(t)
}
func TestGetDevice(t *testing.T)     {
	deviceKey := "test-device"
	repoMock := new(mock_test.DeviceRepoMock)
	srv := DeviceService{
		deviceRepo: repoMock,
	}
	repoMock.On("GetDevice", deviceKey).Times(1).Return(&model.Device{}, nil)

	srv.GetDevice(deviceKey)
	repoMock.AssertExpectations(t)
}
func TestGetDevices(t *testing.T)    {
	repoMock := new(mock_test.DeviceRepoMock)
	srv := DeviceService{
		deviceRepo: repoMock,
	}
	repoMock.On("GetDevices").Times(1).Return(&model.Devices{}, nil)
	srv.GetDevices()
	repoMock.AssertExpectations(t)
}
func TestGetDevicesDto(t *testing.T) {
	devices := model.Devices{
		Devices: []model.Device{
			{
				Caption: "dev1",
			},
			{
				Caption:"dev2",
			},
		},
	}
	repoMock := new(mock_test.DeviceRepoMock)
	srv := DeviceService{
		deviceRepo: repoMock,
	}
	repoMock.On("GetDevices").Times(1).Return(&devices, nil)
	filtered, err := srv.GetDevicesDto(func(device *model.Device) bool { return device.Caption == "dev1"})
	assert.Len(t, filtered.Devices, 1)
	assert.Nil(t, err)
	repoMock.AssertExpectations(t)
}
func TestDeviceUpdateBrodcast(t *testing.T)  {
	deviceKey := "test-device"
	devState := model.Dictionary{
		"property": "value1",
	}
	oldState, _ := json.Marshal(model.Dictionary{
		"property": "value2",
	})
	dev := model.Device{
		DeviceKey: deviceKey,
	}
	broadcasted := false
	broadcastWebsocket := func (dev *model.Device, devState interface{}) {
		broadcasted = true
	}
	sharedMemoryMock := new(mock_test.SharedMemoryMock)
	srv := DeviceService{
		sharedMemory: sharedMemoryMock,
		webSocketBroadcast: broadcastWebsocket,
	}

	sharedMemoryMock.On("GetDeviceMem", deviceKey).Times(1).Return(devState)
	srv.DeviceUpdate(&dev, string(oldState))

	assert.True(t, broadcasted)
	sharedMemoryMock.AssertExpectations(t)
}

func TestDeviceUpdateNoBroadcast(t *testing.T)  {
	deviceKey := "test-device"
	devState := model.Dictionary{
		"property": "value",
	}
	oldState, _ := json.Marshal(devState)
	dev := model.Device{
		DeviceKey: deviceKey,
	}
	broadcasted := false
	broadcastWebsocket := func (dev *model.Device, devState interface{}) {
		broadcasted = true
	}

	sharedMemoryMock := new(mock_test.SharedMemoryMock)
	srv := DeviceService{
		sharedMemory: sharedMemoryMock,
		webSocketBroadcast: broadcastWebsocket,
	}
	sharedMemoryMock.On("GetDeviceMem", deviceKey).Times(1).Return(devState)
	srv.DeviceUpdate(&dev, string(oldState))

	assert.False(t, broadcasted)
	sharedMemoryMock.AssertExpectations(t)
}

