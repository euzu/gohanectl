package app

import (
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/rs/zerolog/log"
	"gohanectl/hanectl/config"
	"gohanectl/hanectl/model"
	"gohanectl/hanectl/utils"
	"gohanectl/hanectl/websocket"
)

var mqttClientFactory = mqtt.NewClient

func startMqttClient(cfg config.IConfiguration, services model.IServiceFactory) {
	mqttService := services.GetMqttService()
	var broker = cfg.GetStr(config.MqttHost, config.DefMqttHost)
	var port = cfg.GetInt(config.MqttPort, config.DefMqttPort)
	var username = cfg.GetStr(config.MqttUsername, "")
	var password = cfg.GetStr(config.MqttPassword, "")
	opts := mqtt.NewClientOptions()
	opts.AutoReconnect = true
	opts.AddBroker(fmt.Sprintf("tcp://%s:%d", broker, port))
	opts.SetClientID(cfg.GetStr(config.MqttClientId, config.DefMqttClientId))
	opts.SetUsername(username)
	opts.SetPassword(password)
	opts.SetDefaultPublishHandler(func(client mqtt.Client, msg mqtt.Message) {
		//log.Info().Msgf("Received message: %s from topic: %s\n", msg.Payload(), msg.Topic())
		services.GetMqttService().HandleMessage(msg.Topic(), msg.Payload())
	})
	opts.OnConnect = func(client mqtt.Client) {
		log.Debug().Msg("Mqtt connected")
		services.GetConfigService().SetMqttStatus(1)
		websocket.Broadcast(utils.Json(map[string]int{"mqtt": 1}))
	}
	opts.OnConnectionLost = func(client mqtt.Client, err error) {
		log.Info().Msgf("Mqtt connect lost: %v", err)
		services.GetConfigService().SetMqttStatus(0)
		websocket.Broadcast(utils.Json(map[string]int{"mqtt": 0}))
	}

	client := mqttClientFactory(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatal().Msgf("Failed to connect to MQTT: %v", token.Error())
	}

	mqttService.SetClient(client)

	// get topics
	var topics = cfg.GetList(config.MqttTopics)
	if topics == nil {
		mqttService.Subscribe(config.DefMqttTopic)
	} else {
		for i := range topics {
			mqttService.Subscribe(topics[i].(string))
		}
	}
}
