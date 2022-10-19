package gin

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/didip/tollbooth"
	"github.com/google/uuid"
	"github.com/zeromicro/go-zero/core/load"
	"go.uber.org/zap"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

func ExtInit() HandlerFunc {
	return func(ctx *Context) {
		ctx.ServiceName = globalServiceName
	}
}

func ExtLogger() HandlerFunc {
	return func(ctx *Context) {
		trace := ctx.GetHeader(TraceID)
		if trace == "" {
			trace = uuid.New().String()
		}
		ctx.TraceID = trace
		ctx.Set(TraceID, trace)
		ctx.Log = newLog(trace)
		ctx.Config = newConfig(ctx.Log)
	}
}

func ExtRecovery(ctx *Context) {
	defer func() {
		if err := recover(); err != nil {
			message := fmt.Sprintf("%s", err)
			ctx.Log.WithOptions(zap.AddCallerSkip(1)).Error(message,
				zap.Any("panic", panicErr()),
				zap.Any("params", requestParams(ctx)),
			)
			CustomResponseError(ctx, http.StatusInternalServerError, "Internal Server Fail")
			ctx.Abort()
		}
	}()
	ctx.Next()
}

// ExtTimeout 客户端请求超时
func ExtTimeout() HandlerFunc {
	return func(ctx *Context) {
		t := globalSystemConfig.ClientTimeout
		if !t.Enable || t.Timeout == 0 {
			return
		}

		c, cancel := context.WithTimeout(ctx, t.Timeout)
		defer cancel()
		ctx.Request.WithContext(c)

		done := make(chan struct{}, 1)
		go func() {
			ExtRecovery(ctx)
			close(done)
		}()

		select {
		case <-c.Done():
			CustomResponseError(ctx, http.StatusInternalServerError, "Internal Server Timeout")
			ctx.Abort()
		case <-done:
		}
	}
}

// ExtCpuLoad Cpu 自适应降载
func ExtCpuLoad() HandlerFunc {
	sd := load.NewAdaptiveShedder(load.WithCpuThreshold(globalSystemConfig.CpuThreshold.Threshold))
	return func(ctx *Context) {
		if !globalSystemConfig.CpuThreshold.Enable {
			return
		}
		promise, err := sd.Allow()
		if err != nil {
			CustomResponseError(ctx, http.StatusInternalServerError, "Internal Server Over Max Request")
			ctx.Abort()
		}
		promise.Pass()
	}
}

func ExtLimit() HandlerFunc {
	max := globalSystemConfig.ClientLimit.Threshold
	limit := tollbooth.NewLimiter(float64(max), nil)
	return func(ctx *Context) {
		if !globalSystemConfig.ClientLimit.Enable {
			return
		}
		if httpError := tollbooth.LimitByRequest(limit, ctx.Writer, ctx.Request); httpError != nil {
			CustomResponseError(ctx, http.StatusInternalServerError, "Internal Server Limit")
			ctx.Abort()
		}
	}
}

// Success 健康检查
func Success() HandlerFunc {
	return func(ctx *Context) {
		ctx.RespSuccess()
	}
}

func Resp404() HandlerFunc {
	return func(ctx *Context) {
		ctx.RespError(errors.New("此接口不存在"))
	}
}

// ExtCors 允许跨域
func ExtCors() HandlerFunc {
	return func(c *Context) {
		// 允许 Origin 字段中的域发送请求
		c.Writer.Header().Add("Access-Control-Allow-Origin", "*")
		// 设置预验请求有效期为 86400 秒
		c.Writer.Header().Set("Access-Control-Max-Age", "86400")
		// 设置允许请求的方法
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE, PATCH")
		// 设置允许请求的 Header
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Token, X-Token, Authorization")
		// 设置拿到除基本字段外的其他字段，如上面的Apitoken, 这里通过引用Access-Control-Expose-Headers，进行配置，效果是一样的。
		c.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Headers")
		// 配置是否可以带认证信息
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		// OPTIONS请求返回200
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(200)
		} else {
			c.Next()
		}
	}
}

func ExtRequestInfo() HandlerFunc {
	return func(ctx *Context) {
		method := strings.ToLower(ctx.Request.Method)
		path := ctx.Request.URL.String()
		if strings.Contains(path, "?") {
			path = strings.Split(path, "?")[0]
		}
		if globalSystemConfig.SkipRequestLog[method+":"+path] {
			return
		}
		now := time.Now()
		blw := &ResponseWriterWrapper{Body: bytes.NewBufferString(""), ResponseWriter: ctx.Writer}
		ctx.Writer = blw
		req := requestParams(ctx)
		ctx.Request.Body = ioutil.NopCloser(bytes.NewBuffer([]byte(req["body"])))
		ctx.Next()
		//请求结束
		end := time.Now()
		ctx.Log.Info("request",
			zap.Any("path", ctx.Request.URL.Path),
			zap.Any("method", ctx.Request.Method),
			zap.Any("timestamp", end.Sub(now)),
			zap.Any("req", req),
			zap.Any("res", blw.Body.String()),
			zap.Any("status", ctx.Writer.Status()),
		)
	}
}

func ExtRequestTokenAuth() HandlerFunc {
	return func(ctx *Context) {
		token := ctx.Config.Get("request-token")
		if ctx.Request.Header.Get("token") == token {
			ctx.Next()
		} else {
			ctx.RespError(errors.New("token 验证失败"))
			ctx.Abort()
		}
	}
}
