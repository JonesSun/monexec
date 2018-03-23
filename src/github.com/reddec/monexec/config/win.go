// +build windows

package config

import (
	"os"
	"log"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
	"github.com/golang/sys/windows/svc/mgr"
	"github.com/golang/sys/windows"
	"github.com/golang/sys/windows/svc/eventlog"
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

			file, _ := exec.LookPath(os.Args[0])
			path, _ := filepath.Abs(file)
			rst := strings.Replace(path, "\\", "/", -1)

			//install windows service
			conf := mgr.Config{
				DisplayName: serviceName,
				ServiceType: windows.SERVICE_AUTO_START}

			s, err = m.CreateService(serviceName, rst, conf, "", "")
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
