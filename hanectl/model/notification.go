package model

type NotificationTrigger = string

const (
	TriggerChange NotificationTrigger = "change"
	TriggerEvent  NotificationTrigger = "event"
)

type Notification struct {
	DeviceKey string   `json:"device_key" yaml:"device_key"`
	Trigger   NotificationTrigger `json:"trigger" yaml:"trigger"`
	Keys      []string `json:"keys" yaml:"keys"`
	Caption   string   `json:"caption" yaml:"caption"`
	Script    string   `json:"script" yaml:"script"`
}

type Notifications struct {
	Notifications []Notification `json:"notifications" yaml:"notifications"`
}
