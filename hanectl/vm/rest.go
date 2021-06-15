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

func getRestScriptTemplatePath(cfg config.IConfiguration, templateName string) (string, error) {
	return getScriptTemplatesNestedPath(cfg, templateName, config.ScriptsTemplatesEventsRestDirectory, config.DefScriptsTemplatesEventsRestDirectory)
}

func getRestScriptForTemplate(cfg config.IConfiguration, dev *model.Device, scriptName string) (string, error) {
	if templateFile, err := getRestScriptTemplatePath(cfg, scriptName); err == nil {
		return getParsedScriptTemplate(dev, templateFile)
	}
	return "", errors.New("not found")
}

func (j *JScriptVM) getRestScript(dev *model.Device) (*otto.Script, string, error) {
	scriptKey := fmt.Sprintf("%s/%s", dev.DeviceKey, dev.Rest.HandlerTemplate)
	if val, ok := j.restScripts[scriptKey]; ok {
		return val, dev.Rest.HandlerTemplate, nil
	} else {
		scriptName := dev.Rest.HandlerTemplate
		if content, err := getRestScriptForTemplate(j.config, dev, scriptName); err == nil {
			script := j.compile(scriptName, content)
			j.restScripts[scriptKey] = script
			return script, scriptName, nil
		}
	}
	return nil, "", errors.New("not found")
}

func (j *JScriptVM) getRestHandlerScript(device *model.Device) (*HandlerScript , error) {
	if script, scriptName, err := j.getRestScript(device); err == nil {
		return &HandlerScript{device: device, script: script, template: scriptName}, nil
	}
	return nil, errors.New(fmt.Sprintf("Could not find script for device: %s", device.DeviceKey))
}

func (j *JScriptVM) handleRestResponse(device *model.Device, payload string) {
	log.Debug().Msgf("Rest response: %s", payload)
	params := make(model.Dictionary)
	params["payload"] = strings.ToLower(payload)
	if hs, err := j.getRestHandlerScript(device); err == nil {
		hs.template = fmt.Sprintf("rest/%s", hs.template)
		j.runHandlerScript(hs, params)
	} else {
		log.Debug().Msgf("Failed to handle rest response: %v", err)
	}
}

