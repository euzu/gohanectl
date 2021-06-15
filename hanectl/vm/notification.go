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

func (j *JScriptVM)  getScriptNotificationPath(notifyName string) (string, error) {
	return getScriptNestedPath(j.config, notifyName, config.ScriptsNotificationDirectory, config.DefScriptsNotificationDirectory)
}

func (j *JScriptVM) getScriptForNotification(dev *model.Device, notification *model.Notification) (string, error) {
	if templateFile, err := j.getScriptNotificationPath(notification.Script); err == nil {
		return getParsedScriptTemplate(dev, templateFile)
	}
	return "", errors.New("not found")
}

func (j *JScriptVM) getNotifyScript(dev *model.Device, notification *model.Notification) (*otto.Script, error) {
	scriptKey := fmt.Sprintf("%s/%s/%s", notification.DeviceKey, strings.Join(notification.Keys, "_"), notification.Script)
	if val, ok := j.notifyScripts[scriptKey]; ok {
		return val, nil
	} else {
		if content, err := j.getScriptForNotification(dev, notification); err == nil {
			script := j.compile(scriptKey, content)
			j.notifyScripts[scriptKey] = script
			return script, nil
		}
	}
	return nil, errors.New("cant find notification script")
}

func (j *JScriptVM) notify(deviceKey string, key string, newValue interface{}, oldValue interface{}) {
	if notifications, err := j.notificationService.GetNotifications(deviceKey, key); err == nil {
		if dev, err := j.deviceService.GetDevice(deviceKey); err == nil {
			for _, n := range notifications {
				if script, err := j.getNotifyScript(dev, n); err == nil {
					params := make(model.Dictionary)
					params["key"] = key
					params["newValue"] = newValue
					params["oldValue"] = oldValue
					//params["scriptName"] = n.Template
					j.runNotifyScript(&HandlerScript{device: dev, script: script, template: n.Caption}, params)
				}
			}
		} else {
			log.Error().Msgf("could not find device %s for notification", deviceKey)
		}
	}
}
