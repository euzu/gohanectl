package service

import (
	"github.com/stretchr/testify/assert"
	"gohanectl/hanectl/config"
	"gohanectl/hanectl/model"
	"gohanectl/hanectl/test/mock_test"
	"testing"
)

func TestGetRoomsConfig(t *testing.T) {
	cfg := new(mock_test.ConfigurationMock)
	srv := new(ConfigService)
	srv.config = cfg
	roomCfg := make(map[interface {}]interface {})

	cfg.On("GetMap", config.Room).Return(roomCfg)

    dict, err := srv.GetRoomsConfig()
    assert.NotNil(t, dict)
	assert.Nil(t, err)
}

func TestSetWebsocketStatus(t *testing.T) {
	srv := new(ConfigService)
    serverStatus := new(model.ServerStatus)
    srv.serverStatus = serverStatus
	srv.SetWebsocketStatus(true)

    assert.True(t, serverStatus.Websocket)
}

func TestSetMqttStatus(t *testing.T) {
	srv := new(ConfigService)
	serverStatus := new(model.ServerStatus)
	srv.serverStatus = serverStatus
	srv.SetMqttStatus(1)

	assert.Equal(t, serverStatus.Mqtt, 1)
}