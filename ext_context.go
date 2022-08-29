package gin

import (
	"context"
	"github.com/casbin/casbin/v2"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
	"sync"
)

func NewContext() *Context {
	traceId := uuid.New().String()
	log := newLog(traceId)
	return &Context{
		mu:          sync.RWMutex{},
		Keys:        map[string]any{},
		ServiceName: globalServiceName,
		TraceID:     traceId,
		Log:         log,
		Config:      newConfig(log),
	}
}

func (c *Context) RespSuccess() {
	ResponseData(c, "")
}

func (c *Context) RespError(err error) {
	ResponseError(c, err)
}

func (c *Context) RespData(data interface{}) {
	ResponseData(c, data)
}

func (c *Context) RespList(page, count, total int, data interface{}) {
	ResponseList(c, page, count, total, data)
}

func (c *Context) Mysql(name string) *gorm.DB {
	return globalMysql[name].WithContext(c.Context())
}

func (c *Context) Mongo(name string) *mongo.Client {
	return globalMongo[name]
}

func (c *Context) Rbac() *casbin.Enforcer {
	return globalRbac
}

func (c *Context) Rsa(name string) *ExtRsa {
	return globalRsa[name]
}

func (c *Context) Redis(name string) *redis.Client {
	return globalRedis[name]
}

func (c *Context) Context() context.Context {
	ctx := context.TODO()
	for key, val := range c.Keys {
		ctx = context.WithValue(ctx, key, val)
	}
	return ctx
}
