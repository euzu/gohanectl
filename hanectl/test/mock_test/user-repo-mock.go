package mock_test

import (
	"github.com/stretchr/testify/mock"
	"gohanectl/hanectl/model"
)

type UserRepoMock struct {
	mock.Mock
}

func (m *UserRepoMock) FindByUsername(userName string) (*model.User, error) {
	args := m.Called(userName)
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *UserRepoMock) GetUsers() (*model.Users, error) {
	args := m.Called()
	return args.Get(0).(*model.Users), args.Error(1)
}

func (m *UserRepoMock) Close() {
	m.Called()
}