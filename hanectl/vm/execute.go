package vm

import (
	"encoding/json"
	"github.com/robertkrimen/otto"
	"github.com/rs/zerolog/log"
	"gohanectl/hanectl/model"
	"sync"
)

type HandlerScript struct {
	device *model.Device
	script *otto.Script
	template string
}


var scriptMutex sync.RWMutex

//func scriptFilename(s *otto.Script) string {
//	v := reflect.ValueOf(*s)
//	y := v.FieldByName("filename")
//	return fmt.Sprint(y.Interface())
//}

func (j *JScriptVM) runHandlerScript(hs *HandlerScript, params model.Dictionary) {
	scriptMutex.Lock()
	defer scriptMutex.Unlock()
	log.Debug().Msgf("Running %s for device %s with params: %v", hs.template, hs.device.DeviceKey, params)
	err := j.vm.Set(KeyParams, params)
	if err != nil {
		log.Error().Msgf("Failed (cant set params) running handler script for %s: %v.", hs.device.DeviceKey, err)
	} else {
		var oldState string
		if state := j.deviceService.DeviceState(hs.device.DeviceKey); state != nil {
			if oState, err := json.Marshal(state); err == nil {
				oldState = string(oState)
			}
		}
		if _, err := j.vm.Run(hs.script); err != nil {
			log.Error().Msgf("Failed running handler script for %s: %v: %v", hs.device.DeviceKey, err, params)
			log.Error().Msgf("Failed running handler script content %v", hs.script)
		} else {
			j.deviceService.DeviceUpdate(hs.device, oldState)
			j.sharedMemory.MarkAsUpdated(hs.device.DeviceKey)
		}
	}
}

func (j *JScriptVM) runNotifyScript(hs *HandlerScript, params model.Dictionary) {
	scriptMutex.Lock()
	defer scriptMutex.Unlock()
	log.Debug().Msgf("Running notify script %s for %s with params: %v", hs.template, hs.device.DeviceKey, params)
	err := j.vm.Set(KeyParams, params)
	if err != nil {
		log.Error().Msgf("Failed (cant set params) running notify script for %s: %v", hs.device.DeviceKey, err)
	} else {
		if _, err := j.vm.Run(hs.script); err != nil {
			log.Error().Msgf("Failed running notify script for %s: %v", hs.device.DeviceKey, err)
		}
	}
}
