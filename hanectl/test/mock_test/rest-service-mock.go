package mock_test

import (
	"github.com/stretchr/testify/mock"
	"gohanectl/hanectl/model"
)

type RestServiceMock struct {
	mock.Mock
}

func (m *RestServiceMock) GetRequest(url string, device *model.Device) bool {
	args := m.Called(url, device)
	return args.Bool(0)
}

func (m *RestServiceMock) SetMessageHandler(handler model.RestMessageHandler) {
	m.Called(handler)
}
func (m *RestServiceMock) HandleMessage(device *model.Device, payload string) {
	m.Called(device, payload)
}
