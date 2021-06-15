package service

import (
	"gohanectl/hanectl/config"
	"gohanectl/hanectl/model"
	"gohanectl/hanectl/utils"
)


type ConfigService struct {
	config       config.IConfiguration
	serverStatus *model.ServerStatus
}

func (c ConfigService) GetServerStatus() (*model.ServerStatus, error) {
	return c.serverStatus, nil
}

func (c ConfigService) GetRoomsConfig() (model.Dictionary, error) {
	cfg := c.config.GetMap(config.Room)
	if cfg == nil {
		return make(model.Dictionary), nil
	}
	var result map[interface{}]interface{} = cfg
	return utils.MapToStringKey(result), nil
}

func (c ConfigService) SetWebsocketStatus(connected bool) {
	c.serverStatus.Websocket = connected
}

func (c ConfigService) SetMqttStatus(connectionStatus int) {
	c.serverStatus.Mqtt = connectionStatus
}

func NewConfigService(cfg config.IConfiguration) model.IConfigService {
	return &ConfigService{
		config:       cfg,
		serverStatus: new(model.ServerStatus),
	}
}
