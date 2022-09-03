package config_drive

import (
	"bytes"
	"context"
	"errors"
	"github.com/spf13/viper"
	"go.etcd.io/etcd/api/v3/mvccpb"
	etcdApi "go.etcd.io/etcd/client/v3"
	"strings"
	"time"
)

type etcd struct {
	client *etcdApi.Client
	wait   *uint64
	path   string
	tp     string
}

func NewEtcd(conf *Config) (ConfigService, error) {
	client, err := etcdApi.New(etcdApi.Config{
		Endpoints:   strings.Split(conf.Host, ","),
		DialTimeout: 10 * time.Second,
		Username:    conf.Username,
		Password:    conf.Password,
	})
	if err != nil {
		return nil, err
	}
	ctx, _ := context.WithTimeout(context.TODO(), 10*time.Second)
	if _, err = client.Get(ctx, "test"); err != nil {
		return nil, errors.New("etcd connect fail")
	}
	return &etcd{client: client, path: conf.Path, tp: conf.Type}, nil
}

func (c *etcd) Init() *viper.Viper {
	v := viper.New()
	v.SetConfigType(c.tp)
	if err := c.GetViper(v); err != nil {
		panic("etcd config get err:" + err.Error())
	}
	go c.Watch(v)
	return v
}

func (c *etcd) Set(value string) error {
	ctx, _ := context.WithTimeout(context.TODO(), 10*time.Second)
	_, err := c.client.KV.Put(ctx, c.path, value)
	return err
}

// Get 从中间件中获取配置
func (c *etcd) GetViper(v *viper.Viper) error {
	data, err := c.Get()
	if err != nil {
		return err
	}
	if err = v.ReadConfig(bytes.NewBuffer(data)); err != nil {
		return err
	}
	return nil
}

// Get 从中间件中获取配置
func (c *etcd) Get() ([]byte, error) {
	data, err := c.client.KV.Get(context.TODO(), c.path)
	if err != nil {
		return nil, err
	}
	if len(data.Kvs) == 0 {
		return nil, errors.New("not exist configure")
	}
	return data.Kvs[0].Value, nil
}

// Watch 监听配置变更
func (c *etcd) Watch(v *viper.Viper) {
	for {
		ch := c.client.Watch(context.TODO(), c.path)
		for ws := range ch {
			for _, event := range ws.Events {
				switch event.Type {
				case mvccpb.PUT:
					if err := v.ReadConfig(bytes.NewBuffer(event.Kv.Value)); err != nil {
						continue
					}
					if CallBack != nil {
						CallBack(v)
					}
				}
			}
		}
	}
}

func (c *etcd) SetPath(key string) {
	c.path += key
}
