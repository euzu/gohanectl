package app

import (
	"github.com/rs/zerolog/log"
	"gohanectl/hanectl/config"
	"gohanectl/hanectl/model"
)

// returns true if mqtt configured device exists, else false
func initDevices(cfg config.IConfiguration, services model.IServiceFactory) {
	deviceService := services.GetDeviceService()
	if devices, err := deviceService.GetDevices(); err == nil {
		useMqtt := false
		for _, x := range devices.Devices {
			if len(x.Mqtt.ListenTopics) > 0 {
				useMqtt = true
				startMqttClient(cfg, services)
				break
			}
		}
		if useMqtt {
			mqttService := services.GetMqttService()
			for i := range devices.Devices {
				xDev := &devices.Devices[i]
				mqttListenTopics(xDev, mqttService)
				mqttDeviceStatus(xDev, mqttService)
			}
		} else {
			services.GetConfigService().SetMqttStatus(2)
		}
	} else {
		log.Fatal().Msgf("could not initialize devices: %v", err)
	}
}

func mqttListenTopics(x *model.Device, mqttService model.IMqttService) {
	for _, t := range x.Mqtt.ListenTopics {
		if t.Topic != "" {
			topic := t.Topic
			//if !strings.HasSuffix(t.Topic, "/") {
			//	topic = fmt.Sprintf("%s/", topic)
			//}
			//if !strings.HasSuffix(t.Topic, "#") {
			//	topic = fmt.Sprintf("%s#", topic)
			//}
			mqttService.Subscribe(topic)
		}
	}
}

func mqttDeviceStatus(x *model.Device, mqttService model.IMqttService) {
	if x.Mqtt.CommandTopics != nil {
		if statusCfg, ok := x.Mqtt.CommandTopics["status"]; ok {
			if templateMap, asserted := statusCfg.(model.Dictionary); asserted {
				if topic, topicOk := templateMap["topic"]; topicOk {
					if payload, payloadOk := templateMap["payload"]; payloadOk {
						formattedPayload := payload.(string)
						if len(formattedPayload) > 0 {
							//log.Warn().Msgf("cant get initial status for %s. Status payload is not empty. %s", x.DeviceKey, formattedPayload)
							mqttService.Publish(topic.(string), formattedPayload)
						} else {
							mqttService.Publish(topic.(string), "")
						}
					}
				}
			} else {
				log.Warn().Msgf("Cant convert map for mqtt.CommandTopics: %s", x.DeviceKey)
			}
		}
	}
}
