package config_drive

import (
	"github.com/spf13/viper"
)

type Config struct {
	Drive    string //中间件 etcd/consul/zk
	Host     string //连接地址
	Type     string //配置数据格式 json\yaml...
	Username string //连接用户名
	Password string //连接密码
	Path     string //配置存储目录
	Token    string //连接的token
}

type ConfigService interface {
	Init() *viper.Viper
	GetViper(v *viper.Viper) error
	Get() ([]byte, error)
	Watch(v *viper.Viper)
	Set(value string) error
	SetPath(key string)
}

type CallFunc func(v *viper.Viper)

var CallBack func(v *viper.Viper)

func Init(conf *Config) *viper.Viper {
	var cs ConfigService
	var err error
	switch conf.Drive {
	case "etcd":
		cs, err = NewEtcd(conf)
	case "zk":
		cs, err = NewZK(conf)
	case "consul":
		cs, err = NewConsul(conf)
	case "local":
		cs, err = NewLocal(conf)
	default:
		panic("config type is fail")
	}
	if err != nil {
		panic("config drive fail:" + err.Error())
	}
	return cs.Init()
}
