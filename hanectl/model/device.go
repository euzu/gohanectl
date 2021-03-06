package model

import (
	"github.com/markphelps/optional"
	_ "gopkg.in/yaml.v2"
	"strings"
)

type OptionalBool struct {
	optional.Bool
}

func (b *OptionalBool) UnmarshalYAML(unmarshal func(interface{}) error) error {
	value := false
	if err := unmarshal(&value); err != nil {
		return err
	}
	b.Bool.Set(value)
	return nil
}

func (b *OptionalBool) MarshalYAML() (interface{}, error) {
	return b.OrElse(false), nil
}

func NewBool(v bool) OptionalBool {
	ob := optional.NewBool(v)
	return OptionalBool{ob}
}

type DeviceTopic struct {
	topic           string
	Topic           string `json:"topic" yaml:"topic"`
	HandlerTemplate string `json:"handler_template" yaml:"handler_template"`
}

func (t *DeviceTopic) Prepare() {
	if strings.HasSuffix(t.Topic, "/#") {
		t.topic = t.Topic[0 : len(t.Topic)-2]
	} else if strings.HasSuffix(t.Topic, "/") {
		t.topic = t.Topic[0 : len(t.Topic)-1]
	} else {
		t.topic = t.Topic
	}
}

func (t *DeviceTopic) Matches(topic string) bool {
	return len(t.topic) > 0 && strings.HasPrefix(topic, t.topic)
}

type DeviceMqttTemplate struct {
	ListenTopics  []DeviceTopic `json:"listen_topics" yaml:"listen_topics"`
	CommandTopics Dictionary    `json:"command_topics" yaml:"command_topics"`
}

type DeviceMqtt struct {
	DeviceID      string        `json:"device_id" yaml:"device_id"`
	ListenTopics  []DeviceTopic `json:"listen_topics" yaml:"listen_topics"`
	CommandTopics Dictionary    `json:"command_topics" yaml:"command_topics"`
}

type DeviceRestTemplate struct {
	CommandPaths    Dictionary `json:"command_paths" yaml:"command_paths"`
	HandlerTemplate string     `json:"handler_template" yaml:"handler_template"`
}

type DeviceRest struct {
	Url             string     `json:"url" yaml:"url"`
	CommandPaths    Dictionary `json:"command_paths" yaml:"command_paths"`
	HandlerTemplate string     `json:"handler_template" yaml:"handler_template"`
}

type Supplemental struct {
	Field    string `json:"field" yaml:"field"`
	Caption  string `json:"caption" yaml:"caption"`
	Format   string `json:"format" yaml:"format"`
	Renderer string `json:"renderer" yaml:"renderer"`
}

type DeviceTemplate struct {
	Mqtt         DeviceMqttTemplate  `json:"mqtt" yaml:"mqtt"`
	Rest         DeviceRestTemplate  `json:"rest" yaml:"rest"`
	Optimistic   bool                `json:"optimistic" yaml:"optimistic"`
	Timeout      int64               `json:"timeout" yaml:"timeout"`
	Icon         string              `json:"icon" yaml:"icon"`
	Supplemental []Supplemental      `json:"supplemental" yaml:"supplemental"`
}

type Device struct {
	Type         string         `json:"type" yaml:"type"`
	DeviceKey    string         `json:"device_key" yaml:"device_key"`
	Caption      string         `json:"caption" yaml:"caption"`
	Confirm      bool           `json:"confirm" yaml:"confirm"`
	Template     string         `json:"template" yaml:"template"`
	Mqtt         DeviceMqtt     `json:"mqtt" yaml:"mqtt"`
	Rest         DeviceRest     `json:"rest" yaml:"rest"`
	Optimistic   OptionalBool   `json:"optimistic" yaml:"optimistic"`
	Timeout      int64          `json:"timeout" yaml:"timeout"`
	Invert       Dictionary     `json:"invert" yaml:"invert"`
	Room         string         `json:"room" yaml:"room"`
	Supplemental []Supplemental `json:"supplemental" yaml:"supplemental"`
	Groups       []string       `json:"groups" yaml:"groups"`
	Authorities  []string       `json:"authorities" yaml:"authorities"`
	Expanded     bool           `json:"expanded" yaml:"expanded"`
	Icon         string         `json:"icon" yaml:"icon"`
}

func (d *Device) hasMqtt() bool {
	return len(d.Mqtt.ListenTopics) > 0
}

type Devices struct {
	Devices []Device `json:"devices" yaml:"devices"`
}

type DeviceDto struct {
	Type         string         `json:"type"`
	DeviceKey    string         `json:"deviceKey"`
	Caption      string         `json:"caption"`
	Confirm      bool           `json:"confirm"`
	Optimistic   bool           `json:"optimistic"`
	Url          string         `json:"url"`
	Timeout      int64          `json:"timeout"`
	Invert       Dictionary     `json:"invert"`
	Room         string         `json:"room"`
	Groups       []string       `json:"groups"`
	Supplemental []Supplemental `json:"supplemental"`
	Expanded     bool           `json:"expanded"`
	Icon         string         `json:"icon"`
}

type DevicesDto struct {
	Devices []DeviceDto `json:"devices"`
}

type DeviceFilter func(device *Device) bool
