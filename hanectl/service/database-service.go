package service

//	"database/sql"
//	_ "github.com/mattn/go-sqlite3"

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
	if jsonString, err := json.Marshal(s.settings); err == nil {
		fileName := s.cfg.GetStr(config.DatabaseSettingsName, config.DefDatabaseSettingsName)
		log.Debug().Msgf("Writing database to disk: %s", fileName)
		_ = ioutil.WriteFile(fileName, jsonString, 0644)
	}
	dbMutex.RUnlock()
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

//const (
//	TblUserSettings = "user_settings"
//	ColUserName     = "uname"
//	ColKey          = "skey"
//	ColValue        = "sval"
//)

//func (s *DatabaseService) GetUserSettings(userName string) (*model.UserSettings, error) {
//	dbMutex.Lock()
//	defer dbMutex.Unlock()
//	dbName := s.cfg.GetStr(config.DatabaseSettingsName, config.DefDatabaseSettingsName)
//	db, err := sql.Open("sqlite3", dbName)
//	defer db.Close()
//	if err != nil {
//		return nil, err
//	}
//	sqlStmt := fmt.Sprintf("SELECT * FROM %s WHERE %s='%s'", TblUserSettings, ColUserName, userName)
//
//	settings := make([]model.UserSetting, 0)
//	rows, err := db.Query(sqlStmt)
//	if err != nil {
//		return nil, err
//	}
//
//	var uname string
//	for rows.Next() {
//		setting := model.UserSetting{}
//		err = rows.Scan(&uname, &setting.Key, &setting.Value)
//		if err != nil {
//			return nil, err
//		}
//		settings = append(settings, setting)
//	}
//
//    return &model.UserSettings{Settings: settings}, nil
//}
//
//func (s *DatabaseService) SaveUserSettings(userName string, settings *model.UserSettings) error {
//	dbMutex.Lock()
//	defer dbMutex.Unlock()
//	dbName := s.cfg.GetStr(config.DatabaseSettingsName, config.DefDatabaseSettingsName)
//	db, err := sql.Open("sqlite3", dbName)
//	defer db.Close()
//	if err != nil {
//		return err
//	}
//	sqlStmt := fmt.Sprintf("INSERT INTO %s (%s, %s, %s) VALUES (?,?,?) ON CONFLICT (%s, %s) DO UPDATE SET %s=? WHERE %s=? AND %s=?",
//		TblUserSettings, ColUserName, ColKey, ColValue, ColUserName, ColKey, ColValue, ColUserName, ColKey)
//
//	tx, err := db.Begin()
//	if err != nil {
//		return err
//	}
//
//	stmt, err := tx.Prepare(sqlStmt)
//	defer tx.Commit()
//
//	if err != nil {
//		return err
//	}
//	defer stmt.Close()
//
//	for i := range settings.Settings {
//		setting := settings.Settings[i]
//		_, err = stmt.Exec(userName, setting.Key, setting.Value, setting.Value, userName, setting.Key)
//		if err != nil {
//			tx.Rollback()
//			return err
//		}
//	}
//	return nil
//}
//
//func checkDatabase(cfg config.IConfiguration) {
//	dbName := cfg.GetStr(config.DatabaseSettingsName, config.DefDatabaseSettingsName)
//	db, err := sql.Open("sqlite3", dbName)
//	defer db.Close()
//	if err != nil {
//		log.Fatal().Msgf("Cant access database %v", err)
//	}
//	sqlStmt := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (%s TEXT not null, %s TEXT not null, %s TEXT, PRIMARY KEY (%s, %s));",
//		TblUserSettings, ColUserName, ColKey, ColValue, ColUserName, ColKey)
//	_, err = db.Exec(sqlStmt)
//	if err != nil {
//		log.Fatal().Msgf("Cant create database table %v", err)
//	}
//}
