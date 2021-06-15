package repository

import (
	"github.com/stretchr/testify/assert"
	"gohanectl/hanectl/config"
	"gohanectl/hanectl/test/mock_test"
	"testing"
)

var usersConfigDir = "../../config"
var usersFile = "../../config/users.yml"

func TestGetAllUsers(t *testing.T) {

	cfg := new(mock_test.ConfigurationRepoMock)
	cfg.On("GetStr", config.ConfigDirectory, config.DefConfigDirectory).Return(usersConfigDir)
	cfg.On("GetStr", config.UserConfig, "").Return(usersFile)

	repo := NewUserRepo(cfg)
	users, err := repo.GetUsers()
	assert.Nil(t, err)
	assert.True(t, len(users.Users) > 0)

	cfg.AssertExpectations(t)
}

func TestReloadUsers(t *testing.T) {

	cfg := new(mock_test.ConfigurationRepoMock)
	cfg.On("GetStr", config.ConfigDirectory, config.DefConfigDirectory).Return(usersConfigDir)
	cfg.On("GetStr", config.UserConfig, "").Return(usersFile)

	repo := NewUserRepo(cfg)
	err := repo.ReloadUsers()
	assert.Nil(t, err)

	cfg.AssertExpectations(t)
}

func TestGetUsers(t *testing.T) {

	cfg := new(mock_test.ConfigurationRepoMock)
	cfg.On("GetStr", config.ConfigDirectory, config.DefConfigDirectory).Return(usersConfigDir)
	cfg.On("GetStr", config.UserConfig, "").Return(usersFile)

	repo := NewUserRepo(cfg)
	user, err := repo.FindByUsername("sarcon")
	assert.Nil(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, user.Username, "sarcon")

	user, err = repo.FindByUsername("unknown")
	assert.Errorf(t, err, "failed to find user with username unknown")
	assert.Nil(t, user)

	cfg.AssertExpectations(t)
}