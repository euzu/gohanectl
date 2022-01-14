package model

type IDeviceRepo interface {
	Close()
	GetDevice(deviceKey string) (*Device, error)
	GetDevices() (*Devices, error)
}

type INotificationRepo interface {
	Close()
	GetNotifications(deviceKey string, key string) ([]*Notification, error)
	GetAllNotifications() (*Notifications, error)
}

type IUserRepo interface {
	Close()
	FindByUsername(userName string) (*User, error)
	GetUsers() (*Users, error)
}

