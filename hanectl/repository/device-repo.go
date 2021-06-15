package repository

import (
	"errors"
	"fmt"
	"github.com/rs/zerolog/log"
	"gohanectl/hanectl/config"
	"gohanectl/hanectl/model"
	"gohanectl/hanectl/utils"
	"strings"
	"sync"
)

// keep everything in destination, source is only optional if not set in dest
func mergeMaps(dest model.Dictionary, source model.Dictionary) model.Dictionary {
	if dest == nil {
		return source
	}

	if source == nil {
		return dest
	}

	result := make(model.Dictionary)
	for key, srcValue := range source {
		if dstValue, exists := dest[key]; exists {
			if mapDestValue, okDst := dstValue.(model.Dictionary); okDst {
				if mapSrcValue, okSrc := srcValue.(model.Dictionary); okSrc {
					result[key] = mergeMaps(mapDestValue, mapSrcValue)
				} else {
					log.Error().Msgf("Failed to merge configuration:  %srcValue and %srcValue", mapDestValue, srcValue)
				}
			} else {
				result[key] = dstValue
			}
		} else {
			result[key] = srcValue
		}
	}

	return result
}

type DeviceRepo struct {
	devices     *model.Devices
	devicesLock sync.RWMutex
	config      config.IConfiguration
}

func (d *DeviceRepo) assignMqttTemplate(dev *model.Device, template *model.DeviceMqtt) {
	mqtt := &(dev.Mqtt)
	if len(mqtt.ListenTopics) == 0 {
		mqtt.ListenTopics = template.ListenTopics
	}
	template.CommandTopics = utils.MapContentToStringKey(template.CommandTopics)
	if mqtt.CommandTopics != nil {
		mqtt.CommandTopics = mergeMaps(utils.MapContentToStringKey(mqtt.CommandTopics), template.CommandTopics)
	} else {
		mqtt.CommandTopics = template.CommandTopics
	}
}

func (d *DeviceRepo) assignRestTemplate(dev *model.Device, template *model.DeviceRest) {
	rest := &(dev.Rest)
	if len(rest.Url) == 0 {
		rest.Url = template.Url
	}
	template.CommandPaths = utils.MapContentToStringKey(template.CommandPaths)
	if rest.CommandPaths != nil {
		rest.CommandPaths = mergeMaps(utils.MapContentToStringKey(rest.CommandPaths), template.CommandPaths)
	} else {
		rest.CommandPaths = template.CommandPaths
	}
	if utils.IsBlank(rest.HandlerTemplate) {
		rest.HandlerTemplate = template.HandlerTemplate
	}
}

func (d *DeviceRepo) loadMqttTemplate(devices *model.Devices) {
	for i := range devices.Devices {
		xDev := &devices.Devices[i]
		xMqtt := &xDev.Mqtt
		if utils.IsNotBlank(xMqtt.Template.Name) {
			if templatePath, err := getDeviceMqttPath(d.config, xMqtt.Template.Name); err == nil {
				var mqttTemplate = new(model.DeviceMqtt)
				if content, err := getDeviceTemplateContent(templatePath, xMqtt.Template); err == nil {
					if _, err := readConfigurationFromContent(templatePath, content, mqttTemplate); err == nil {
						d.assignMqttTemplate(&devices.Devices[i], mqttTemplate)
					} else {
						log.Fatal().Msgf("could not read mqtt template  %s for %s", xMqtt.Template.Name, xDev.DeviceKey)
					}
				} else {
					log.Fatal().Msgf("could not read mqtt template  %s for %s", xMqtt.Template.Name, xDev.DeviceKey)
				}
			} else {
				log.Fatal().Msgf("could not find mqtt template  %s for %s", xMqtt.Template.Name, xDev.DeviceKey)
			}
		}
	}
}

func (d *DeviceRepo) loadRestTemplate(devices *model.Devices) {
	for i := range devices.Devices {
		xDev := &devices.Devices[i]
		xRest := xDev.Rest
		if utils.IsNotBlank(xRest.Template.Name) {
			if templatePath, err := getDeviceRestPath(d.config, xRest.Template.Name); err == nil {
				var restTemplate = new(model.DeviceRest)
				if content, err := getDeviceTemplateContent(templatePath, xRest.Template); err == nil {
					if _, err := readConfigurationFromContent(templatePath, content, restTemplate); err == nil {
						d.assignRestTemplate(&devices.Devices[i], restTemplate)
					} else {
						log.Fatal().Msgf("could not read mqtt template  %s for %s", xRest.Template.Name, xDev.DeviceKey)
					}
				} else {
					log.Fatal().Msgf("could not read mqtt template  %s for %s", xRest.Template.Name, xDev.DeviceKey)
				}
			} else {
				log.Fatal().Msgf("could not find mqtt template  %s for %s", xRest.Template.Name, xDev.DeviceKey)
			}
		}
	}
}

func (d *DeviceRepo) loadTemplates(devices *model.Devices) {
	d.loadMqttTemplate(devices)
	d.loadRestTemplate(devices)
}

func (d *DeviceRepo) prepare(devices *model.Devices) {
	d.loadTemplates(devices)

	defaultTimeout := d.config.GetInt(config.DeviceTimeout, config.DefDeviceTimeout)
	for i := range devices.Devices {
		xDev := &devices.Devices[i]
		if xDev.Timeout == 0 {
			xDev.Timeout = int64(defaultTimeout)
		}
		for j := range xDev.Mqtt.ListenTopics {
			xDev.Mqtt.ListenTopics[j].Prepare()
		}
	}
}

func (d *DeviceRepo) loadDevices() (*model.Devices, error) {
	devices := &model.Devices{}
	if _, err := readConfiguration(d.config, config.DeviceConfig, "", devices); err != nil {
		log.Error().Msgf("Failed to read devices file: %v", err)
		return nil, errors.New("failed to read devices file")
	}

	d.prepare(devices)

	d.devicesLock.Lock()
	d.devices = devices
	d.devicesLock.Unlock()
	return devices, nil
}

func (d *DeviceRepo) ReloadDevices() error {
	if _, err := d.loadDevices(); err == nil {
		log.Info().Msg("Devices reloaded")
		return nil
	} else {
		return err
	}
}

func (d *DeviceRepo) GetDevice(deviceKey string) (*model.Device, error) {
	if devices, err := d.GetDevices(); err == nil {
		for _, x := range devices.Devices {
			if strings.Compare(x.DeviceKey, deviceKey) == 0 {
				return &x, nil
			}
		}
	}
	return nil, errors.New(fmt.Sprintf("cant find device with key %s", deviceKey))
}

func (d *DeviceRepo) GetDevices() (*model.Devices, error) {
	if d.devices != nil {
		d.devicesLock.RLock()
		defer d.devicesLock.RUnlock()
		return d.devices, nil
	}
	return d.loadDevices()
}

func NewDeviceRepo(cfg config.IConfiguration) model.IDeviceRepo {
	return &DeviceRepo{
		config: cfg,
	}
}
