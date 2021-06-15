package vm

import (
	"fmt"
	"github.com/robertkrimen/otto"
	"github.com/rs/zerolog/log"
	"gohanectl/hanectl/config"
	"gohanectl/hanectl/model"
	"io/ioutil"
	"os"
	"path"
)

const JsSuffix = ".js"

const KeyParams = "_otto_params"
const KeyGetMem = "_otto_getMem"
const KeySetMem = "_otto_setMem"
const KeyGetParam = "_otto_getParam"
const KeyLog = "_otto_log"
const KeyTelegram = "_otto_telegram"
const KeyMqtt = "_otto_mqtt"

type ScriptMap = map[string]*otto.Script

var jScriptVM *JScriptVM

type JScriptVM struct {
	vm                  *otto.Otto
	config              config.IConfiguration
	mqttScripts         ScriptMap
	restScripts         ScriptMap
	notifyScripts       ScriptMap
	deviceService       model.IDeviceService
	notificationService model.INotificationService
	sharedMemory        model.ISharedMemory
}

func newJScriptVM(cfg config.IConfiguration, serviceFactory model.IServiceFactory) *JScriptVM {
	vm := JScriptVM{
		vm:                  otto.New(),
		config:              cfg,
		mqttScripts:         make(ScriptMap),
		restScripts:         make(ScriptMap),
		notifyScripts:       make(ScriptMap),
		deviceService:       serviceFactory.GetDeviceService(),
		notificationService: serviceFactory.GetNotificationService(),
		sharedMemory:        serviceFactory.GetSharedMemory(),
	}
	vm.init(serviceFactory.GetTelegramService(), serviceFactory.GetMqttService())
	return &vm
}

func ReloadScripts() error {
	jScriptVM.clearScriptCache()
	return nil
}

func HandleMqttMessage(topic string, payload []byte) {
	jScriptVM.handleMqttMessage(topic, payload)
}

func HandleRestMessage(device *model.Device, payload string) {
	jScriptVM.handleRestResponse(device, payload)
}

func (j *JScriptVM) clearScriptCache() {
	j.mqttScripts = make(ScriptMap)
	j.restScripts = make(ScriptMap)
	j.notifyScripts = make(ScriptMap)
}

func getSendTelegram(telegramService model.ITelegramService) func(call otto.FunctionCall) otto.Value {
	return func(call otto.FunctionCall) otto.Value {
		if message, err := call.Argument(0).ToString(); err == nil {
			telegramService.SendMessage(message)
			return otto.TrueValue()
		}
		return otto.FalseValue()
	}
}

func getSendMqtt(mqttService model.IMqttService) func(call otto.FunctionCall) otto.Value {
	return func(call otto.FunctionCall) otto.Value {
		if topic, err := call.Argument(0).ToString(); err == nil {
			if payload, err := call.Argument(1).ToString(); err == nil {
				mqttService.Publish(topic, payload)
				return otto.TrueValue()
			}
		}
		return otto.FalseValue()
	}
}

func logJS(call otto.FunctionCall) otto.Value {
	if level, err := call.Argument(0).ToString(); err == nil {
		if message, err := call.Argument(1).ToString(); err == nil {
			switch level {
			case "error":
				log.Error().Msgf("vm: %s", message)
			case "warn":
				log.Warn().Msgf("vm: %s", message)
			case "info":
				log.Info().Msgf("vm: %s", message)
			default:
				log.Debug().Msgf("vm: %s", message)
			}
			return otto.TrueValue()
		}
	}
	return otto.FalseValue()
}

func getParam(call otto.FunctionCall) otto.Value {
	key, _ := call.Argument(0).ToString()
	params := call.Argument(1).Object()
	if result, err := params.Get(key); err == nil {
		return result
	}
	return otto.NullValue()
}

func getSharedMemSetter(memory model.ISharedMemory) func(call otto.FunctionCall) otto.Value {
	return func(call otto.FunctionCall) otto.Value {
		deviceKey, _ := call.Argument(0).ToString()
		field, _ := call.Argument(1).ToString()
		value, _ := call.Argument(2).Export()
		//log.Debug().Msgf("SetSharedMemory deviceKey: %s, field: %s, value: %v", deviceKey, field, value)
		memory.SetMem(deviceKey, field, value)
		return otto.TrueValue()
	}
}

func getSharedMemGetter(memory model.ISharedMemory) func(call otto.FunctionCall) otto.Value {
	return func(call otto.FunctionCall) otto.Value {
		deviceKey, _ := call.Argument(0).ToString()
		field, _ := call.Argument(1).ToString()
		value, _ := otto.ToValue(memory.GetMem(deviceKey, field))
		return value
	}
}

func (j *JScriptVM) init(telegramService model.ITelegramService, mqttService model.IMqttService) {
	if err := j.vm.Set(KeyGetMem, getSharedMemGetter(j.sharedMemory)); err != nil {
		log.Error().Msgf("Could not set getter for shared memory: %v", err)
	}
	if err := j.vm.Set(KeySetMem, getSharedMemSetter(j.sharedMemory)); err != nil {
		log.Error().Msgf("Could not set setter for shared memory: %v", err)
	}
	if err := j.vm.Set(KeyGetParam, getParam); err != nil {
		log.Error().Msgf("Could not set getter for params: %v", err)
	}
	if err := j.vm.Set(KeyLog, logJS); err != nil {
		log.Error().Msgf("Could not set default log: %v", err)
	}
	if err := j.vm.Set(KeyTelegram, getSendTelegram(telegramService)); err != nil {
		log.Error().Msgf("Could not set default telegram: %v", err)
	}
	if err := j.vm.Set(KeyMqtt, getSendMqtt(mqttService)); err != nil {
		log.Error().Msgf("Could not set default mqtt: %v", err)
	}
	defaultLib := j.config.GetStr(config.ScriptsDefaultLib, config.DefScriptsDefaultLib)
	lib := fmt.Sprintf("%s%s", path.Join(getScriptPath(j.config), defaultLib), JsSuffix)
	if _, err := os.Stat(lib); err == nil {
		libContent, err := ioutil.ReadFile(lib)
		_, err = j.vm.Run(string(libContent))
		if err != nil {
			log.Error().Msgf("Could not set default lib: %v", err)
		}
	} else {
		log.Error().Msgf("Could not load default lib: %s", lib)
	}
}

func InitVM(cfg config.IConfiguration, services model.IServiceFactory) {
	jScriptVM = newJScriptVM(cfg, services)
	services.GetMqttService().SetMessageHandler(HandleMqttMessage)
	services.GetRestService().SetMessageHandler(HandleRestMessage)
	services.GetSharedMemory().LoadSharedMem()
	services.GetSharedMemory().SetNotifyCallback(jScriptVM.notify)
}

func (j *JScriptVM) compile(name string, content string) *otto.Script {
	if script, err := j.vm.Compile(name, content); err == nil {
		return script
	} else {
		log.Fatal().Msgf("Could not compile script %s, %v: %s", name, err, content)
		return nil
	}
}
