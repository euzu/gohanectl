package app

import (
	"github.com/rs/zerolog/log"
	"gohanectl/hanectl/config"
	"gohanectl/hanectl/model"
	"gohanectl/hanectl/service"
	"gohanectl/hanectl/utils"
	jscript "gohanectl/hanectl/vm"
)

func Boot(cfg config.IConfiguration) {
	dropPrivileges(cfg.GetStr(config.RunUser, ""))
	serviceFactory := service.NewServiceFactory(cfg)
	utils.CatchOsSignals(func() {
		serviceFactory.Finalize()
	}, reloadCallback(serviceFactory))

	go initDevices(cfg, serviceFactory)

	jscript.InitVM(cfg, serviceFactory)
	startServer(cfg, serviceFactory)
}

func reloadCallback(serviceFactory model.IServiceFactory) func() {
	return func() {
		//TODO Reloading the configurations currently has no effect on MQTT subscriptions.
		//  The devices are loaded and cached initially.
		//if err := serviceFactory.GetDeviceService().ReloadDevices(); err != nil {
		//	log.Err(err)
		//}
		if err := serviceFactory.GetNotificationService().ReloadNotifications(); err != nil {
			log.Err(err)
		}
		if err := jscript.ReloadScripts(); err != nil {
			log.Err(err)
		}
		if err := serviceFactory.GetUserService().ReloadUsers(); err != nil {
			log.Err(err)
		}
	}
}
