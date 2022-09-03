package config_drive

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/spf13/viper"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
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

var configFile = flag.String("c", "config/dev.json", "the Config file path")

func GetConfig(srv string) *viper.Viper {
	flag.Parse()
	conf := &Config{}
	if configFile == nil {
		addr := os.Getenv("CONFIG_ADDR")
		token := os.Getenv("CONFIG_TOKEN")
		if addr == "" {
			panic("环境变量CONFIG_ADDR未配置")
		}
		if token == "" {
			panic("环境变量CONFIG_TOKEN未配置")
		}
		url := fmt.Sprintf("%v/configure/config?service=%v&token=%v", addr, srv, token)
		client := http.Client{Timeout: 10 * time.Second}
		response, err := client.Get(url)
		if err != nil {
			panic("请求配置中心信息异常" + err.Error())
		}
		defer response.Body.Close()
		respData := struct {
			Code int64   `json:"code"`
			Msg  string  `json:"msg"`
			Data *Config `json:"data"`
		}{}
		b, _ := io.ReadAll(response.Body)
		if json.Unmarshal(b, &respData) != nil {
			panic("解析配置中心失败")
		}
		if respData.Code != 200 || respData.Data == nil {
			panic("获取配置连接信息失败:" + respData.Msg)
		}
		conf = respData.Data
	} else {
		temp := strings.Split(*configFile, ".")
		conf = &Config{
			Drive: "local",
			Type:  temp[len(temp)-1],
			Path:  *configFile,
		}
	}
	return Init(conf)
}
