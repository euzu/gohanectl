package service

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/rs/zerolog/log"
	"gohanectl/hanectl/model"
	"time"
)

type MqttService struct {
	client         mqtt.Client
	messageHandler model.MqttMessageHandler
}

func (m *MqttService) SetClient(client mqtt.Client) {
	m.client = client
}

func (m *MqttService) Publish(topic string, payload interface{}) bool {
	if m.client != nil && m.client.IsConnected() {
		token := m.client.Publish(topic, 0, false, payload)
		return token.WaitTimeout(5 * time.Second)
	}
	return false
}

func (m *MqttService) Subscribe(topic string) bool {
	if m.client != nil {
		token := m.client.Subscribe(topic, 1, nil)
		if result := token.WaitTimeout(10 * time.Second); result {
			log.Debug().Msgf("Subscribed to topic: %s", topic)
			return true
		} else {
			log.Error().Msgf("Cant subscribe to topic: %s", topic)
		}
	}
	return false
}

func (m *MqttService) Unsubscribe(topic string) bool {
	if m.client != nil {
		token := m.client.Unsubscribe(topic)
		if result := token.WaitTimeout(5 * time.Second); result {
			log.Debug().Msgf("Unubscribed from topic: %s", topic)
			return true
		} else {
			log.Error().Msgf("Cant unsubscribe from topic: %s", topic)
		}
	}
	return false
}

func (m *MqttService) SetMessageHandler(handler model.MqttMessageHandler) {
	m.messageHandler = handler
}

func (m *MqttService) HandleMessage(topic string, payload []byte) {
	if m.messageHandler != nil {
		m.messageHandler(topic, payload)
	}
}

func (m *MqttService) Close() {
	if m.client != nil {
		m.client.Disconnect(0)
		m.client = nil
	}
}

func NewMqttService() model.IMqttService {
	return new(MqttService)
}
