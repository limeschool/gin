package gin

import (
	"github.com/limeschool/gin/config_drive"
	"go.uber.org/zap"
	"sync"
	"time"
)

type IConfig interface {
	Get(key string) interface{}
	GetString(key string) string
	GetDefaultString(key string, def string) string
	GetBool(key string) bool
	GetInt(key string) int
	GetDefaultInt(key string, def int) int
	GetInt32(key string) int32
	GetInt64(key string) int64
	GetUint(key string) uint
	GetUint32(key string) uint32
	GetUint64(key string) uint64
	GetFloat64(key string) float64
	GetTime(key string) time.Time
	GetDuration(key string) time.Duration
	GetDefaultDuration(key string, def time.Duration) time.Duration
	GetIntSlice(key string) []int
	GetStringSlice(key string) []string
	GetStringMap(key string) map[string]interface{}
	GetStringMapString(key string) map[string]string
	GetStringMapStringSlice(key string) map[string][]string
	UnmarshalKey(key string, val interface{}) error
	Unmarshal(val interface{}) error
}

type Config struct {
	logger *zap.Logger
}

var configPool = sync.Pool{New: func() any {
	return &Config{
		logger: nil,
	}
}}

func newConfig(log *zap.Logger) IConfig {
	conf := configPool.Get().(*Config)
	conf.logger = log.WithOptions(zap.AddCaller(), zap.AddCallerSkip(1))
	configPool.Put(conf)
	return conf
}

func WatchConfig(f config_drive.CallFunc) {
	config_drive.CallBack = f
	// 这只之后进行初始化执行
	f(globalConfig)
}

func (c *Config) Set(key string, value interface{}) {
	globalConfig.Set(key, value)
	c.logger.Info(SetConfigTip, zap.Any(key, value))
}

func (c *Config) Get(key string) interface{} {
	res := globalConfig.Get(key)
	c.logger.Info(GetConfigTip, zap.Any(key, res))
	return res
}

func (c *Config) GetString(key string) string {
	res := globalConfig.GetString(key)
	c.logger.Info(GetConfigTip, zap.Any(key, res))
	return res
}

func (c *Config) GetDefaultString(key string, def string) string {
	res := globalConfig.GetString(key)
	if res == "" {
		res = def
	}
	c.logger.Info(GetConfigTip, zap.Any(key, res))
	return res
}

func (c *Config) GetBool(key string) bool {
	res := globalConfig.GetBool(key)
	c.logger.Info(GetConfigTip, zap.Any(key, res))
	return res
}

func (c *Config) GetInt(key string) int {
	res := globalConfig.GetInt(key)
	c.logger.Info(GetConfigTip, zap.Any(key, res))
	return res
}

func (c *Config) GetDefaultInt(key string, def int) int {
	res := globalConfig.GetInt(key)
	if res == 0 {
		res = def
	}
	c.logger.Info(GetConfigTip, zap.Any(key, res))
	return res
}

func (c *Config) GetInt32(key string) int32 {
	res := globalConfig.GetInt32(key)
	c.logger.Info(GetConfigTip, zap.Any(key, res))
	return res
}

func (c *Config) GetInt64(key string) int64 {
	res := globalConfig.GetInt64(key)
	c.logger.Info(GetConfigTip, zap.Any(key, res))
	return res
}

func (c *Config) GetUint(key string) uint {
	res := globalConfig.GetUint(key)
	c.logger.Info(GetConfigTip, zap.Any(key, res))
	return res
}

func (c *Config) GetUint32(key string) uint32 {
	res := globalConfig.GetUint32(key)
	c.logger.Info(GetConfigTip, zap.Any(key, res))
	return res
}

func (c *Config) GetUint64(key string) uint64 {
	res := globalConfig.GetUint64(key)
	c.logger.Info(GetConfigTip, zap.Any(key, res))
	return res
}

func (c *Config) GetFloat64(key string) float64 {
	res := globalConfig.GetFloat64(key)
	c.logger.Info(GetConfigTip, zap.Any(key, res))
	return res
}

func (c *Config) GetDefaultFloat64(key string, def float64) float64 {
	res := globalConfig.GetFloat64(key)
	if res == 0 {
		res = def
	}
	c.logger.Info(GetConfigTip, zap.Any(key, res))
	return res
}

func (c *Config) GetTime(key string) time.Time {
	res := globalConfig.GetTime(key)
	c.logger.Info(GetConfigTip, zap.Any(key, res))
	return res
}

func (c *Config) GetDuration(key string) time.Duration {
	res := globalConfig.GetDuration(key)
	c.logger.Info(GetConfigTip, zap.Any(key, res))
	return res
}

func (c *Config) GetDefaultDuration(key string, duration time.Duration) time.Duration {
	res := globalConfig.GetDuration(key)
	if res.String() == "" {
		res = duration
	}
	c.logger.Info(GetConfigTip, zap.Any(key, res))
	return res
}

func (c *Config) GetIntSlice(key string) []int {
	res := globalConfig.GetIntSlice(key)
	c.logger.Info(GetConfigTip, zap.Any(key, res))
	return res
}

func (c *Config) GetStringSlice(key string) []string {
	res := globalConfig.GetStringSlice(key)
	c.logger.Info(GetConfigTip, zap.Any(key, res))
	return res
}

func (c *Config) GetStringMap(key string) map[string]interface{} {
	res := globalConfig.GetStringMap(key)
	c.logger.Info(GetConfigTip, zap.Any(key, res))
	return res
}

func (c *Config) GetStringMapString(key string) map[string]string {
	res := globalConfig.GetStringMapString(key)
	c.logger.Info(GetConfigTip, zap.Any(key, res))
	return res
}

func (c *Config) GetStringMapStringSlice(key string) map[string][]string {
	res := globalConfig.GetStringMapStringSlice(key)
	c.logger.Info(GetConfigTip, zap.Any(key, res))
	return res
}

func (c *Config) UnmarshalKey(key string, val interface{}) error {
	defer c.logger.Info(GetConfigTip, zap.Any(key, val))
	return globalConfig.UnmarshalKey(key, val)
}

func (c *Config) Unmarshal(val interface{}) error {
	defer c.logger.Info(GetConfigTip, zap.Any("res", val))
	return globalConfig.Unmarshal(&val)
}
