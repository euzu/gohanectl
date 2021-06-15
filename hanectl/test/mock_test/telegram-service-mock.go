package mock_test

import "github.com/stretchr/testify/mock"

type TelegramServiceMock struct {
	mock.Mock
}

func (m *TelegramServiceMock) SendMessage(message string) bool {
	args := m.Called(message)
	return args.Bool(0)
}