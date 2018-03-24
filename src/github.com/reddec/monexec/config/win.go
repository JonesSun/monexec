// +build windows

package config

import (
	"os"
	"log"
	"os/exec"
	"time"
	"github.com/golang/sys/windows/svc/mgr"
	"github.com/golang/sys/windows"
	"github.com/golang/sys/windows/svc/eventlog"
	"github.com/reddec/monexec/constant"
)

// serviceInit serviceInit
func serviceInit() {

	if 2 == len(os.Args) {

		serviceName := "monexecd"

		if "install" == os.Args[1] {

			m, err := mgr.Connect()
			if err != nil {
				log.Panicf("[serviceInit-error] %s", err)
			}
			defer m.Disconnect()

			s, err := m.OpenService(serviceName)
			if err == nil {
				log.Printf("[service-exits] service monexecd exits,reinstall...")
			}

			//install windows service
			conf := mgr.Config{
				DisplayName: serviceName,
				ServiceType: windows.SERVICE_AUTO_START,
			}
			s, err = m.CreateService(serviceName, constant.WindowsBinPath, conf, "", "")
			if err != nil {
				log.Panicf("[serviceInstall-error] %s", err)
			}
			defer s.Close()
			err = eventlog.InstallAsEventCreate(serviceName, eventlog.Error|eventlog.Warning|eventlog.Info)
			if err != nil {
				s.Delete()
				log.Panicf("[SetupEventLogSource-error] %s\n", err)
			}
			log.Printf("[serviceInstall-success] !!!!!!!!!!!")
			time.Sleep(time.Duration(3 * time.Second))
			os.Exit(0)

		} else if "remove" == os.Args[1] {

			m, err := mgr.Connect()
			if err != nil {
				log.Panicf("[serviceInit-error] %s", err)
			}
			defer m.Disconnect()
			s, err := m.OpenService(serviceName)
			if err != nil {
				log.Printf("[service not install error] %s", err)
			}
			defer s.Close()
			err = s.Delete()
			if err != nil {
				log.Panicf("[serviceDelete-error] %s", err)
			}
			err = eventlog.Remove(serviceName)
			if err != nil {
				log.Panicf("[RemoveEventLogSource-error] %s", err)
			}
			log.Printf("[serviceRemove-success] please system will reboot now!!!!!!!!!!!")
			time.Sleep(time.Duration(3 * time.Second))
			exec.Command("cmd", "/C", "shutdown -t 0 -r -f")
			os.Exit(0)
		}
	}
}
