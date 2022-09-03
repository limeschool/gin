package config_drive

import (
	"bytes"
	"github.com/samuel/go-zookeeper/zk"
	"github.com/spf13/viper"
	"strings"
	"time"
)

type zookeeper struct {
	client *zk.Conn
	wait   *uint64
	path   string
	tp     string
}

func NewZK(conf *Config) (ConfigService, error) {
	client, _, err := zk.Connect(strings.Split(conf.Host, ","), 10*time.Second)
	if err != nil {
		return nil, err
	}
	if conf.Username != "" {
		if err = client.AddAuth("digest", []byte(conf.Username+":"+conf.Password)); err != nil {
			panic(err)
		}
	}
	return &zookeeper{client: client, path: conf.Path, tp: conf.Type}, nil
}

func (z *zookeeper) Init() *viper.Viper {
	v := viper.New()
	//v := viper.NewWithOptions(viper.)
	v.SetConfigType(z.tp)
	if err := z.GetViper(v); err != nil {
		panic("zk config get err:" + err.Error())
	}
	go z.Watch(v)
	return v
}

func (z *zookeeper) GetViper(v *viper.Viper) error {
	data, err := z.Get()
	if err != nil {
		return err
	}
	if err = v.ReadConfig(bytes.NewBuffer(data)); err != nil {
		return err
	}
	return nil
}

func (z *zookeeper) Get() ([]byte, error) {
	data, _, err := z.client.Get(z.path)
	return data, err
}

func (z *zookeeper) Set(value string) error {
	_, err := z.client.Set(z.path, []byte(value), 1.0)
	return err
}

func (z *zookeeper) Watch(v *viper.Viper) {
	for {
		data, _, event, err := z.client.GetW(z.path)
		if err != nil {
			break
		}
		evt := <-event
		if evt.Type == zk.EventNodeDataChanged {
			if err = v.ReadConfig(bytes.NewBuffer(data)); err != nil {
				continue
			}
			if CallBack != nil {
				CallBack(v)
			}
		}
	}
}

func (z *zookeeper) SetPath(key string) {
	z.path += key
}
