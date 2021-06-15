package service

import (
	"gohanectl/hanectl/config"
	"gohanectl/hanectl/model"
	"gohanectl/hanectl/repository"
)

type UserService struct {
	users           *model.Users
	userRepo        model.IUserRepo
	databaseService model.IDatabaseService
}

func (u *UserService) FindByUsername(userName string) (*model.User, error) {
	return u.userRepo.FindByUsername(userName)
}

func (u *UserService) ReloadUsers() error {
	return u.userRepo.ReloadUsers()
}

func (u *UserService) SaveSettings(userName string, settings *model.UserSettings) error {
	return u.databaseService.SaveUserSettings(userName, settings)
}

func (u *UserService) GetSettings(userName string) (*model.UserSettings, error) {
	return u.databaseService.GetUserSettings(userName)
}

func NewUserService(cfg config.IConfiguration, databaseService model.IDatabaseService) model.IUserService {
	return &UserService{
		userRepo:        repository.NewUserRepo(cfg),
		databaseService: databaseService,
	}
}
