package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/rs/zerolog/log"
	"gohanectl/hanectl/config"
	"gohanectl/hanectl/model"
	"gohanectl/hanectl/repository"
	"gohanectl/hanectl/utils"
	"gohanectl/hanectl/websocket"
	"strings"
	"sync"
	"time"
)

const KeyPayload = "payload"
const KeyPayloadOff = "payload_off"
const KeyPayloadOn = "payload_on"
const KeyPayloadToggle = "payload_toggle"

const KeyPowerOff = "power_off"
const KeyPowerOn = "power_on"
const KeyPowerToggle = "power_toggle"

func getMqttPowerPayload(value interface{}) (string, bool) {
	payloadKey := KeyPayloadOff
	if v, ok := value.(bool); ok {
		if v {
			payloadKey = KeyPayloadOn
		}
	} else {
		if vs, ok := value.(string); ok {
			vs = strings.ToLower(vs);
			if strings.Compare(vs, "off") == 0 {
				payloadKey = KeyPayloadOff
			} else if strings.Compare(vs, "on") == 0 {
				payloadKey = KeyPayloadOn
			} else if strings.Compare(vs, "toggle") == 0 {
				payloadKey = KeyPayloadToggle
			}
		} else {
			return "", false
		}
	}
	return payloadKey, true
}

func getRestPowerQuery(value interface{}) (string, bool) {
	powerKey := KeyPowerOff
	if v, ok := value.(bool); ok {
		if v {
			powerKey = KeyPowerOn
		}
	} else {
		if vs, ok := value.(string); ok {
			vs = strings.ToLower(vs);
			if strings.Compare(vs, "off") == 0 {
				powerKey = KeyPowerOff
			} else if strings.Compare(vs, "on") == 0 {
				powerKey = KeyPowerOn
			} else if strings.Compare(vs, "toggle") == 0 {
				powerKey = KeyPowerToggle
			}
		} else {
			return "", false
		}
	}
	return powerKey, true
}


type DeviceService struct {
	restService        model.IRestService
	mqttService        model.IMqttService
	deviceRepo         model.IDeviceRepo
	sharedMemory       model.ISharedMemory
	devices            *model.Devices
	webSocketBroadcast model.FWebSocketBroadcast
}

func (d *DeviceService) mqttDevicePower(topic string, templateMap model.Dictionary, params model.Dictionary) bool {
	if value, exists := params[KeyPayload]; exists {
		if payloadKey, ok := getMqttPowerPayload(value); ok {
			if payload, payloadOk := templateMap[payloadKey]; payloadOk {
				return d.mqttService.Publish(topic, payload)
			}
		}
	}
	return false
}

func (d *DeviceService) mqttDeviceProperty(topic string, templateMap model.Dictionary, params model.Dictionary) bool {
	if payload, payloadOk := templateMap[KeyPayload]; payloadOk {
		formattedPayload := payload.(string)
		if len(formattedPayload) > 0 {
			formattedPayload = utils.Sprintf(formattedPayload, params)
		}
		return d.mqttService.Publish(topic, formattedPayload)
	}
	return false
}

func (d *DeviceService) mqttDeviceCommand(dev *model.Device, params model.Dictionary) bool {
	if commandKey, cmdKeyOk := params["command"]; cmdKeyOk {
		cmdKey := commandKey.(string)
		if commandCfg, cfgOk := dev.Mqtt.CommandTopics[cmdKey]; cfgOk {
			if commandMap, asserted := commandCfg.(model.Dictionary); asserted {
				if topic, topicOk := commandMap["topic"]; topicOk {
					if cmdKey == "power" {
						return d.mqttDevicePower(topic.(string), commandMap, params)
					} else {
						return d.mqttDeviceProperty(topic.(string), commandMap, params)
					}
				}
			} else {
				log.Warn().Msgf("Cant convert map for mqtt.CommandTopics: %s", dev.DeviceKey)
			}
		}
	}
	return false
}

func (d *DeviceService) restDevicePower(dev *model.Device, templateMap model.Dictionary, params model.Dictionary) bool {
	if value, exists := params[KeyPayload]; exists {
		if pathKey, ok := getRestPowerQuery(value); ok {
			if path, payloadOk := templateMap[pathKey]; payloadOk {
				url := fmt.Sprintf("%s%s", dev.Rest.Url, path)
				return d.restService.GetRequest(url, dev)
			}
		}
	}
	return false
}

func (d *DeviceService) restDeviceProperty(dev *model.Device, templateMap model.Dictionary, params model.Dictionary) bool {
	if path, payloadOk := templateMap["path"]; payloadOk {
		formattedPath := path.(string)
		if len(formattedPath) > 0 {
			formattedPath = utils.Sprintf(formattedPath, params)
		}
		url := fmt.Sprintf("%s%s", dev.Rest.Url, formattedPath)
		log.Debug().Msgf("requesting: %s", url)
		return d.restService.GetRequest(url, dev)
	}
	return false
}

func (d *DeviceService) restDeviceCommand(dev *model.Device, params model.Dictionary) bool {
	if utils.IsNotBlank(dev.Rest.Url) {
		if commandKey, cmdKeyOk := params["command"]; cmdKeyOk {
			if commandCfg, cfgOk := dev.Rest.CommandPaths[commandKey.(string)]; cfgOk {
				if commandMap, asserted := commandCfg.(model.Dictionary); asserted {
					if commandKey == "power" {
						return d.restDevicePower(dev, commandMap, params)
					} else {
						return d.restDeviceProperty(dev, commandMap, params)
					}
				} else {
					log.Warn().Msgf("Cant convert map for rest.CommandPaths: %s", dev.DeviceKey)
				}
			}
		}
	}
	return false
}

