package mock_test

import (
	"github.com/stretchr/testify/mock"
	"gohanectl/hanectl/model"
)

type UserServiceMock struct {
   mock.Mock
}

func (m *UserServiceMock) FindByUsername(userName string) (*model.User, error) {
	args := m.Called(userName)
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *UserServiceMock) ReloadUsers() error {
	args := m.Called()
	return args.Error(0)
}