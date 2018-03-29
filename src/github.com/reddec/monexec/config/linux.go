// +build linux

package config

import (
	"os"
	"log"
	"time"
	"github.com/reddec/monexec/constant"
	"fmt"
	"github.com/reddec/monexec/util"
)

const shell = `#!/bin/sh
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
func ServiceInit() {

	if 2 == len(os.Args) {

		log.SetFlags(log.LstdFlags)
		log.SetOutput(os.Stdout)

		serviceName := "monexecd"
		path := constant.LinuxStartShellFile + serviceName
		log.SetFlags(log.LstdFlags)

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
				_, err = util.RunCmd("/bin/bash", "-c", "chmod +x "+constant.LinuxStartShellFile+serviceName)
				if err != nil {
					log.Panicf("[chmodx-error] %s\n", err)
				}
				out, err := util.RunCmd("/bin/bash", "-c", "cd "+constant.LinuxStartShellFile+";update-rc.d "+serviceName+" defaults 91")
				if err != nil {
					log.Panicf("[update-rc.d-error] %s\n", err)
				}
				log.Println(out)
				log.Printf("[serviceInstall-success] !!!!!!!!!!!")
			}
			time.Sleep(time.Duration(3 * time.Second))
			os.Exit(0)

		} else if "remove" == os.Args[1] {

			_, err := os.Stat(path)
			if err != nil && os.IsNotExist(err) {

				log.Printf("[service-%s-not-exits]", serviceName)

			} else {

				out, err := util.RunCmd("/bin/bash", "-c", "update-rc.d -f "+serviceName+" remove")
				if err != nil {
					log.Panicf("[update-rc.d-remove-error] %s\n", err)
				}
				log.Println(out)
				if _, err = util.RunCmd("/bin/bash", "-c", "rm -rf "+path); err != nil {
					log.Panicf("[deleteServiceFile-error] %s\n", err)
				}
			}

			log.Printf("[serviceRemove-success] !!!!!!!!!!!")
			time.Sleep(time.Duration(3 * time.Second))
			os.Exit(0)
		}
	}
}
