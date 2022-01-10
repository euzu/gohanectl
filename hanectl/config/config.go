package config

import (
	"bytes"
	"encoding/json"
	"github.com/rs/zerolog/log"
	"gohanectl/hanectl/utils"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path"
	"strconv"
	"strings"
)

type ConfigKey string

const (
	LogLevel   = ConfigKey("log.level")
	LogFile    = ConfigKey("log.file")
	LogConsole = ConfigKey("log.console")

	RunUser = ConfigKey("runas")

	WorkingDir = ConfigKey("working_dir")

	Mqtt              = ConfigKey("mqtt")
	MqttHost          = ConfigKey("mqtt.host")
	MqttPort          = ConfigKey("mqtt.port")
	MqttWebsocketPort = ConfigKey("mqtt.websocket_port")
	MqttUsername      = ConfigKey("mqtt.username")
	MqttPassword      = ConfigKey("mqtt.password")
	MqttTopics        = ConfigKey("mqtt.topics")
	MqttClientId      = ConfigKey("mqtt.client_id")

	TelegramBotEnabled = ConfigKey("telegram.enabled")
	TelegramBotToken   = ConfigKey("telegram.bot_token")
	TelegramChatIds    = ConfigKey("telegram.chat_ids")

	ConfigDirectory = ConfigKey("config.directory")
	CommandToken    = ConfigKey("config.command_token")

	ScriptsDirectory          = ConfigKey("config.scripts.directory")
	ScriptsTemplatesDirectory = ConfigKey("config.scripts.templates.directory")

	ScriptsTemplatesEventsMqttDirectory   = ConfigKey("config.scripts.templates.events.mqtt")
	ScriptsTemplatesEventsRestDirectory   = ConfigKey("config.scripts.templates.events.rest")
	ScriptsTemplatesDevicesDirectory = ConfigKey("config.scripts.templates.devices")

	ScriptsNotificationDirectory = ConfigKey("config.scripts.notifications")
	ScriptsDefaultLib            = ConfigKey("config.scripts.default_lib")

	DatabaseStatesName    = ConfigKey("database.states.name")
	DatabaseStatesPersist = ConfigKey("database.states.persist")
	DatabaseSettingsName  = ConfigKey("database.settings.name")
	DatabaseSettingsPersist = ConfigKey("database.settings.persist")

	ListenPort         = ConfigKey("web.listen.port")
	ListenHost         = ConfigKey("web.listen.host")
	WebFilesDir        = ConfigKey("web.web_files")
	JwtSecret          = ConfigKey("jwt.secret")
	CorsAllowedOrigins = ConfigKey("cors.allowed_origins")

	DeviceTimeout      = ConfigKey("config.device_timeout")
	DeviceConfig       = ConfigKey("config.devices")
	UserConfig         = ConfigKey("config.users")
	NotificationConfig = ConfigKey("config.notifications")

	Room = ConfigKey("config.room")
)

const (
	DefLogLevel = "info"

	DefMqttHost          = "localhost"
	DefMqttPort          = 1883
	DefMqttWebsocketPort = 8083
	DefMqttTopic         = "tele/#"
	DefMqttClientId      = "hanectl_mqtt_client"

	DefConfigDirectory                       = "config"
	DefScriptsDirectory                      = "scripts"
	DefScriptsNotificationDirectory          = "notification"
	DefScriptsDefaultLib                     = "_lib"
	DefScriptsTemplatesDirectory             = "templates"
	DefScriptsTemplatesEventsMqttDirectory   = "events/mqtt"
	DefScriptsTemplatesEventsRestDirectory   = "events/rest"
	DefScriptsTemplatesDevicesDirectory      = "devices"

	DefDatabaseStatesName    = "hanectl_states.db"
	DefDatabaseStatesPersist = true
	DefDatabaseSettingsName  = "hanectl_settings.db"
	DefDatabaseSettingsPersist = true

	DefPort = 8500
	DefHost = "127.0.0.1"

	DefWebFilesDir = "web"
	DefWebUrl      = "/"

	DefDeviceTimeout = 60
)

const (
	EnvMqttHost          = "MQTT_HOST"
	EnvMqttPort          = "MQTT_PORT"
	EnvMqttWebsocketPort = "MQTT_WEBSOCKET_PORT"
)

type IConfiguration interface {
	ReadConfiguration(configFile string)
	GetStr(key ConfigKey, defaultValue string) string
	GetInt(key ConfigKey, defaultValue int) int
	GetBool(key ConfigKey, defaultValue bool) bool
	GetList(key ConfigKey) []interface{}
	GetMap(key ConfigKey) map[interface{}]interface{}
}

