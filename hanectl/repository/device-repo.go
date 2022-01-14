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

func (d *DeviceRepo) assignMqttTemplate(dev *model.Device, template *model.DeviceTemplate) {
	mqtt := &(dev.Mqtt)
	if len(mqtt.ListenTopics) == 0 {
		mqtt.ListenTopics = template.Mqtt.ListenTopics
	}
	template.Mqtt.CommandTopics = utils.MapContentToStringKey(template.Mqtt.CommandTopics)
	if mqtt.CommandTopics != nil {
		mqtt.CommandTopics = mergeMaps(utils.MapContentToStringKey(mqtt.CommandTopics), template.Mqtt.CommandTopics)
	} else {
		mqtt.CommandTopics = template.Mqtt.CommandTopics
	}
}

func (d *DeviceRepo) assignRestTemplate(dev *model.Device, template *model.DeviceTemplate) {
	rest := &(dev.Rest)
	template.Rest.CommandPaths = utils.MapContentToStringKey(template.Rest.CommandPaths)
	if rest.CommandPaths != nil {
		rest.CommandPaths = mergeMaps(utils.MapContentToStringKey(rest.CommandPaths), template.Rest.CommandPaths)
	} else {
		rest.CommandPaths = template.Rest.CommandPaths
	}
	if utils.IsBlank(rest.HandlerTemplate) {
		rest.HandlerTemplate = template.Rest.HandlerTemplate
	}
}

func (d *DeviceRepo) assignDeviceTemplate(dev *model.Device, template *model.DeviceTemplate) {
	d.assignMqttTemplate(dev, template)
	d.assignRestTemplate(dev, template)

	if len(dev.Supplemental) == 0 {
		dev.Supplemental = template.Supplemental
	}

	if dev.Timeout == 0 {
		if template.Timeout != 0 {
			dev.Timeout = template.Timeout
		}
	}
	if !dev.Optimistic.Present() {
		dev.Optimistic = model.NewBool(template.Optimistic)
	}

	if dev.Icon == "" {
		dev.Icon = template.Icon
	}
}

type DeviceTemplateParams struct {
	DeviceID string
	Url string
}

func (d *DeviceRepo) loadTemplates(devices *model.Devices) {
	for i := range devices.Devices {
		xDev := &devices.Devices[i]
		devParams := DeviceTemplateParams{
			DeviceID: xDev.Mqtt.DeviceID,
			Url: xDev.Rest.Url,
		}
		xTemplate := xDev.Template
		if utils.IsNotBlank(xTemplate) {
			if templatePath, err := getDevicesConfigPath(d.config, xTemplate); err == nil {

				var deviceTemplate = new(model.DeviceTemplate)
				if content, err := getDeviceTemplateContent(templatePath, &devParams); err == nil {
					if _, err := readConfigurationFromContent(templatePath, content, deviceTemplate); err == nil {
						d.assignDeviceTemplate(&devices.Devices[i], deviceTemplate)
					} else {
						log.Fatal().Msgf("could not read device template  %s for %s", xTemplate, xDev.DeviceKey)
					}
				} else {
					log.Fatal().Msgf("could not read device template  %s for %s", xTemplate, xDev.DeviceKey)
				}
			} else {
				log.Fatal().Msgf("could not find device template  %s for %s", xTemplate, xDev.DeviceKey)
			}
		}
	}
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
	d.devices = devices
	return devices, nil
}

func (d *DeviceRepo) Close() {
	d.devicesLock.Lock()
	defer d.devicesLock.Unlock()
	d.devices = nil
	log.Info().Msg("Devices cleared")
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
	d.devicesLock.RLock()
	if d.devices != nil {
		d.devicesLock.RUnlock()
		return d.devices, nil
	}
	d.devicesLock.RUnlock()
	d.devicesLock.Lock()
	d.devicesLock.Unlock()
	return d.loadDevices()
}

func NewDeviceRepo(cfg config.IConfiguration) model.IDeviceRepo {
	return &DeviceRepo{
		config: cfg,
	}
}
