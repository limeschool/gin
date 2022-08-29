package gin

import (
	"time"
)

type systemConfig struct {
	ClientTimeout  clientTimeout   `json:"client_timeout" mapstructure:"client_timeout"`
	CpuThreshold   cpuThreshold    `json:"cpu_threshold" mapstructure:"cpu_threshold"`
	ClientLimit    clientLimit     `json:"client_limit" mapstructure:"client_limit"`
	SkipRequestLog map[string]bool `json:"skip_request_log" mapstructure:"skip_request_log"`
}

type clientTimeout struct {
	Enable  bool          `json:"enable" mapstructure:"enable"`
	Timeout time.Duration `json:"timeout" mapstructure:"timeout"`
}

type clientLimit struct {
	Enable    bool  `json:"enable" mapstructure:"enable"`
	Threshold int64 `json:"threshold"  mapstructure:"threshold"`
}

type cpuThreshold struct {
	Enable    bool  `json:"enable" mapstructure:"enable"`
	Threshold int64 `json:"threshold"  mapstructure:"threshold"`
}

func initSystemConfig() systemConfig {
	conf := systemConfig{
		ClientTimeout: clientTimeout{
			Enable:  true,
			Timeout: 10 * time.Second,
		},
		CpuThreshold: cpuThreshold{
			Enable:    true,
			Threshold: 900,
		},
		ClientLimit: clientLimit{
			Enable:    true,
			Threshold: 50,
		},
	}
	if err := globalConfig.UnmarshalKey("system", &conf); err != nil {
		panic(err)
	}
	return conf
}
