package model

type Notification struct {
	DeviceKey string   `json:"device_key" yaml:"device_key"`
	Keys      []string `json:"keys" yaml:"keys"`
	Caption   string   `json:"caption" yaml:"caption"`
	Script  string   `json:"script" yaml:"script"`
}

type Notifications struct {
	Notifications []Notification `json:"notifications" yaml:"notifications"`
}
