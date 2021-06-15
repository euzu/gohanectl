package repository

import (
	"github.com/stretchr/testify/assert"
	"gohanectl/hanectl/config"
	"gohanectl/hanectl/test/mock_test"
	"testing"
)

var notificationsConfigDir = "../../config"
var notificationsFile = "../../config/notifications.yml"

func TestGetAllNotifications(t *testing.T) {

	cfg := new(mock_test.ConfigurationRepoMock)
	cfg.On("GetStr", config.ConfigDirectory, config.DefConfigDirectory).Return(notificationsConfigDir)
	cfg.On("GetStr", config.NotificationConfig, "").Return(notificationsFile)

	repo := NewNotificationRepo(cfg)
	notifications, err := repo.GetAllNotifications()
	assert.Nil(t, err)
	assert.True(t, len(notifications.Notifications) > 0)

	cfg.AssertExpectations(t)
}

func TestReloadNotifications(t *testing.T) {

	cfg := new(mock_test.ConfigurationRepoMock)
	cfg.On("GetStr", config.ConfigDirectory, config.DefConfigDirectory).Return(notificationsConfigDir)
	cfg.On("GetStr", config.NotificationConfig, "").Return(notificationsFile)

	repo := NewNotificationRepo(cfg)
	err := repo.ReloadNotifications()
	assert.Nil(t, err)

	cfg.AssertExpectations(t)
}

func TestGetNotifications(t *testing.T) {

	cfg := new(mock_test.ConfigurationRepoMock)
	cfg.On("GetStr", config.ConfigDirectory, config.DefConfigDirectory).Return(notificationsConfigDir)
	cfg.On("GetStr", config.NotificationConfig, "").Return(notificationsFile)

	repo := NewNotificationRepo(cfg)
	notifications, err := repo.GetNotifications("socket-washer", "power")
	assert.Nil(t, err)
	assert.NotNil(t, notifications)
	assert.True(t, len(notifications) > 0)

	assert.Equal(t, notifications[0].Caption, "Waschmaschine")
	assert.Equal(t, notifications[0].Template, "washer_power")


	notifications, err = repo.GetNotifications("socket-washer", "unknown")
	assert.Errorf(t, err, "no notification found")
	assert.Zero(t, len(notifications))

	cfg.AssertExpectations(t)
}