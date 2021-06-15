package model

type IDeviceRepo interface {
	ReloadDevices() error
	GetDevice(deviceKey string) (*Device, error)
	GetDevices() (*Devices, error)
}

type INotificationRepo interface {
	ReloadNotifications() error
	GetNotifications(deviceKey string, key string) ([]*Notification, error)
	GetAllNotifications() (*Notifications, error)
}

type IUserRepo interface {
	ReloadUsers() error
	FindByUsername(userName string) (*User, error)
	GetUsers() (*Users, error)
}

