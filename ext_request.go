package gin

import (
	"encoding/json"
	"github.com/go-resty/resty/v2"
	"go.uber.org/zap"
	"time"
)

type requestConfig struct {
	EnableLog        bool          `json:"enable_log" mapstructure:"enable_log"`
	RetryCount       int           `json:"retry_count" mapstructure:"retry_count"`
	RetryWaitTime    time.Duration `json:"retry_wait_time" mapstructure:"retry_wait_time"`
	MaxRetryWaitTime time.Duration `json:"max_retry_wait_time" mapstructure:"max_retry_wait_time"`
	Timeout          time.Duration `json:"timeout" mapstructure:"timeout"`
	RequestMsg       string        `json:"request_msg" mapstructure:"request_msg"`
	ResponseMsg      string        `json:"response_msg" mapstructure:"response_msg"`
}

func initRequestConfig() requestConfig {
	conf := requestConfig{
		EnableLog:        true,
		Timeout:          10 * time.Second,
		RetryCount:       3,
		RetryWaitTime:    2 * time.Second,
		MaxRetryWaitTime: 4 * time.Second,
		RequestMsg:       "request",
		ResponseMsg:      "response",
	}
	if err := globalConfig.UnmarshalKey("request", &conf); err != nil {
		panic(err)
	}
	return conf
}

type request struct {
	ctx     *Context
	request *resty.Request
}

type RequestFunc func(*resty.Request) *resty.Request

func (h *request) Option(fn RequestFunc) *request {
	h.request = fn(h.request)
	return h
}

func (h *request) log() {
	if !globalRequestConfig.EnableLog {
		return
	}
	logs := []zap.Field{
		zap.Any("method", h.request.Method),
		zap.Any("url", h.request.URL),
		zap.Any("header", h.request.Header),
		zap.Any("body", h.request.Body),
	}
	if len(h.request.FormData) != 0 {
		logs = append(logs, zap.Any("form-data", h.request.FormData))
	}
	if len(h.request.QueryParam) != 0 {
		logs = append(logs, zap.Any("query-data", h.request.QueryParam))
	}
	h.ctx.Log.Info(globalRequestConfig.RequestMsg, logs...)
}

func (h *request) Get(url string) *response {
	defer h.log()
	res := &response{ctx: h.ctx}
	res.response, res.err = h.request.Get(url)
	return res
}

func (h *request) Post(url string, data interface{}) *response {
	defer h.log()
	res := &response{ctx: h.ctx}
	res.response, res.err = h.request.SetBody(data).Post(url)
	return res
}

func (h *request) PostJson(url string, data interface{}) *response {
	defer h.log()
	res := &response{ctx: h.ctx}
	res.response, res.err = h.request.ForceContentType("application/json").SetBody(data).Post(url)
	return res
}

func (h *request) Put(url string, data interface{}) *response {
	defer h.log()
	res := &response{ctx: h.ctx}
	res.response, res.err = h.request.SetBody(data).Put(url)
	return res
}

func (h *request) PutJson(url string, data interface{}) *response {
	defer h.log()
	res := &response{ctx: h.ctx}
	res.response, res.err = h.request.ForceContentType("application/json").SetBody(data).Put(url)
	return res
}

func (h *request) Delete(url string) *response {
	defer h.log()
	res := &response{ctx: h.ctx}
	res.response, res.err = h.request.Delete(url)
	return res
}

type response struct {
	ctx      *Context
	err      error
	response *resty.Response
}

func (r *response) log() {
	if !globalRequestConfig.EnableLog {
		return
	}

	logs := []zap.Field{
		zap.Any("status", r.response.Status()),
		zap.Any("time", r.response.Time()),
		zap.Any("body", string(r.response.Body())),
		zap.Any("error", r.err),
	}
	r.ctx.Log.Info(globalRequestConfig.ResponseMsg, logs...)
}

func (r *response) Body() ([]byte, error) {
	defer r.log()
	return r.response.Body(), r.err
}

func (r *response) Result(val interface{}) error {
	defer r.log()
	if r.err != nil {
		return r.err
	}
	return json.Unmarshal(r.response.Body(), val)
}
