package gin

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
)

type redisConfig struct {
	Enable   bool   `json:"enable" mapstructure:"enable"` //是否启用redis
	Name     string `json:"name"  mapstructure:"name"`
	Network  string `json:"network" mapstructure:"network"`
	Host     string `json:"host" mapstructure:"host"` //redis的连接地址
	Username string `json:"username" mapstructure:"username"`
	Password string `json:"password" mapstructure:"password"` //redis的密码
	DB       int    `json:"db" mapstructure:"db"`
	PoolSize int    `json:"pool_size" mapstructure:"pool_size"`
}

func parseRedisConfig(v *viper.Viper) (conf []redisConfig) {
	if v == nil {
		return
	}
	if err := v.UnmarshalKey("redis", &conf); err != nil {
		panic("redis 配置解析错误" + err.Error())
	}
	return
}

func initRedis() {
	confList := parseRedisConfig(globalConfig)
	clients := make(map[string]*redis.Client)
	for _, conf := range confList {
		if !conf.Enable {
			return
		}
		client := redis.NewClient(&redis.Options{
			Network:  conf.Network,
			Addr:     conf.Host,
			Username: conf.Username,
			Password: conf.Password,
			PoolSize: conf.PoolSize,
			DB:       conf.DB,
		})
		if err := client.Ping(context.TODO()).Err(); err != nil {
			panic(fmt.Sprintf("redis %v 连接失败%v", conf.Name, err.Error()))
		}
		clients[conf.Name] = client
	}

	globalRedis = clients
}
