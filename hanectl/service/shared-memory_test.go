package service

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"gohanectl/hanectl/config"
	"gohanectl/hanectl/model"
	"gohanectl/hanectl/test/mock_test"
	"gohanectl/hanectl/utils"
	"testing"
)

func TestSetMem(t *testing.T) {
	smem := SharedMemory{
		sharedMem: make(model.Dictionary),
	}
	smem.SetMem("deviceKey", "key", "value")
	assert.Equal(t, smem.GetMem("deviceKey", "key"), "value")
}

func TestMarkAsUpdated(t *testing.T) {
	smem := SharedMemory{
		sharedMem: make(model.Dictionary),
	}
	smem.MarkAsUpdated("deviceKey")
	assert.NotNil(t, smem.GetMem("deviceKey", KeyLastUpdated))
	assert.Less(t, utils.NowTimestamp()-smem.GetMem("deviceKey", KeyLastUpdated).(int64), int64(100))
}

func TestGetLastUpdated(t *testing.T) {
	smem := SharedMemory{
		sharedMem: make(model.Dictionary),
	}

	value := smem.GetLastUpdated("deviceKey")
	assert.Equal(t, value, int64(0))

	smem.MarkAsUpdated("deviceKey")
	value = smem.GetLastUpdated("deviceKey")
	assert.Less(t, utils.NowTimestamp()-value, int64(100))

}

func TestGetMem(t *testing.T) {
	smem := SharedMemory{
		sharedMem: make(model.Dictionary),
	}

	value := smem.GetMem("deviceKey", "key")
	assert.Nil(t, value)

	smem.SetMem("deviceKey", "key", "data")
	value = smem.GetMem("deviceKey", "key")

	assert.Equal(t, value, "data")
}

func TestGetDeviceMem(t *testing.T) {
	smem := SharedMemory{
		sharedMem: make(model.Dictionary),
	}

	value := smem.GetDeviceMem("deviceKey")
	assert.Nil(t, value)

	smem.SetMem("deviceKey", "key", "data")
	value = smem.GetDeviceMem("deviceKey")

	v1, _ := json.Marshal(model.Dictionary{"key": "data"})
	v2, _ := json.Marshal(value)
	assert.Equal(t, string(v1), string(v2))
}

func TestGetMemory(t *testing.T) {
	smem := SharedMemory{
		sharedMem: make(model.Dictionary),
	}

	value := smem.GetMemory()
	assert.Empty(t, value)

	smem.SetMem("deviceKey", "key", "data")
	value = smem.GetMemory()

	v1, _ := json.Marshal(model.Dictionary{ "deviceKey" : model.Dictionary{"key": "data"}})
	v2, _ := json.Marshal(value)
	assert.Equal(t, string(v1), string(v2))

}

func TestLoadSharedMem(t *testing.T) {
	cfgMock := new (mock_test.ConfigurationMock)
	smem := SharedMemory{
		config: cfgMock,
		sharedMem: make(model.Dictionary),
	}
	cfgMock.On("GetStr", config.DatabaseName, config.DefDatabaseName).Return("")
	smem.LoadSharedMem()
	cfgMock.AssertExpectations(t)
}

func TestFinalizeSharedMem(t *testing.T) {
	cfgMock := new (mock_test.ConfigurationMock)
	smem := SharedMemory{
		config: cfgMock,
		sharedMem: make(model.Dictionary),
	}
	cfgMock.On("GetBool", config.DatabasePersist, config.DefDatabasePersist).Return(false)
	smem.FinalizeSharedMem()
	cfgMock.AssertExpectations(t)
}
