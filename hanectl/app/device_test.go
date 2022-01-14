package app

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/stretchr/testify/mock"
	"gohanectl/hanectl/config"
	"gohanectl/hanectl/model"
	"gohanectl/hanectl/test/mock_test"
	"testing"
)

func TestInitDevices(t *testing.T) {
	cfgMock := new(mock_test.ConfigurationMock)
	serviceFactoryMock := new(mock_test.ServiceFactoryMock)
	deviceServiceMock := new(mock_test.DeviceServiceMock)
	mqttServiceMock := new(mock_test.MqttServiceMock)
	mqttClientMock := new(mock_test.MqttClientMock)
	mqttTockenMock := new(mock_test.MqttTokenMock)
	devices := &model.Devices{
		Devices: []model.Device{
			{
				Mqtt: model.DeviceMqtt{
					ListenTopics: []model.DeviceTopic{
						{
							Topic: "topic",
						},
					},
				},
			},
		},
	}

	cfgMock.On("GetStr", config.MqttHost, config.DefMqttHost).Return("localhost")
	cfgMock.On("GetInt", config.MqttPort, config.DefMqttPort).Return(1883)
	cfgMock.On("GetStr", config.MqttUsername, "").Return("")
	cfgMock.On("GetStr", config.MqttPassword, "").Return("")
	cfgMock.On("GetStr", config.MqttClientId, config.DefMqttClientId).Return("client")
	cfgMock.On("GetList", config.MqttTopics).Return([]interface{}{"topic"})
	serviceFactoryMock.On("GetDeviceService").Return(deviceServiceMock)
	serviceFactoryMock.On("GetMqttService").Return(mqttServiceMock)
	deviceServiceMock.On("GetDevices").Return(devices, nil)
	mqttServiceMock.On("Subscribe", mock.AnythingOfTypeArgument("string")).Return(true)
	mqttServiceMock.On("SetClient", mqttClientMock)
	mqttClientMock.On("Connect").Return(mqttTockenMock)
	mqttTockenMock.On("Wait").Return(true)
	mqttTockenMock.On("Error").Return(nil)

	mqttClientFactory = func(options *mqtt.ClientOptions) mqtt.Client {
		return mqttClientMock
	}
	initDevices(cfgMock, serviceFactoryMock)

	cfgMock.AssertExpectations(t)
	serviceFactoryMock.AssertExpectations(t)
	deviceServiceMock.AssertExpectations(t)
	mqttServiceMock.AssertExpectations(t)
	mqttClientMock.AssertExpectations(t)
	mqttTockenMock.AssertExpectations(t)
}

func TestInitDevicesNoMqtt(t *testing.T) {
	cfgMock := new(mock_test.ConfigurationMock)
	serviceFactoryMock := new(mock_test.ServiceFactoryMock)
	deviceServiceMock := new(mock_test.DeviceServiceMock)
	configServiceMock := new(mock_test.ConfigServiceMock)
	devices := &model.Devices{}

	serviceFactoryMock.On("GetConfigService").Return(configServiceMock)
	serviceFactoryMock.On("GetDeviceService").Return(deviceServiceMock)
	deviceServiceMock.On("GetDevices").Return(devices, nil)
	configServiceMock.On("SetMqttStatus", 2)

	initDevices(cfgMock, serviceFactoryMock)

	cfgMock.AssertExpectations(t)
	serviceFactoryMock.AssertExpectations(t)
	configServiceMock.AssertExpectations(t)
	deviceServiceMock.AssertExpectations(t)
}

func TestMqttListenTopics(t *testing.T) {
	dev := model.Device{
		Mqtt: model.DeviceMqtt{
			ListenTopics: []model.DeviceTopic{
				{
					Topic: "topic",
				},
			},
		},
	}
	mqttServiceMock := new(mock_test.MqttServiceMock)
	mqttServiceMock.On("Subscribe", "topic").Return(true)

	mqttSubscribeTopics(&dev, mqttServiceMock)
}

func TestMqttDeviceStatus(t *testing.T) {
	dev := model.Device{
		Mqtt: model.DeviceMqtt{
			CommandTopics: model.Dictionary{
				"status": model.Dictionary{
					"topic":   "topic",
					"payload": "",
				},
			},
			ListenTopics: []model.DeviceTopic{
				{
					Topic: "topic",
				},
			},
		},
	}
	mqttServiceMock := new(mock_test.MqttServiceMock)
	mqttServiceMock.On("Publish", "topic", "").Return(true)

	mqttDeviceStatus(&dev, mqttServiceMock)
}
