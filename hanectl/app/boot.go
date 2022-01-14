package app

import (
	"flag"
	"github.com/rs/zerolog/log"
	"gohanectl/hanectl/config"
	"gohanectl/hanectl/model"
	"gohanectl/hanectl/service"
	"gohanectl/hanectl/utils"
	jscript "gohanectl/hanectl/vm"
	"net"
	"sync/atomic"
)

var (
	realoadLocker uint32
	httpListener *net.Listener
	configFile *string
	logLevel *string
)

func Boot() {

	configFile = flag.String("config", "", "Sets the config file to default")
	logLevel = flag.String("log", "", "Sets log level. Valid values are debug, info, warn, error, fatal, panic")
	flag.Parse()

	loadConfig(true)
}

func loadConfig(startHttpServer bool) {
	cfg := new(config.Configuration)
	cfg.ReadConfiguration(*configFile)

	if *logLevel == "" || logLevel == nil {
		*logLevel = cfg.GetStr(config.LogLevel, config.DefLogLevel)
	}

	utils.ChangeWorkingDir(cfg.GetStr(config.WorkingDir, ""))
	utils.InitLogger(logLevel, cfg.GetStr(config.LogFile, ""), cfg.GetBool(config.LogConsole, false))

	dropPrivileges(cfg.GetStr(config.RunUser, ""))
	serviceFactory := service.NewServiceFactory(cfg)
	utils.CatchOsSignals(func() {
		serviceFactory.Finalize()
	}, reloadCallback(serviceFactory))

	initDevices(cfg, serviceFactory)
	jscript.InitVM(cfg, serviceFactory)
	if startHttpServer {
		httpListener = startServer(cfg, serviceFactory)
	}
}

func reloadCallback(serviceFactory model.IServiceFactory) func() {
	return func() {
		if !atomic.CompareAndSwapUint32(&realoadLocker, 0, 1) {
			return
		}
		defer atomic.StoreUint32(&realoadLocker, 0)
		log.Info().Msg("Reloading config")
		// due to problems with socket release and binding, disabled
		//if httpListener !=  nil {
		//	if err := (*httpListener).Close(); err != nil {
		//		log.Err(err)
		//	}
		//}
		serviceFactory.Finalize()
		jscript.Close()
		loadConfig(false)
	}
}
