package gin

import (
	"github.com/casbin/casbin/v2"
	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var (
	globalLog           *zap.Logger
	globalConfig        *viper.Viper
	globalServiceName   string //服务名
	globalRequestConfig requestConfig
	globalSystemConfig  systemConfig
	globalMysql         map[string]*gorm.DB
	globalMongo         map[string]*mongo.Client
	globalRedis         map[string]*redis.Client
	globalRbac          *casbin.Enforcer
	globalRsa           map[string]*ExtRsa
)

func initGlobal() {
	// 初始化全局配置
	globalConfig = initConfig()
	// 初始化全局服务名
	globalServiceName = globalConfig.GetString("service")
	if globalServiceName == "" {
		panic("Config service field not found")
	}
	// 初始化全局日志
	globalLog = initLog()
	// 初始化请求工具
	globalRequestConfig = initRequestConfig()
	// 初始化系统配置
	globalSystemConfig = initSystemConfig()
	// 初始化数据库
	initMysql()
	// 初始化mongo
	initMongo()
	// 初始化redis
	initRedis()
	// 初始化rbac
	initRbac()
	// 初始化rsa
	initRsa()
}
