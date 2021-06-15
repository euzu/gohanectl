package service

import (
	"gohanectl/hanectl/config"
	"gohanectl/hanectl/model"
	"gohanectl/hanectl/repository"
)

type NotificationService struct {
	mqttService      *MqttService
	notifications    *model.Notifications
	notificationRepo model.INotificationRepo
}

func (n *NotificationService) GetNotifications(deviceKey string, key string) ([]*model.Notification, error) {
	return n.notificationRepo.GetNotifications(deviceKey, key)
}

func (n *NotificationService) ReloadNotifications() error {
	return n.notificationRepo.ReloadNotifications()
}

func NewNotificationService(cfg config.IConfiguration) model.INotificationService {
	return &NotificationService{
		notificationRepo: repository.NewNotificationRepo(cfg),
	}
}
