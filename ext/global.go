package ext

import (
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var (
	globalLog         *zap.Logger
	globalConfig      *viper.Viper
	globalServiceName string //服务名
)

func Init() {
	initConfig()
	initLog()
}