func (d *DeviceService) DeviceCommand(deviceKey string, params model.Dictionary) bool {
	if dev, err := d.GetDevice(deviceKey); err == nil {
		if !d.mqttDeviceCommand(dev, params) {
			return d.restDeviceCommand(dev, params)
		} else {
			return true
		}
	}
	return false
}

func (d *DeviceService) DeviceGroupCommand(groupKey string, params model.Dictionary) bool {
	done := false
	devices, err := d.GetDevices()
	if err == nil {
		grpKey := strings.ToLower(groupKey)
		for _, x := range devices.Devices {
			for _, g := range x.Groups {
				if strings.Compare(strings.ToLower(g), grpKey) == 0 {
					params["device"] = x.DeviceKey
					done = d.DeviceCommand(x.DeviceKey, params)
					break
				}
			}
		}
	}
	return done
}

func (d *DeviceService) DeviceStates() model.Dictionary {
	return d.sharedMemory.GetMemory()
}

func (d *DeviceService) DeviceState(deviceKey string) interface{} {
	return d.sharedMemory.GetDeviceMem(deviceKey)
}

func (d *DeviceService) GetDevice(deviceKey string) (*model.Device, error) {
	return d.deviceRepo.GetDevice(deviceKey)
}

func (d *DeviceService) GetDevices() (*model.Devices, error) {
	return d.deviceRepo.GetDevices()
}

func mapDevice(entity model.Device) model.DeviceDto {
	return model.DeviceDto{
		Type:         entity.Type,
		DeviceKey:    entity.DeviceKey,
		Caption:      entity.Caption,
		Confirm:      entity.Confirm,
		Optimistic:   entity.Optimistic,
		Url:          entity.Rest.Url,
		Timeout:      entity.Timeout,
		Invert:       entity.Invert,
		Room:         entity.Room,
		Groups:       entity.Groups,
		Supplemental: entity.Supplemental,
		Expanded:     entity.Expanded,
		Icon:         entity.Icon,
	}
}

func (d *DeviceService) GetDevicesDto(filter model.DeviceFilter) (*model.DevicesDto, error) {
	devices, err := d.GetDevices()
	if err == nil {
		var devicesDto model.DevicesDto
		for _, x := range devices.Devices {
			if filter(&x) {
				devicesDto.Devices = append(devicesDto.Devices, mapDevice(x))
			}
		}
		return &devicesDto, nil
	}
	return nil, errors.New("failed to read devices configuration")
}

func (d *DeviceService) DeviceUpdate(dev *model.Device, oldState string) {
	if devState := d.DeviceState(dev.DeviceKey); devState != nil {
		newState, err := json.Marshal(devState)
		if err == nil && strings.Compare(string(newState), oldState) != 0 {
			d.webSocketBroadcast(dev, devState)
		}
	}
}

func broadcastDeviceState(dev *model.Device, devState interface{}) {
	mem := model.Dictionary{
		"deviceKey": dev.DeviceKey,
		"state":     devState,
	}
	if value, err := json.Marshal(mem); err == nil {
		websocket.Broadcast(string(value))
	} else {
		log.Error().Msgf("Cant convert device mem: %v", err)
	}
}

func (d *DeviceService) ReloadDevices() error {
	return d.deviceRepo.ReloadDevices()
}

var doOnce sync.Once

func checkDeviceTimeout(cfg config.IConfiguration, deviceService *DeviceService, sharedMemory model.ISharedMemory) {
	doOnce.Do(func() {
		go func() {
			defaultTimeout := cfg.GetInt(config.DeviceTimeout, config.DefDeviceTimeout)
			for range time.Tick(time.Duration(defaultTimeout) * time.Second) {
				timeStamp := utils.NowTimestamp()
				if devices, err := deviceService.GetDevices(); err == nil && devices != nil {
					for i, x := range devices.Devices {
						if x.Timeout >= 0 {
							xDev := &devices.Devices[i]
							deviceService.restDeviceCommand(xDev, model.Dictionary{"command": "status"})
							lastUpdated := sharedMemory.GetLastUpdated(xDev.DeviceKey)
							if timeStamp-lastUpdated > (x.Timeout * 1000) {
								mem := sharedMemory.GetDeviceMem(xDev.DeviceKey)
								if mem == nil {
									mem = model.Dictionary{KeyLastUpdated: lastUpdated}
								}
								broadcastDeviceState(xDev, mem)
							}
						}
					}
				}
			}
		}()
	})
}

func NewDeviceService(cfg config.IConfiguration, restService model.IRestService, mqttService model.IMqttService, sharedMemory model.ISharedMemory) model.IDeviceService {
	deviceService := &DeviceService{
		restService:        restService,
		mqttService:        mqttService,
		sharedMemory:       sharedMemory,
		deviceRepo:         repository.NewDeviceRepo(cfg),
		webSocketBroadcast: broadcastDeviceState,
	}
	checkDeviceTimeout(cfg, deviceService, sharedMemory)
	return deviceService
}
