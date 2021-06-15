package model

type User struct {
	ID          int32    `json:"id" yaml:"id"`
	Username    string   `json:"username" yaml:"username"`
	Password    string   `json:"password" yaml:"password"`
	Authorities []string `json:"authorities" yaml:"authorities"`
	Enabled     bool     `json:"enabled" yaml:"enabled"`
}

type Users struct {
	Users []User `json:"users" yaml:"users"`
}

type UserSetting struct {
	User string `json:"user"`
	Key   string `json:"key"`
	Value string `json:"value"`
}

type UserSettings struct {
	Settings []UserSetting `json:"settings"`
}

func NewUserSettings() *UserSettings {
	return &UserSettings{Settings: make([]UserSetting, 0)}
}
