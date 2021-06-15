package mock_test

import (
	"github.com/stretchr/testify/mock"
	"gohanectl/hanectl/config"
)

type ConfigurationRepoMock struct {
	mock.Mock
}

func (m *ConfigurationRepoMock) ReadConfiguration(configFile string) {

}
func (m *ConfigurationRepoMock) GetStr(key config.ConfigKey, defaultValue string) string {
	args := m.Called(key, defaultValue)
	return args.String(0)
}
func (m *ConfigurationRepoMock) GetInt(key config.ConfigKey, defaultValue int) int {
	args := m.Called(key, defaultValue)
	return args.Int(0)
}
func (m *ConfigurationRepoMock) GetBool(key config.ConfigKey, defaultValue bool) bool {
	args := m.Called(key, defaultValue)
	return args.Bool(0)
}
func (m *ConfigurationRepoMock) GetList(key config.ConfigKey) []interface{} {
	args := m.Called(key)
	return args.Get(0).([]interface{})
}
func (m *ConfigurationRepoMock) GetMap(key config.ConfigKey) map[interface{}]interface{} {
	args := m.Called(key)
	return args.Get(0).(map[interface{}]interface{})
}
