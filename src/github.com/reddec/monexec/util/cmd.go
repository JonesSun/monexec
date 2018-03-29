package util

import (
	"os/exec"
	"io/ioutil"
	"log"
	"unsafe"
)

//RunCmd runCmd
func RunCmd(name string, arg ... string) (string, error) {

	cmd := exec.Command(name, arg...)
	stdout, err := cmd.StdoutPipe()
	if err := cmd.Start(); err != nil {
		log.Printf("[runCmd-error] %s\n", err)
		return "", err
	}
	content, err := ioutil.ReadAll(stdout)
	if err != nil {
		log.Printf("[runCmd-error] %s\n", err)
		return "", err
	}

	return B2S(content), nil
}

// B2S byte to string
func B2S(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}
