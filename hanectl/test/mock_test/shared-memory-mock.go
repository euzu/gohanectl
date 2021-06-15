 package mock_test

 import (
	 "github.com/stretchr/testify/mock"
	 "gohanectl/hanectl/model"
 )


 type SharedMemoryMock struct {
	 mock.Mock
 }

 func (m *SharedMemoryMock) SetMem(deviceKey string, key string, value interface{}) {
	 m.Called(deviceKey, key, value)
 }
 func (m *SharedMemoryMock) MarkAsUpdated(deviceKey string) {
	 m.Called(deviceKey)
 }
 func (m *SharedMemoryMock) GetLastUpdated(deviceKey string) int64 {
	 args := m.Called(deviceKey)
	 return args.Get(0).(int64)
 }
 func (m *SharedMemoryMock) GetMem(deviceKey string, key string) interface{} {
	 args := m.Called(deviceKey, key)
	 return args.Get(0)
 }
 func (m *SharedMemoryMock) GetDeviceMem(deviceKey string) interface{} {
	 args := m.Called(deviceKey)
	 return args.Get(0)
 }
 func (m *SharedMemoryMock) GetMemory() model.Dictionary {
	 args := m.Called()
	 return args.Get(0).(model.Dictionary)
 }
 func (m *SharedMemoryMock) LoadSharedMem() {
	 m.Called()
 }

 func (m *SharedMemoryMock) SetNotifyCallback(notifyFunc model.NotifyFunc) {
	 m.Called(notifyFunc)
 }

 func (m *SharedMemoryMock) FinalizeSharedMem() {
	 m.Called()
 }
