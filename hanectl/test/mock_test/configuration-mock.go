package mock_test

import (
	"github.com/stretchr/testify/mock"
	"gohanectl/hanectl/config"
)

type ConfigurationMock struct {
	mock.Mock
}

func (m *ConfigurationMock) ReadConfiguration(configFile string) {
	m.Called(configFile)
}

func (m *ConfigurationMock) GetStr(key config.ConfigKey, defaultValue string) string {
	args := m.Called(key, defaultValue)
	return args.String(0)
}
func (m *ConfigurationMock) GetInt(key config.ConfigKey, defaultValue int) int {
	args := m.Called(key, defaultValue)
	return args.Int(0)
}
func (m *ConfigurationMock) GetBool(key config.ConfigKey, defaultValue bool) bool {
	args := m.Called(key, defaultValue)
	return args.Bool(0)
}
func (m *ConfigurationMock) GetList(key config.ConfigKey) []interface{} {
	args := m.Called(key)
	return args.Get(0).([]interface{})
}
func (m *ConfigurationMock) GetMap(key config.ConfigKey) map[interface{}]interface{} {
	args := m.Called(key)
	return args.Get(0).(map[interface{}]interface{})
}
