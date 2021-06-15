package service

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gohanectl/hanectl/test/mock_test"
	"testing"
	"time"
)

func TestPublish(t *testing.T) {
	srv := new(MqttService)
	clientMock := new (mock_test.MqttClientMock)
	tokenMock := new(mock_test.MqttTokenMock)
	srv.SetClient(clientMock)

	clientMock.On("IsConnected").Return(true)
	clientMock.On("Publish", "topic", byte(0), false, "payload").Return(tokenMock)
	tokenMock.On("WaitTimeout", 5 * time.Second).Return(true)

	success := srv.Publish("topic", "payload")
	assert.True(t, success)

	tokenMock.AssertExpectations(t)
	clientMock.AssertExpectations(t)
}

func TestPublishNotConnected(t *testing.T) {
	srv := new(MqttService)
	clientMock := new (mock_test.MqttClientMock)
	srv.SetClient(clientMock)

	clientMock.On("IsConnected").Return(false)

	success := srv.Publish("topic", "payload")
	assert.False(t, success)

	clientMock.AssertExpectations(t)
}

func TestSubscribe(t *testing.T) {
	srv := new(MqttService)
	clientMock := new (mock_test.MqttClientMock)
	tokenMock := new(mock_test.MqttTokenMock)
	srv.SetClient(clientMock)

	clientMock.On("Subscribe", "topic", byte(1), mock.AnythingOfType("mqtt.MessageHandler")).Return(tokenMock)
	tokenMock.On("WaitTimeout", 10 * time.Second).Return(true)

	success := srv.Subscribe("topic")
	assert.True(t, success)

	tokenMock.AssertExpectations(t)
	clientMock.AssertExpectations(t)
}

func TestHandleMessage(t *testing.T) {
	called := false
	handler := func(string, []byte) {
		called = true
	}
	srv := new(MqttService)
	clientMock := new (mock_test.MqttClientMock)
	srv.SetClient(clientMock)
	srv.SetMessageHandler(handler)

	srv.HandleMessage("topic", nil)

	assert.True(t, called)
}
