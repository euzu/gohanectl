package vm

import (
	"errors"
	"fmt"
	"github.com/robertkrimen/otto"
	"github.com/rs/zerolog/log"
	"gohanectl/hanectl/config"
	"gohanectl/hanectl/model"
	"strings"
)

func getMqttScriptTemplatePath(cfg config.IConfiguration, templateName string) (string, error) {
	return getScriptTemplatesNestedPath(cfg, templateName, config.ScriptsTemplatesEventsMqttDirectory, config.DefScriptsTemplatesEventsMqttDirectory)
}

func getMqttScriptForTemplate(cfg config.IConfiguration, dev *model.Device, scriptName string) (string, error) {
	if templateFile, err := getMqttScriptTemplatePath(cfg, scriptName); err == nil {
		return getParsedScriptTemplate(dev, templateFile)
	}
	return "", errors.New("not found")
}

func (j *JScriptVM) getMqttScript(cfg config.IConfiguration, dev *model.Device, topic *model.DeviceTopic) (*otto.Script, error) {
	scriptKey := fmt.Sprintf("%s/%s", dev.DeviceKey, topic.HandlerTemplate)
	if val, ok := j.mqttScripts[scriptKey]; ok {
		return val, nil
	} else {
		scriptName := topic.HandlerTemplate
		if content, err := getMqttScriptForTemplate(cfg, dev, scriptName); err == nil {
			script := j.compile(scriptName, content)
			j.mqttScripts[scriptKey] = script
			return script, nil
		}
	}
	return nil, errors.New("not found")
}

func (j *JScriptVM) getMqttHandlerScript(cfg config.IConfiguration, topic string) ([]HandlerScript, error) {
    var result []HandlerScript
	if devices, err := j.deviceService.GetDevices(); err == nil {
		for i := range devices.Devices {
			xDev := &devices.Devices[i]
			for _, t := range xDev.Mqtt.ListenTopics {
				if t.Matches(topic) {
					if script, err := j.getMqttScript(cfg, xDev, &t); err == nil {
						result = append(result, HandlerScript{ device: xDev, script: script, template: t.HandlerTemplate})
					}
				}
			}
		}
		return result, nil
	}
	return result, errors.New(fmt.Sprintf("Could not find script for topic: %s", topic))
}

func (j *JScriptVM) handleMqttMessage(topic string, payload []byte) {
	params := make(model.Dictionary)
	params["topic"] = strings.ToLower(topic)
	params["payload"] = strings.ToLower(string(payload))
	if handlerScripts, err := j.getMqttHandlerScript(j.config, topic); err == nil {
		for i, hs := range handlerScripts {
			hs.template = fmt.Sprintf("mqtt/%s", hs.template)
			j.runHandlerScript(&handlerScripts[i], params)
		}
	} else {
		log.Debug().Msgf("Failed to handle mqtt message: %v", err)
	}
}
