package main

import (
	"flag"
	"gohanectl/hanectl/app"
	"gohanectl/hanectl/config"
	"gohanectl/hanectl/utils"
)

func main() {
	configFile := flag.String("config", "", "Sets the config file to default")
	logLevel := flag.String("log", "", "Sets log level. Valid values are debug, info, warn, error, fatal, panic")
	flag.Parse()

	cfg := new(config.Configuration)
	cfg.ReadConfiguration(*configFile)

	if *logLevel == "" || logLevel == nil {
		*logLevel = cfg.GetStr(config.LogLevel, config.DefLogLevel)
	}

	utils.ChangeWorkingDir(cfg.GetStr(config.WorkingDir, ""))
	utils.InitLogger(logLevel, cfg.GetStr(config.LogFile, ""), cfg.GetBool(config.LogConsole, false))

	app.Boot(cfg)
}

