// +build linux

package config

import (
	"os"
	"log"
	"time"
	"github.com/reddec/monexec/constant"
	"fmt"
	"os/exec"
)

const shell = `
#!/bin/sh
# monexec

start() {
        start-stop-daemon --start -b --quiet --oknodo --pidfile /var/run/monexec.pid --exec ` + constant.LinuxBinPath + `
}

stop() {
        ps aux|grep monexec|grep grep -v|awk '{print $2}'|xargs kill -9
}

restart() {
    stop
    start
}

case "$1" in
    start)
        $1
        ;;
    stop)
        $1
        ;;
    *)
        echo "Usage: $0 {start|stop}"
        exit 2
esac`

// serviceInit serviceInit
func serviceInit() {

	if 2 == len(os.Args) {

		serviceName := "monexecd"
		path := constant.LinuxStartShellFile + serviceName

		if "install" == os.Args[1] {

			_, err := os.Stat(path)
			if err != nil && os.IsExist(err) {

				log.Printf("[service-%s-exits]", serviceName)

			} else {

				file, err := os.Create(path)
				if err != nil {
					log.Panicf("[CannotCreateFile] %s", err)
				}
				defer file.Close()

				if _, err := fmt.Fprintf(file, shell); err != nil {
					log.Panicf("[WriteFile-error] %s", err)
				}
				exec.Command("/bin/bash", "-c", "chmod", "+x", constant.LinuxBinPath)
				exec.Command("/bin/bash", "-c", "update-rc.d", serviceName, "defaults", "91")

				log.Printf("[serviceInstall-success] !!!!!!!!!!!")
			}
			time.Sleep(time.Duration(3 * time.Second))
			os.Exit(0)

		} else if "remove" == os.Args[1] {

			_, err := os.Stat(path)
			if err != nil && os.IsNotExist(err) {

				log.Printf("[service-%s-not-exits]", serviceName)

			} else {

				exec.Command("/bin/bash", "-c", "update-rc.d", "-f", serviceName, "remove")
				exec.Command("/bin/bash", "-c", "rm", "-rf", path)
			}

			log.Printf("[serviceRemove-success] !!!!!!!!!!!")
			time.Sleep(time.Duration(3 * time.Second))
			os.Exit(0)
		}
	}
}
