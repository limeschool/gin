package gin

import (
	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	adapter "github.com/casbin/gorm-adapter/v3"
)

type ExtRbacConfig struct {
	Enable bool   `json:"enable" mapstructure:"enable"`
	DB     string `json:"db" mapstructure:"db"`
}

func initRbac() {
	conf := ExtRbacConfig{}
	_ = globalConfig.UnmarshalKey("rbac", &conf)
	if !conf.Enable {
		return
	}

	m := model.NewModel()
	m.AddDef("r", "r", "sub, obj, act")
	m.AddDef("p", "p", "sub, obj, act")
	m.AddDef("g", "g", "_, _")
	m.AddDef("e", "e", "some(where (p.eft == allow))")
	m.AddDef("m", "m", "g(r.sub, p.sub) && r.obj == p.obj && r.act == p.act")

	a, err := adapter.NewAdapterByDB(globalMysql[conf.DB])
	if err != nil {
		panic(err)
	}
	e, _ := casbin.NewEnforcer(m, a)
	if err = e.LoadPolicy(); err != nil {
		panic(err)
	}
	globalRbac = e
}
