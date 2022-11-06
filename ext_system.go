package gin

import (
	"time"
)

type systemConfig struct {
	Timeout        time.Duration   `json:"timeout" mapstructure:"timeout"`
	SkipRequestLog map[string]bool `json:"skip_request_log"  mapstructure:"skip_request_log"`
}

func initSystemConfig() systemConfig {
	conf := systemConfig{
		Timeout: 10 * time.Second,
	}
	if err := globalConfig.UnmarshalKey("system", &conf); err != nil {
		panic(err)
	}
	return conf
}
