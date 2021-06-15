package mock_test

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/stretchr/testify/mock"
	"time"
)

type MqttTokenMock struct {
	mock.Mock
}

func (m *MqttTokenMock) Wait() bool {
	args := m.Called()
	return args.Bool(0)
}

func (m *MqttTokenMock) WaitTimeout(duration time.Duration) bool {
	args := m.Called(duration)
	return args.Bool(0)
}

func (m *MqttTokenMock) Done() <-chan struct{} {
	args := m.Called()
	return args.Get(0).(<-chan struct{})
}

func (m *MqttTokenMock) Error() error {
	args := m.Called()
	return args.Error(0)
}

type MqttClientMock struct {
	mock.Mock
}

func (m *MqttClientMock) IsConnected() bool {
	args := m.Called()
	return args.Bool(0)
}
func (m *MqttClientMock) IsConnectionOpen() bool {
	args := m.Called()
	return args.Bool(0)
}
func (m *MqttClientMock) Connect() mqtt.Token {
	args := m.Called()
	return args.Get(0).(mqtt.Token)
}
func (m *MqttClientMock) Disconnect(quiesce uint) {
	m.Called(quiesce)
}
func (m *MqttClientMock) Publish(topic string, qos byte, retained bool, payload interface{}) mqtt.Token {
	args := m.Called(topic, qos, retained, payload)
	return args.Get(0).(mqtt.Token)
}
func (m *MqttClientMock) Subscribe(topic string, qos byte, callback mqtt.MessageHandler) mqtt.Token {
	args := m.Called(topic, qos, callback)
	return args.Get(0).(mqtt.Token)
}
func (m *MqttClientMock) SubscribeMultiple(filters map[string]byte, callback mqtt.MessageHandler) mqtt.Token {
	args := m.Called(filters, callback)
	return args.Get(0).(mqtt.Token)
}
func (m *MqttClientMock) Unsubscribe(topics ...string) mqtt.Token {
	args := m.Called(topics)
	return args.Get(0).(mqtt.Token)
}
func (m *MqttClientMock) AddRoute(topic string, callback mqtt.MessageHandler) {
	m.Called(topic, callback)
}
func (m *MqttClientMock) OptionsReader() mqtt.ClientOptionsReader {
	args := m.Called()
	return args.Get(0).(mqtt.ClientOptionsReader)
}
