package main

import (
	"os"
	"github.com/reddec/container"
	"context"
	"log"
	"github.com/reddec/monexec/monexec"
	"os/signal"
	"syscall"
	"runtime"
)

func run() {

	path := "/home/i5/bin/agentServer"
	if "windows" == runtime.GOOS {
		path = "c:/bin/agentServer.exe"
	}

	config := DefaultConfig()

	config.Services = append(config.Services, monexec.Executable{
		Name:    "agentServer",
		Command: path,
		Restart: 0,
	})
	FillDefaultExecutable(&config.Services[0])

	sv := container.NewSupervisor(log.New(os.Stderr, "[supervisor] ", log.LstdFlags))
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
