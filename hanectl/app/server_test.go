package app

import (
	"gohanectl/hanectl/config"
	"gohanectl/hanectl/test/mock_test"
	"testing"
)

func TestCreateListener(t *testing.T) {
	cfgMock := new(mock_test.ConfigurationMock)
	cfgMock.On("GetStr", config.ListenHost, config.DefHost).Return("localhost")
	cfgMock.On("GetInt", config.ListenPort, config.DefPort).Return(8888)

	createListener(cfgMock)

	cfgMock.AssertExpectations(t)
}