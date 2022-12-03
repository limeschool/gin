package config_drive

import (
	"bytes"
	"errors"
	consulApi "github.com/hashicorp/consul/api"
	"github.com/spf13/viper"
	"log"
	"strings"
)

type consul struct {
	client *consulApi.Client
	wait   *uint64
	path   string
	tp     string
}

// NewConsul 创建consul对象，用于获取和监听日志
func NewConsul(conf *Config) (ConfigService, error) {
	client, err := consulApi.NewClient(&consulApi.Config{
		Address: conf.Host,
		Token:   conf.Token,
	})
	if err != nil {
		return nil, err
	}

	if _, err = client.Status().Leader(); err != nil {
		return nil, errors.New("consul connect fail")
	}

	return &consul{client: client, path: conf.Path, tp: conf.Type}, nil
}

// Init 初始化consul配置信息
func (c *consul) Init() *viper.Viper {
	v := viper.New()
	v.SetConfigType(c.tp)
	if err := c.GetViper(v); err != nil {
		panic("consul config get err:" + err.Error())
	}
	go c.Watch(v)
	return v
}

// GetViper 从中间件中获取配置
func (c *consul) GetViper(v *viper.Viper) error {
	data, err := c.Get()
	if err != nil {
		return err
	}
	if err = v.ReadConfig(bytes.NewBuffer(data)); err != nil {
		return err
	}
	return nil
}

func (c *consul) Get() ([]byte, error) {
	var q *consulApi.QueryOptions
	if c.wait != nil {
		q = &consulApi.QueryOptions{
			WaitIndex: *c.wait,
		}
	}

	data, meta, err := c.client.KV().Get(c.path, q)
	if err != nil {
		return nil, err
	}

	c.wait = &meta.LastIndex
	if data == nil {
		return nil, nil
	}
	return data.Value, nil
}

func (c *consul) Set(value string) error {
	path := strings.TrimPrefix(c.path, "/")
	kv := &consulApi.KVPair{Key: path, Value: []byte(value)}
	_, err := c.client.KV().Put(kv, nil)
	return err
}

// Watch 监听配置变更
func (c *consul) Watch(v *viper.Viper) {
	for {
		if err := c.GetViper(v); err != nil {
			log.Println("consul监听变更信息获取失败")
			continue
		}
		if CallBack != nil {
			CallBack(v)
		}
	}
}

func (c *consul) SetPath(key string) {
	c.path += key
}
