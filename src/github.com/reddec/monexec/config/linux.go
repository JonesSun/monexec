// +build linux

package config

import "os"

// serviceInit serviceInit
func serviceInit() {

	if 2 == len(os.Args) {

		//serviceName := "monexecd"

		if "install" == os.Args[1] {

			//log.Printf("[serviceInstall-success] !!!!!!!!!!!")
			//time.Sleep(time.Duration(3 * time.Second))
			//os.Exit(0)

		} else if "remove" == os.Args[1] {

			//log.Printf("[serviceRemove-success] please system will reboot now!!!!!!!!!!!")
			//time.Sleep(time.Duration(3 * time.Second))
			//exec.Command("cmd", "/C", "shutdown -t 0 -r -f")
			//os.Exit(0)
		}
	}
}
