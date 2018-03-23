package util

import (
	"github.com/magiconair/properties"
	"os"
	"bufio"
	"github.com/reddec/monexec/constant"
	"errors"
	"runtime"
	"sync"
)

// Properties
type config struct {
	*properties.Properties
	lock *sync.RWMutex
}

var Config = config{
	lock: new(sync.RWMutex)}

var writeBoxConfLock = new(sync.Mutex)

// GetBoxConf
func (conf config) GetBoxConf(key string) (string, bool) {

	conf.lock.RLock()
	defer conf.lock.RUnlock()
	return conf.Properties.Get(key)
}

// SetBoxConf
func (conf config) SetBoxConf(key, val string) (error) {

	conf.lock.Lock()
	defer conf.lock.Unlock()

	if _, _, err := conf.Properties.Set(key, val); err != nil {
		return errors.New("[setBoxConfig-error] " + err.Error())
	}
	if err := conf.WriteBoxConf(); err != nil {
		return err
	}
	return nil
}

// SetBoxConfigs
func (conf config) SetBoxConfigs(req map[string]string) (error) {

	conf.lock.Lock()
	defer conf.lock.Unlock()

	for k, v := range req {
		if _, _, err := conf.Properties.Set(k, v); err != nil {
			return errors.New("[setBoxConfigs-error] " + err.Error())
		}
	}
	if err := conf.WriteBoxConf(); err != nil {
		return err
	}
	return nil
}

// WriteFile
func (conf config) WriteBoxConf() error {

	writeBoxConfLock.Lock()
	defer writeBoxConfLock.Unlock()

	var path string

	if "windows" == runtime.GOOS {
		path = constant.WinConfPath
	} else {
		path = constant.LinuxConfPath
	}

	outputFile, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {

		return err
	}
	defer outputFile.Close()

	outputWriter := bufio.NewWriter(outputFile)
	defer outputWriter.Flush()

	if _, err := conf.Properties.Write(outputWriter, properties.UTF8); err != nil {
		return err
	}
	return nil
}
