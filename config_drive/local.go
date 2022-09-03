package config_drive

import (
	"bytes"
	"errors"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"io/ioutil"
	"os"
)

type local struct {
	client *os.File
	wait   *uint64
	path   string
	tp     string
}

func NewLocal(conf *Config) (ConfigService, error) {
	client, err := os.Open(conf.Path)
	if err != nil {
		return nil, err
	}
	return &local{client: client, path: conf.Path, tp: conf.Type}, nil
}

func (c *local) Init() *viper.Viper {
	v := viper.New()
	v.SetConfigType(c.tp)
	if err := c.GetViper(v); err != nil {
		panic("local config get err:" + err.Error())
	}
	go c.Watch(v)
	return v
}

// GetViper 从中间件中获取配置
func (c *local) GetViper(v *viper.Viper) error {
	data, err := c.Get()
	if err != nil {
		return errors.New("get local config fail:" + err.Error())
	}
	if err = v.ReadConfig(bytes.NewBuffer(data)); err != nil {
		return err
	}
	return nil
}

func (c *local) Get() ([]byte, error) {
	c.client.Seek(0, 0) //重制到最开始
	return ioutil.ReadAll(c.client)
}

func (c *local) Set(value string) error {
	file, err := os.OpenFile(c.path, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return err
	}
	_, err = file.Write([]byte(value))
	return err
}

func (c *local) Watch(v *viper.Viper) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		panic(err)
	}

	if err = watcher.Add(c.path); err != nil {
		panic(err)
	}

	defer watcher.Close()

	for {
		select {
		case event, ok := <-watcher.Events: // 正常的事件的处理逻辑
			if !ok {
				continue
			}
			if event.Op&fsnotify.Write == fsnotify.Write {
				if err = c.GetViper(v); err != nil && CallBack != nil {
					CallBack(v)
				}
			}
		case _, ok := <-watcher.Errors:
			if !ok {
				break
			}
		}
	}
}

func (c *local) SetPath(key string) {
	c.path += key
}
