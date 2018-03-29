package config

import (
	"github.com/magiconair/properties"
	"log"
	"runtime"
	"github.com/reddec/monexec/constant"
	"github.com/reddec/monexec/util"
	"io"
	"os"
	"gopkg.in/natefinch/lumberjack.v2"
)

// AgentConfig is global config
var AgentConfig config

type config struct {
	Services string `properties:"services,default="`
}

func Init() {
	initConfig()
	logInit()
}

//InitConfig
func initConfig() {

	var err error
	var path string
	if "windows" == runtime.GOOS {
		path = constant.WinConfPath
	} else {
		path = constant.LinuxConfPath
	}
	util.Config.Properties, err = properties.LoadFiles([]string{path}, properties.UTF8, true)
	if err != nil {
		panic(err)
	}
	properties.ErrorHandler = func(e error) {
		log.Printf("[setBoxConfig-error] %s\n", e)
	}
	if err = util.Config.Properties.Decode(&AgentConfig); err != nil {
		panic(err)
	}
}

//LogInit
func logInit() {

	lumberjackLogger := &lumberjack.Logger{
		MaxSize:    10, // megabytes
		MaxBackups: 5,
		MaxAge:     3, //days
		Compress:   true,
	}

	if "windows" == runtime.GOOS {

		lumberjackLogger.Filename = constant.WinLogPath

	} else {

		lumberjackLogger.Filename = constant.LinuxLogPath
	}
	util.LogOutput = io.MultiWriter(os.Stdout, lumberjackLogger)
	log.SetOutput(util.LogOutput)
}
