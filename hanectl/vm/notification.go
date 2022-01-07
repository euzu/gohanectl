package vm

import (
	"errors"
	"fmt"
	"github.com/robertkrimen/otto"
	"github.com/rs/zerolog/log"
	"gohanectl/hanectl/config"
	"gohanectl/hanectl/model"
	"reflect"
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

func compareValues(v1 interface{}, v2 interface{}) bool {
	t1 := reflect.TypeOf(v1)
	t2 := reflect.TypeOf(v2)
	if t1 != t2 {
		return fmt.Sprintf("%v", v1) != fmt.Sprintf("%v", v2)
	} else {
		return v1 != v2
	}
}

func (j *JScriptVM) notify(deviceKey string, key string, newValue interface{}, oldValue interface{}) {
	if notifications, err := j.notificationService.GetNotifications(deviceKey, key); err == nil {
		if dev, err := j.deviceService.GetDevice(deviceKey); err == nil {
			var valuesNotEqual = compareValues(newValue, oldValue)
			for _, n := range notifications {
				if n.Trigger == model.TriggerEvent  || (n.Trigger == model.TriggerChange && valuesNotEqual) {
					if script, err := j.getNotifyScript(dev, n); err == nil {
						params := make(model.Dictionary)
						params["key"] = key
						params["newValue"] = newValue
						params["oldValue"] = oldValue
						//params["scriptName"] = n.Template
						j.runNotifyScript(&HandlerScript{device: dev, script: script, template: n.Caption}, params)
					}
				}
			}
		} else {
			log.Error().Msgf("could not find device %s for notification", deviceKey)
		}
	}
}
