package service

import (
	"gohanectl/hanectl/config"
	"gohanectl/hanectl/model"
	"sync/atomic"
	"time"
)

var (
	finalizeLocker uint32
)

type ServiceFactory struct {
	configService       model.IConfigService
	telegramService     model.ITelegramService
	restService         model.IRestService
	mqttService         model.IMqttService
	deviceService       model.IDeviceService
	userService         model.IUserService
	notificationService model.INotificationService
	sharedMemory        model.ISharedMemory
	databaseService     model.IDatabaseService
}

func (s *ServiceFactory) GetDeviceService() model.IDeviceService {
	return s.deviceService
}

func (s *ServiceFactory) GetConfigService() model.IConfigService {
	return s.configService
}

func (s *ServiceFactory) GetMqttService() model.IMqttService {
	return s.mqttService
}

func (s *ServiceFactory) GetTelegramService() model.ITelegramService {
	return s.telegramService
}

func (s *ServiceFactory) GetUserService() model.IUserService {
	return s.userService
}

func (s *ServiceFactory) GetNotificationService() model.INotificationService {
	return s.notificationService
}

func (s *ServiceFactory) GetRestService() model.IRestService {
	return s.restService
}

func (s *ServiceFactory) GetDatabaseService() model.IDatabaseService {
	return s.databaseService
}

func (s *ServiceFactory) GetSharedMemory() model.ISharedMemory {
	return s.sharedMemory
}

func (s *ServiceFactory) Finalize() {
	if !atomic.CompareAndSwapUint32(&finalizeLocker, 0, 1) {
		return
	}
	defer atomic.StoreUint32(&finalizeLocker, 0)

	s.sharedMemory.Persist()
	s.databaseService.Persist()
	s.deviceService.Close()
	s.notificationService.Close()
	s.userService.Close()
	s.mqttService.Close()
}

func (s *ServiceFactory) runPersistor(cfg config.IConfiguration) {
	if cfg.GetBool(config.DatabaseStatesPersist, config.DefDatabaseStatesPersist) || cfg.GetBool(config.DatabaseSettingsPersist, config.DefDatabaseSettingsPersist) {
		ticker := time.NewTicker(15 * time.Minute)
		quit := make(chan struct{})
		go func() {
			for {
				select {
				case <-ticker.C:
					s.Finalize()
				case <-quit:
					ticker.Stop()
					return
				}
			}
		}()
	}
}

func NewServiceFactory(cfg config.IConfiguration) model.IServiceFactory {
	sharedMemory := NewSharedMemory(cfg)
	restService := NewRestService()
	mqttService := NewMqttService()
	databaseService := NewDatabaseService(cfg)
	return &ServiceFactory{
		databaseService:     databaseService,
		configService:       NewConfigService(cfg),
		telegramService:     NewTelegramService(cfg),
		sharedMemory:        sharedMemory,
		restService:         restService,
		mqttService:         mqttService,
		deviceService:       NewDeviceService(cfg, restService, mqttService, sharedMemory),
		userService:         NewUserService(cfg, databaseService),
		notificationService: NewNotificationService(cfg),
	}

}