type Configuration struct {
	config map[interface{}]interface{}
}

func (c *Configuration) ReadConfiguration(configFile string) {
	fileName := configFile
	if configFile == "" {
		candidates := []string{"config.yml", path.Join(DefConfigDirectory, "config.yml"), "config.json", path.Join(DefConfigDirectory, "config.json")}
		for _, x := range candidates {
			if utils.FileExists(x) {
				fileName = x
				break
			}
		}
		if utils.IsBlank(fileName) {
			log.Fatal().Msg("cant find any configuration file")
			return
		}
	} else if !utils.FileExists(configFile) {
		log.Fatal().Msgf("cant find any configuration file %s", configFile)
		return
	}

	content, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Fatal().Msgf("cant read configuration file %s: %v", fileName, err)
		return
	}
	if cwd, err := os.Getwd(); err == nil {
		log.Info().Msgf("Using configuration file: %s", path.Join(cwd, fileName))
	} else {
		log.Info().Msgf("Using configuration file: %s", fileName)
	}

	if strings.HasSuffix(fileName, ".yml") {
		if err := yaml.Unmarshal(content, &c.config); err != nil {
			log.Fatal().Msgf("Failed to read configuration %v", err)
		}
	} else if strings.HasSuffix(fileName, ".json") {
		decoder := json.NewDecoder(bytes.NewReader(content))
		if err := decoder.Decode(&c.config); err != nil {
			log.Fatal().Msgf("Failed to read configuration %v", err)
		}
	}

	mqttCfg, _ := c.getConfig(Mqtt)
	value := os.Getenv(EnvMqttPort)
	if value != "" {
		v, err := strconv.Atoi(value)
		if err == nil {
			mqttCfg[string(MqttPort)] = v
		}
	}
	value = os.Getenv(EnvMqttWebsocketPort)
	if value != "" {
		v, err := strconv.Atoi(value)
		if err == nil {
			mqttCfg[string(MqttWebsocketPort)] = v
		}
	}

	value = os.Getenv(EnvMqttHost)
	if value != "" {
		mqttCfg[string(MqttHost)] = value
	}
}

func (c *Configuration) getConfig(key ConfigKey) (map[interface{}]interface{}, string) {
	var result map[interface{}]interface{} = nil
	keys := strings.Split(string(key), ".")
	length := len(keys)
	result = c.config
	for i := 0; i < length-1; i++ {
		t, ok := result[keys[i]].(map[interface{}]interface{})
		if !ok {
			return nil, ""
		}
		result = t
	}
	return result, keys[length-1]
}

func (c *Configuration) getValue(key ConfigKey) interface{} {
	cfg, ckey := c.getConfig(key)
	if cfg != nil {
		return cfg[interface{}(ckey)]
	}
	return nil
}

func (c *Configuration) GetStr(key ConfigKey, defaultValue string) string {
	if value := c.getValue(key); value != nil {
		if t, ok := value.(string); ok {
			return t
		} else {
			log.Error().Msgf("Config values does not match string type %s: %v", key, value)
		}
	}
	return defaultValue
}

func (c *Configuration) GetInt(key ConfigKey, defaultValue int) int {
	if value := c.getValue(key); value != nil {
		if t, ok := value.(int); ok {
			return int(t)
		} else if t, ok := value.(float64); ok {
			return int(t)
		} else {
			log.Error().Msgf("Config values does not match int type %s: %v", key, value)
		}
	}
	return defaultValue
}

func (c *Configuration) GetBool(key ConfigKey, defaultValue bool) bool {
	if value := c.getValue(key); value != nil {
		if t, ok := value.(bool); ok {
			return t
		} else {
			log.Error().Msgf("Config values does not match bool type %s: %v", key, value)
		}
	}
	return defaultValue
}

func (c *Configuration) GetList(key ConfigKey) []interface{} {
	if value := c.getValue(key); value != nil {
		if t, ok := value.([]interface{}); ok {
			return t
		} else {
			log.Error().Msgf("Config values does not match list type %s: %v", key, value)
		}
	}
	return nil
}

func (c *Configuration) GetMap(key ConfigKey) map[interface{}]interface{} {
	if value := c.getValue(key); value != nil {
		if t, ok := value.(map[interface{}]interface{}); ok {
			return t
		} else {
			log.Error().Msgf("Config values does not match map type %s: %v", key, value)
		}
	}
	return nil
}
