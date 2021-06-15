package repository

import (
	"errors"
	"fmt"
	"github.com/rs/zerolog/log"
	"gohanectl/hanectl/config"
	"gohanectl/hanectl/model"
	"strings"
	"sync"
)

type UserRepo struct {
	users     *model.Users
	usersLock sync.RWMutex
	config    config.IConfiguration
}

func (u *UserRepo) loadUsers() (*model.Users, error) {
	users := &model.Users{}
	if _, err := readConfiguration(u.config, config.UserConfig, "", users); err != nil {
		log.Error().Msgf("Failed to read users file: %v", err)
		return nil, errors.New("failed to read users file")
	}
	u.usersLock.Lock()
	u.users = users
	u.usersLock.Unlock()
	return users, nil
}

func (u *UserRepo) ReloadUsers() error {
	if _, err := u.loadUsers(); err == nil {
		log.Info().Msg("Users reloaded")
		return nil
	} else {
		return err
	}
}

func (u *UserRepo) FindByUsername(userName string) (*model.User, error) {
	if users, err := u.GetUsers(); err == nil {
		for _, x := range users.Users {
			if strings.Compare(x.Username, userName) == 0 {
				return &x, nil
			}
		}
	}
	return nil, errors.New(fmt.Sprintf("failed to find user with username %s", userName))
}

func (u *UserRepo) GetUsers() (*model.Users, error) {
	if u.users != nil {
		u.usersLock.RLock()
		defer u.usersLock.RUnlock()
		return u.users, nil
	}
	return u.loadUsers()
}

func NewUserRepo(cfg config.IConfiguration) model.IUserRepo {
	return &UserRepo{
		config: cfg,
	}
}
