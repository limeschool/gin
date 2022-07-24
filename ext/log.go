package ext

import (
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

type LogConfig struct {
	Level      int8   `json:"level" mapstructure:"level"`
	TraceKey   string `json:"trace_key" mapstructure:"trace_key"`
	ServiceKey string `json:"service_key" mapstructure:"service_key"`
}

func parseLog(v *viper.Viper) LogConfig {
	conf := LogConfig{}
	if err := v.UnmarshalKey("log", &conf); err != nil {
		panic(err)
	}
	return conf
}

func initLogKeys(v *viper.Viper) {
	if v.GetString("trace_key") != "" {
		TraceID = v.GetString("trace_key")
	}
	if v.GetString("service_key") != "" {
		ServiceID = v.GetString("service_key")
	}
}

// newLog 链路日志
func newLog(id string) *zap.Logger {
	return globalLog.With(zap.Any(TraceID, id))
}

func initLog() *zap.Logger {
	conf := parseLog(globalConfig)
	initLogKeys(globalConfig)
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,                          // 小写编码器
		EncodeTime:     zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05.000"), // ISO8601 UTC 时间格式
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.FullCallerEncoder, // 全路径编码器
	}

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),  // 编码器配置
		zapcore.NewMultiWriteSyncer(os.Stdout), // 输出方式
		zapcore.Level(conf.Level),              // 设置日志级别
	)

	return zap.New(core,
		zap.AddCaller(),   // 开启文件及行号,设置初始化字段
		zap.Development(), // 开启开发模式，堆栈跟踪
		zap.Fields(zap.String(ServiceID, globalServiceName)),
	)
}
