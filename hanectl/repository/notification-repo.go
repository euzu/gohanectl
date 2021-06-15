package repository

import (
	"errors"
	"github.com/rs/zerolog/log"
	"gohanectl/hanectl/config"
	"gohanectl/hanectl/model"
	"gohanectl/hanectl/utils"
	"strings"
	"sync"
)

type NotificationRepo struct {
	notifications     *model.Notifications
	notificationsLock sync.RWMutex
	config            config.IConfiguration
}

func (n *NotificationRepo) loadNotifications() (*model.Notifications, error) {
	notifications := &model.Notifications{}
	if _, err := readConfiguration(n.config, config.NotificationConfig, "", notifications); err != nil {
		log.Error().Msgf("Failed to read notifications file: %v", err)
		return nil, errors.New("failed to read notifications file")
	}
	n.notificationsLock.Lock()
	n.notifications = notifications
	n.notificationsLock.Unlock()
	return notifications, nil
}

func (n *NotificationRepo) ReloadNotifications() error {
	if _, err := n.loadNotifications(); err == nil {
		log.Info().Msg("Notifications reloaded")
		return nil
	} else {
		return err
	}
}

func (n *NotificationRepo) GetNotifications(deviceKey string, key string) ([]*model.Notification, error) {
	if notifications, err := n.GetAllNotifications(); err == nil {
		var result []*model.Notification
		for i, x := range notifications.Notifications {
			if strings.Compare(x.DeviceKey, deviceKey) == 0 {
				if len(x.Keys) == 0 || utils.ContainsStr(x.Keys, key) {
					result = append(result, &notifications.Notifications[i])
				}
			}
		}
		if result == nil {
			return nil, errors.New("no notification found")
		}
		return result, nil
	}
	return nil, errors.New("cant find notification")
}

func (n *NotificationRepo) GetAllNotifications() (*model.Notifications, error) {
	if n.notifications != nil {
		n.notificationsLock.RLock()
		defer n.notificationsLock.RUnlock()
		return n.notifications, nil
	}
	return n.loadNotifications()
}

func NewNotificationRepo(cfg config.IConfiguration) *NotificationRepo {
	return &NotificationRepo{
		config: cfg,
	}
}
