package service

import (
	"gohanectl/hanectl/config"
	"gohanectl/hanectl/model"
	"gohanectl/hanectl/repository"
)

type NotificationService struct {
	notificationRepo model.INotificationRepo
}

func (n *NotificationService) GetNotifications(deviceKey string, key string) ([]*model.Notification, error) {
	return n.notificationRepo.GetNotifications(deviceKey, key)
}

func (n *NotificationService) Close() {
	n.notificationRepo.Close()
}

func NewNotificationService(cfg config.IConfiguration) model.INotificationService {
	return &NotificationService{
		notificationRepo: repository.NewNotificationRepo(cfg),
	}
}
