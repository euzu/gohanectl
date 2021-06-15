package mock_test

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/stretchr/testify/mock"
	"gohanectl/hanectl/model"
)

type MqttServiceMock struct {
	mock.Mock
}

func (m *MqttServiceMock) SetClient(client mqtt.Client) {
	m.Called(client)
}
func (m *MqttServiceMock) Publish(topic string, payload interface{}) bool {
	args := m.Called(topic, payload)
	return args.Bool(0)
}
func (m *MqttServiceMock) Subscribe(topic string) bool {
	args := m.Called(topic)
	return args.Bool(0)
}

func (m *MqttServiceMock) SetMessageHandler(handler model.MqttMessageHandler) {
	m.Called(handler)
}

func (m *MqttServiceMock) HandleMessage(topic string, payload []byte) {
	m.Called(topic, payload)
}
