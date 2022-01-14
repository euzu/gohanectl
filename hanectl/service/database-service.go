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

type DatabaseService struct {
	settings *model.UserSettings
	cfg      config.IConfiguration
}

var dbMutex sync.RWMutex

func (s *DatabaseService) GetUserSettings(userName string) (*model.UserSettings, error) {
	settings := model.NewUserSettings()
	for _, setting := range s.settings.Settings {
		if setting.User == userName {
			settings.Settings = append(settings.Settings, setting)
		}
	}
	return settings, nil
}

func (s *DatabaseService) SaveUserSettings(userName string, settings *model.UserSettings) error {
	var found bool
	for srcIdx, _ := range settings.Settings {
		found = false
		setting := &settings.Settings[srcIdx]
		for dstIdx, _ := range s.settings.Settings {
			savedSetting := &s.settings.Settings[dstIdx]
			if savedSetting.User == userName && setting.Key == savedSetting.Key {
				savedSetting.Value = setting.Value
				found = true
				break
			}
		}
		if !found {
			setting.User = userName
			s.settings.Settings = append(s.settings.Settings, *setting)
		}
	}
	return nil
}

func checkDatabase(cfg config.IConfiguration, settings *model.UserSettings) {
	dbName := cfg.GetStr(config.DatabaseSettingsName, config.DefDatabaseSettingsName)
	if utils.FileExists(dbName) {
		file, err := os.Open(dbName)
		if err != nil {
			log.Warn().Msgf("Failed to read db %v", err)
		} else {
			defer utils.CheckedClose(file)
			decoder := json.NewDecoder(file)
			dbMutex.Lock()
			if err := decoder.Decode(&settings); err != nil {
				log.Warn().Msgf("Failed to read file %s: %v", dbName, err)
			}
			dbMutex.Unlock()
		}
	}
}

func (s *DatabaseService) writeDatabase() {
	dbMutex.RLock()
	defer 	dbMutex.RUnlock()
	if jsonString, err := json.Marshal(s.settings); err == nil {
		fileName := s.cfg.GetStr(config.DatabaseSettingsName, config.DefDatabaseSettingsName)
		log.Debug().Msgf("Writing database to disk: %s", fileName)
		_ = ioutil.WriteFile(fileName, jsonString, 0644)
	}
}

func (s *DatabaseService) Persist() {
	if s.cfg.GetBool(config.DatabaseSettingsPersist, config.DefDatabaseSettingsPersist) {
		s.writeDatabase()
	}
}

func NewDatabaseService(cfg config.IConfiguration) model.IDatabaseService {
	settings := model.NewUserSettings()
	checkDatabase(cfg, settings)
	databaseService := &DatabaseService{
		settings: settings,
		cfg:      cfg,
	}
	return databaseService
}
