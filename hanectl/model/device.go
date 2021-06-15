package model

import (
	"strings"
)

type DeviceMqttTemplate struct {
	Name     string `json:"name" yaml:"name"`
	DeviceID string `json:"device_id" yaml:"device_id"`
}

type DeviceRestTemplate struct {
	Name     string `json:"name" yaml:"name"`
	DeviceID string `json:"device_id" yaml:"device_id"`
	Url      string `json:"url" yaml:"url"`
	Template string `json:"template" yaml:"template"`
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

type DeviceMqtt struct {
	Template      DeviceMqttTemplate `json:"template" yaml:"template"`
	ListenTopics  []DeviceTopic      `json:"listen_topics" yaml:"listen_topics"`
	CommandTopics Dictionary         `json:"command_topics" yaml:"command_topics"`
}

type DeviceRest struct {
	Url             string             `json:"url" yaml:"url"`
	Template        DeviceRestTemplate `json:"template" yaml:"template"`
	CommandPaths    Dictionary         `json:"command_paths" yaml:"command_paths"`
	HandlerTemplate string             `json:"handler_template" yaml:"handler_template"`
}

type Supplemental struct {
	Field    string `json:"field" yaml:"field"`
	Caption  string `json:"caption" yaml:"caption"`
	Format   string `json:"format" yaml:"format"`
	Renderer string `json:"renderer" yaml:"renderer"`
}

type Device struct {
	Type         string         `json:"type" yaml:"type"`
	DeviceKey    string         `json:"device_key" yaml:"device_key"`
	Caption      string         `json:"caption" yaml:"caption"`
	Confirm      bool           `json:"confirm" yaml:"confirm"`
	Mqtt         DeviceMqtt     `json:"mqtt" yaml:"mqtt"`
	Rest         DeviceRest     `json:"rest" yaml:"rest"`
	Optimistic   bool           `json:"optimistic" yaml:"optimistic"`
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
