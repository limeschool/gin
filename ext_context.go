package gin

import (
	"context"
	"github.com/casbin/casbin/v2"
	"github.com/go-redis/redis/v8"
	"github.com/go-resty/resty/v2"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
	"sync"
)

type Option struct {
	TraceID string
}

type option func(*Option)

func (u *Option) Option(opts ...option) {
	for _, opt := range opts {
		opt(u)
	}
}

func WithTraceID(id string) option {
	return func(o *Option) {
		o.TraceID = id
	}
}

func NewContext(opts ...option) *Context {
	opt := &Option{}
	opt.Option(opts...)

	if opt.TraceID == "" {
		opt.TraceID = uuid.New().String()
	}
	log := newLog(opt.TraceID)
	return &Context{
		mu:          sync.RWMutex{},
		Keys:        map[string]any{},
		ServiceName: globalServiceName,
		TraceID:     opt.TraceID,
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

func (c *Context) RespJson(data interface{}) {
	ResponseJson(c, data)
}

func (c *Context) RespXml(data interface{}) {
	ResponseXml(c, data)
}

func (c *Context) RespList(page, count, total int, data interface{}) {
	ResponseList(c, page, count, total, data)
}

func (c *Context) Mysql(name string) *gorm.DB {
	db, ok := globalMysql[name]
	if !ok {
		return nil
	}

	return db.WithContext(c.Context())
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

func (c *Context) Http() *request {
	client := resty.New()
	client.SetRetryWaitTime(globalRequestConfig.RetryWaitTime)
	client.SetRetryMaxWaitTime(globalRequestConfig.MaxRetryWaitTime)
	client.SetRetryCount(globalRequestConfig.RetryCount)
	client.SetTimeout(globalRequestConfig.Timeout)
	return &request{ctx: c, request: client.R()}
}
