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

const pid = "PID=`ps -ef |grep \"${NAME_BIN}\" |grep -v \"grep\" |grep -v \"init.d\" |grep -v \"service\" |awk '{print $2}'`"
const shell = `#!/bin/bash
# chkconfig: 2345 90 10
# description: monexec is a daemon process

### BEGIN INIT INFO
# Provides:          monexec
# Required-Start:    $network $syslog
# Required-Stop:     $network
# Default-Start:     2 3 4 5
# Default-Stop:      0 1 6
# Short-Description: monexec
# Description:       start or stop the monexec
### END INIT INFO

ROOT_PATH="/home/i5/bin/"
NAME="monexec"
NAME_BIN="monexec"

check_running(){
    ` + pid + `    
	if [[ ! -z ${PID} ]]; then
		return 0
	else
		return 1
	fi
}

start() {
	check_running
	if [[ $? -eq 0 ]]; then
			echo -e "${NAME} (PID ${PID}) running..." && exit 0
        else
			start-stop-daemon --start -b --quiet --oknodo --pidfile /var/run/${NAME}.pid --exec ${ROOT_PATH}${NAME_BIN}
			sleep 1s
			check_running
			if [[ $? -eq 0 ]]; then
			echo -e "${NAME} (PID ${PID}) start success !"
			else
			echo -e "${NAME} start fail !"
			fi
	fi
	
}

stop() {
	check_running
	if [[ $? -eq 0 ]]; then
		ps aux|grep ${NAME_BIN}|grep -v grep|awk '{print $2}'|xargs kill -9
	fi
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
	restart)
		$1
		;;
    *)
        echo "Usage: $0 {start|stop|restart}"
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
				_, err = util.RunCmd("/bin/bash", "-c", "cd "+constant.LinuxStartShellFile+";update-rc.d "+serviceName+" defaults 91")
				if err != nil {
					log.Panicf("[update-rc.d-error] %s\n", err)
				}
				log.Printf("[serviceInstall-success] !!!!!!!!!!!")
			}
			time.Sleep(time.Duration(2 * time.Second))
			os.Exit(0)

		} else if "remove" == os.Args[1] {

			_, err := os.Stat(path)
			if err != nil && os.IsNotExist(err) {

				log.Printf("[service-%s-not-exits]", serviceName)

			} else {

				out, err := util.RunCmd("/bin/bash", "-c", "update-rc.d", "-f", serviceName, "remove")
				if err != nil {
					log.Panicf("[update-rc.d-remove-error] %s\n", err)
				}
				log.Println(out)
				if _, err = util.RunCmd("/bin/bash", "-c", "rm", "-rf", path); err != nil {
					log.Panicf("[deleteServiceFile-error] %s\n", err)
				}
			}

			log.Printf("[serviceRemove-success] !!!!!!!!!!!")
			time.Sleep(time.Duration(2 * time.Second))
			os.Exit(0)
		}
	}
}
