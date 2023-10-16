package main

import (
	config "github.com/satandyh/ansible-inventory-git-go/internal/config"
	//database "github.com/satandyh/ansible-inventory-git-go/internal/database"
	logging "github.com/satandyh/ansible-inventory-git-go/internal/logger"
)

// Global vars for logs
var logConfig = logging.LogConfig{
	ConsoleLoggingEnabled: true,
	EncodeLogsAsJson:      true,
	FileLoggingEnabled:    true,
	Directory:             "./data",
	Filename:              "ans-inv-git.log",
	MaxSize:               10,
	MaxBackups:            7,
	MaxAge:                7,
	LogLevel:              6,
}

var logger = logging.Configure(logConfig)

func main() {

	conf := config.NewConfig()
	println(conf.Nmap.Ip)   // test only
	println(conf.Nmap.Port) // test only

	logger.Info().
		Str("module", "main").
		Msg("All tasks completed.")

}
