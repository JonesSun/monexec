package main

import (
	"os"
	"github.com/reddec/container"
	"context"
	"log"
	"github.com/reddec/monexec/monexec"
	conf "github.com/reddec/monexec/config"
	"os/signal"
	"syscall"
	"runtime"
	"github.com/reddec/monexec/util"
	"strings"
	"path"
)

func run() {

	conf.ServiceInit()
	conf.Init()
	config := DefaultConfig()

	agentServer := "/home/i5/bin/agentServer"
	if "windows" == runtime.GOOS {
		agentServer = "c:/bin/agentServer.exe"
	}
	config.Services = append(config.Services, monexec.Executable{
		Name:    "agentServer",
		Command: agentServer,
		Restart: 0,
	})

	conf.Config.Services = strings.TrimSpace(conf.Config.Services)
	if "" != conf.Config.Services {

		services := strings.Split(conf.Config.Services, ",")
		for _, service := range services {

			fileName := path.Base(service)
			config.Services = append(config.Services, monexec.Executable{
				Name:    fileName,
				Command: service,
				Restart: 0,
			})
		}
	}

	sv := container.NewSupervisor(log.New(util.LogOutput, "[supervisor] ", log.LstdFlags))
	runConfigInSupervisor(&config, sv)
}

func runConfigInSupervisor(config *Config, sv container.Supervisor) {

	ctx, stop := context.WithCancel(context.Background())
	c := make(chan os.Signal, 2)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		for range c {
			stop()
			break
		}
	}()

	err := config.Run(sv, ctx)

	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	run()
}
