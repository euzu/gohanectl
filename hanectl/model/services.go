package model

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type MqttMessageHandler = func(string, []byte)
type RestMessageHandler = func(device *Device, payload string)
type FWebSocketBroadcast = func(dev *Device, devState interface{})
type NotifyFunc func(deviceKey string, key string, newValue interface{}, oldValue interface{})

type IConfigService interface {
	GetServerStatus() (*ServerStatus, error)
	GetRoomsConfig() (Dictionary, error)
	SetWebsocketStatus(connected bool)
	SetMqttStatus(connectionStatus int)
}

type IDeviceService interface {
	DeviceCommand(deviceKey string, params Dictionary) bool
	DeviceGroupCommand(groupKey string, params Dictionary) bool
	DeviceStates() Dictionary
	DeviceState(deviceKey string) interface{}
	GetDevice(deviceKey string) (*Device, error)
	GetDevices() (*Devices, error)
	GetDevicesDto(filter DeviceFilter) (*DevicesDto, error)
	DeviceUpdate(dev *Device, oldState string)
	ReloadDevices() error
}

type IMqttService interface {
	SetClient(client mqtt.Client)
	Publish(topic string, payload interface{}) bool
	Subscribe(topic string) bool
	SetMessageHandler(handler MqttMessageHandler)
	HandleMessage(topic string, payload []byte)
}

type INotificationService interface {
	GetNotifications(deviceKey string, key string) ([]*Notification, error)
	ReloadNotifications() error
}

type ITelegramBotAPI interface {
	Send(c tgbotapi.Chattable) (tgbotapi.Message, error)
}

type ITelegramService interface {
	SendMessage(message string) bool
}

type IUserService interface {
	FindByUsername(userName string) (*User, error)
	ReloadUsers() error
    SaveSettings(userName string, settings *UserSettings) error
	GetSettings(userName string) (*UserSettings, error)
}

type IRestService interface {
	GetRequest(url string, device *Device) bool
	SetMessageHandler(handler RestMessageHandler)
	HandleMessage(device *Device, payload string)
}

type ISharedMemory interface {
	SetMem(deviceKey string, key string, value interface{})
	MarkAsUpdated(deviceKey string)
	GetLastUpdated(deviceKey string) int64
	GetMem(deviceKey string, key string) interface{}
	GetDeviceMem(deviceKey string) interface{}
	GetMemory() Dictionary
	LoadSharedMem()
	SetNotifyCallback(notifyFunc NotifyFunc)
	Persist()
}

type IDatabaseService interface {
	Persist()
    GetUserSettings(userName string) (*UserSettings, error)
	SaveUserSettings(name string, settings *UserSettings) error
}

type IServiceFactory interface {
	GetDeviceService() IDeviceService
	GetConfigService() IConfigService
	GetMqttService() IMqttService
	GetTelegramService() ITelegramService
	GetUserService() IUserService
	GetNotificationService() INotificationService
	GetRestService() IRestService
	GetSharedMemory() ISharedMemory
	GetDatabaseService() IDatabaseService
	Finalize()
}
