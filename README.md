# Gin Web Framework

### 基于配置编程
- config/dev.json
```
{
  "service": "test",
  "log": {
    "level": 0
  },
  "ip_limit": {
    "max": 10
  },
  "http_tool":{
    "retry_count": 3,
    "retry_wait_time":"100ms",
    "max_retry_wait_time":"2s",
    "timeout": "10s",
    "enable_log": true,
    "request_msg": "http request info",
    "response_msg": "http response res"
  },
  "system": {
    "timeout": "5s"
  },
  "config": {
    "test": "111"
  }
}

```
- main.go
```
package main

import "github.com/limeschool/gin"

func main() {
	e := gin.Default()
	e.GET("/test", func(ctx *gin.Context) {
		ctx.Mysql("devops")
		ctx.Mongo("devops")
		ctx.Redis("devops")
		ctx.Config.Get("config.test")
	})
	e.Run(":8000")
}
```
