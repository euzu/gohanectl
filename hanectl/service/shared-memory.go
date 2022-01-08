package service

import (
	"encoding/json"
	"github.com/rs/zerolog/log"
	"gohanectl/hanectl/config"
	"gohanectl/hanectl/model"
	"gohanectl/hanectl/utils"
	"io/ioutil"
	"os"
	"sync"
)

const KeyLastUpdated = "lastUpdated"

var notificationCallback model.NotifyFunc

var mutex sync.RWMutex

type SharedMemory struct {
	sharedMem model.Dictionary
	config    config.IConfiguration
}

func triggerValueChange(deviceKey string, key string, newValue interface{}, oldValue interface{}) {
	log.Debug().Msgf("Triggered %s.%s: newValue=%v, oldValue=%v", deviceKey, key, newValue, oldValue)
	if notificationCallback != nil {
		notificationCallback(deviceKey, key, newValue, oldValue)
	}
}

func (s *SharedMemory) readSharedMem() {
	fileName := s.config.GetStr(config.DatabaseStatesName, config.DefDatabaseStatesName)
	if utils.FileExists(fileName) {
		file, err := os.Open(fileName)
		if err != nil {
			log.Warn().Msgf("Failed to read configuration %v", err)
		} else {
			defer utils.CheckedClose(file)
			decoder := json.NewDecoder(file)
			mutex.Lock()
			if err := decoder.Decode(&s.sharedMem); err != nil {
				log.Warn().Msgf("Failed to read file %s: %v", fileName, err)
			}
			mutex.Unlock()
		}
	}
}

func (s *SharedMemory) writeSharedMem() {
	mutex.RLock()
	if jsonString, err := json.Marshal(s.sharedMem); err == nil {
		fileName := s.config.GetStr(config.DatabaseStatesName, config.DefDatabaseStatesName)
		log.Debug().Msgf("Writing sharedMem to disk: %s", fileName)
		_ = ioutil.WriteFile(fileName, jsonString, 0644)
	}
	mutex.RUnlock()
}

func (s *SharedMemory) GetMemory() model.Dictionary {
	return s.sharedMem
}

func (s *SharedMemory) setMem(deviceKey string, key string, value interface{}, trigger bool) {
	var mem model.Dictionary
	mutex.Lock()
	defer mutex.Unlock()
	if value, exist := s.sharedMem[deviceKey]; exist {
		mem = value.(model.Dictionary)
	} else {
		mem = make(model.Dictionary)
		s.sharedMem[deviceKey] = mem
	}
	if trigger {
		var oldValue, _ = mem[key]
		//if oldValue, exists := mem[key]; exists {
		go triggerValueChange(deviceKey, key, value, oldValue)
		//}
	}
	mem[key] = value
}

func (s *SharedMemory) SetMem(deviceKey string, key string, value interface{}) {
	s.setMem(deviceKey, key, value, true)
}

func (s *SharedMemory) MarkAsUpdated(deviceKey string) {
	s.setMem(deviceKey, KeyLastUpdated, utils.NowTimestamp(), false)
}

func (s *SharedMemory) GetLastUpdated(deviceKey string) int64 {
	lastUpdated := s.GetMem(deviceKey, KeyLastUpdated)
	if lastUpdated != nil {
		if t, ok := lastUpdated.(int); ok {
			return int64(t)
		} else if t, ok := lastUpdated.(float64); ok {
			return int64(t)
		} else if t, ok := lastUpdated.(int64); ok {
			return t
		}
	}
	return 0
}

func (s *SharedMemory) GetMem(deviceKey string, key string) interface{} {
	mutex.RLock()
	defer mutex.RUnlock()
	if value, exist := s.sharedMem[deviceKey]; exist {
		mem := value.(model.Dictionary)
		if value, exist = mem[key]; exist {
			return value
		}
	}
	return nil
}

func (s *SharedMemory) GetDeviceMem(deviceKey string) interface{} {
	mutex.RLock()
	defer mutex.RUnlock()
	if value, exist := s.sharedMem[deviceKey]; exist {
		return value
	}
	return nil
}

func (s *SharedMemory) LoadSharedMem() {
	s.readSharedMem()
}

func (s *SharedMemory) SetNotifyCallback(notifyFunc model.NotifyFunc) {
	notificationCallback = notifyFunc
}

func (s *SharedMemory) Persist() {
	if s.config.GetBool(config.DatabaseStatesPersist, config.DefDatabaseStatesPersist) {
		s.writeSharedMem()
	}
}

func NewSharedMemory(cfg config.IConfiguration) model.ISharedMemory {
	sharedMemory := &SharedMemory{
		sharedMem: make(model.Dictionary),
		config:    cfg,
	}
	return sharedMemory
}
