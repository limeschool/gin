package gin

import (
	"bytes"
	"encoding/json"
)

type List struct {
	Code int      `json:"code"`
	Msg  string   `json:"msg"`
	Data ListData `json:"data"`
}

type ListData struct {
	Page  int         `json:"page"`
	Count int         `json:"count"`
	Total int         `json:"total"`
	List  interface{} `json:"list"`
}

type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

type CustomError struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

func (c *CustomError) Error() string {
	return c.Msg
}

func ParseResponse(data []byte) *Response {
	var resp = new(Response)
	_ = json.Unmarshal(data, resp)
	return resp
}

func ResponseData(ctx *Context, data interface{}) {
	ctx.mu.Lock()
	defer ctx.mu.Unlock()

	if ctx.Writer.Size() != -1 {
		return
	}
	ctx.JSON(200, &Response{
		Code: 200,
		Msg:  "success",
		Data: data,
	})
}

func ResponseList(ctx *Context, page, count, total int, data interface{}) {
	ctx.mu.Lock()
	defer ctx.mu.Unlock()

	if ctx.Writer.Size() != -1 {
		return
	}
	ctx.JSON(200, &List{
		Code: 200,
		Msg:  "success",
		Data: ListData{
			Page:  page,
			Count: count,
			Total: total,
			List:  data,
		},
	})
}

func ResponseError(ctx *Context, err error) {
	ctx.mu.Lock()
	defer ctx.mu.Unlock()

	if ctx.Writer.Size() != -1 {
		return
	}
	if customErr, is := err.(*CustomError); is {
		ctx.JSON(200, &Response{
			Code: customErr.Code,
			Msg:  customErr.Msg,
		})
	} else {
		ctx.JSON(200, &Response{
			Code: 400,
			Msg:  err.Error(),
		})
	}
}

func CustomResponseError(ctx *Context, code int, msg string) {
	ctx.mu.Lock()
	defer ctx.mu.Unlock()

	if ctx.Writer.Size() != -1 {
		return
	}
	ctx.JSON(code, &Response{
		Code: code,
		Msg:  msg,
	})
}

type ResponseWriterWrapper struct {
	ResponseWriter
	Body *bytes.Buffer // 缓存
}

func (w ResponseWriterWrapper) Write(b []byte) (int, error) {
	w.Body.Write(b)
	return w.ResponseWriter.Write(b)
}

func (w ResponseWriterWrapper) WriteString(s string) (int, error) {
	w.Body.WriteString(s)
	return w.ResponseWriter.WriteString(s)
}
