# Gin Web Framework
在学习此框架之前，首先你应该会使用gin框架，如果你还没有了解到gin框架，那你应该提前学习一下[传送们](https://github.com/gin-gonic/gin)

### 扩展gin框架的好处
在常规的开发过程中，会引入很多包，比如 redis、mysql、mongo、kafka 等，单独开发某一个应用进行使用的时候，其实感觉不出痛点，在企业或者一些大型的项目中，都会将某个系统拆分成一些更小的可划分的系统，而我们基本都是负责一个或者某几个项目，这些划分出来的更小单位的系统负责不同的功能，但是其底层的组件肯定都是一致的，比如：mysql、redis、kafka相关的引入代码，但是事实的情况却是就连这些底层的引入实现，每个人都有每个是不同的实现方式。这样就会出现同一个系统的底层基础组件的实现却大相径庭，从而导致项目中的一些业务代码对基础组件的引用也不一样，从而感觉每个服务都有每个服务不一样的写法。

这样的场景在微服务的场景下是十分恶心的，所以针对于gin框架，我做了一些扩展封装，将一些底层的基础组件实现进行统一管理，针对mysql、redis、mongo等都选用了主流的框架，这对于大多数的开发来说都是实用的。

还有一个非常不好的开发习惯，就是很多同学喜欢把项目中的一些基础共用的数据或者配置信息直接在代码里面写死，这个是非常不理智的行为，当出现一些服务环境的之间的切换，数据变更的时候，你需要去频繁修改代码，所以我更建议将一些固定的数据、统一收归到配置文件中去。 实现配置的方式有很多，比如在项目中新建一个配置文件，或者通过mysql、etcd、consul、zk等存储，当然直接使用etcd、consul、zk存储配置也会遇到很多比较麻烦的问题，比如版本切换，不同环境之间的配置值不一致等，你可以直接使用配置中心进行接入解决[配置中心](https://github.com/limeschool/devops)(目前内测中)。

### 基于配置编程
在gin的框架的基础上，我扩展实现了很多组件的封装，比如链路请求，链路追踪，统一配置中心，ip限流，自适应限流，mysql、mongo、redis、rbac等组件的封装，接下来我们开始学习如何使用他们。
### 快速入门
##### gin安装
```
go get -u github.com/limeschool/gin
```

##### 创建一个配置文件 config/dev.json
```
{
  "service": "服务名"
}
```

##### 创建第一个web应用
首先打开你的编辑器，创建一个main.go

```
package main

import "github.com/limeschool/gin"

func main() {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.Run(":8080")
}
```
##### 运行程序
```
go run main.go
```
当然如果你的配置文件不是config/dev.json 这个路径，你也可以在运行时通过参数进行指定。
```
go run main.go -c config/conf.json
```
##### 访问web应用
浏览器输入 localhost:8080/ping 你可以看到浏览器输出 pong 。


### mysql使用

##### config/dev.json 配置
```
{
  "service": "user-center",
  "mysql": [
    {
    "enable":true,
    "name":"ums",
    "dsn": "root:root@tcp(127.0.0.1:3306)/devops_ums?charset=utf8mb4&parseTime=True&loc=Local"
    }
  ]
}
```

#### mysql配置解析

|参数|说明|类型|是否必填|默认值|
|-|-|-|-|-|
|enable|是否启用|bool|true|false|
|name|数据库标志符，唯一标志数据库|string|true|-|
|dsn|连接dsn|string|true|-|
|conn_max_lifetime|连接最大存活时长|int|false|120(秒)|
|max_open_conn|最大连接数|int|false|10|
|max_idle_conn|最大空闲连接数|int|false|5|
|level|日志等级 1:Silent 2:Error 3:Warn 4:Info|int|false|4|
|slow_threshold|慢查询阈值|int|false|2(秒)|
|table_prefix|表前缀|string|false|-|
|skip_default_transaction|跳过默认事务|bool|false|true|
|singular_table|使用单数表名|bool|false|true| 
|dry_run|测试生成的 SQL|bool|false|false|
|prepare_stmt|执行任何 SQL 时都会创建一个 prepared statement 并将其缓存|bool|false|false|
|disable_foreign_key|GORM 会自动创建外键约束|bool|false|false|

#### 使用示范

```
package main

import "github.com/limeschool/gin"

func main() {
	e := gin.Default()
	e.GET("/test", func(ctx *gin.Context) {
		db := ctx.Mysql("ums")
		// db操作
	})
	e.Run(":8000")
}
```

### mongo 使用

##### config/dev.json 配置
```
{
  "service": "user-center",
  "mongo": [
    {
    "enable":true,
    "name":"ums",
    "dsn": "root:root@tcp(127.0.0.1:3306)/devops_ums?charset=utf8mb4&parseTime=True&loc=Local"
    }
  ]
}
```


#### mongo配置解析

|参数|说明|类型|是否必填|默认值|
|-|-|-|-|-|
|enable|是否启用|bool|true|false|
|name|数据库标志符，唯一标志数据库|string|true|-|
|dsn|连接dsn|string|true|-|
|min_pool_size|最小连接数|int|false|0|
|max_pool_size|最大连接数|int|true|10|

#### 使用示范

```
package main

import "github.com/limeschool/gin"

func main() {
	e := gin.Default()
	e.GET("/test", func(ctx *gin.Context) {
		db := ctx.Mongo("ums")
		// db操作
	})
	e.Run(":8000")
}
```


### redis 使用
##### config/dev.json 配置
```
{
  "service": "user-center",
  "redis": [
    {
    "enable":true,
    "name":"redis",
    }
  ]
}
```

#### reid配置解析

|参数|说明|类型|是否必填|默认值|
|-|-|-|-|-|
|enable|是否启用|bool|true|false|
|name|数据库标志符，唯一标志数据库|string|true|-|
|host|连接地址|string|true|-|
|network|网路类型|string|false|-|
|username|用户名|string|false|""|
|password|用户密码|string|false|""|
|db|连接指定数据库|int|false|0|
|pool_size|连接池数量|int|false|cpu核*10|



#### 使用示范

```
package main

import "github.com/limeschool/gin"

func main() {
	e := gin.Default()
	e.GET("/test", func(ctx *gin.Context) {
		redis := ctx.Redis("redis")
		// redis操作
	})
	e.Run(":8000")
}
```



### rsa 使用
##### config/dev.json 配置
```
{
  "service": "user-center",
  "rsa": [
    {
    "enable":true,
    "name":"redis",
    "path":"config/private.key",
    }
  ]
}
```

#### rsa配置解析

|参数|说明|类型|是否必填|默认值|
|-|-|-|-|-|
|enable|是否启用|bool|true|false|
|name|数据库标志符，唯一标志数据库|string|true|-|
|path|证书地址|string|true|-|

#### 使用示范

```
package main

import "github.com/limeschool/gin"

func main() {
	e := gin.Default()
	e.GET("/test", func(ctx *gin.Context) {
		ctx.Rsa("public").Encode("hello world")
		ctx.Rsa("private").Decode("xxx")
	})
	e.Run(":8000")
}
```


### rbac 使用
##### config/dev.json 配置
```
{
  "service": "user-center",
  "rbac": [
    {
    "enable":true,
    "db":"ums"
    }
  ]
}
```

#### rbac配置解析

|参数|说明|类型|是否必填|默认值|
|-|-|-|-|-|
|enable|是否启用|bool|true|false|
|db|数据库标志符，唯一标志数据库,从mysql配置中的name中获取|string|true|-|

#### 使用示范

```
package main

import "github.com/limeschool/gin"

func main() {
	e := gin.Default()
	e.GET("/test", func(ctx *gin.Context) {
		if is, _ := ctx.Rbac().Enforce("admin", "/test", "GET"); !is {
			ctx.JSON(200,gin.H{"msg":"暂无权限"})
		}else{
		    ctx.JSON(200,gin.H{"msg":"验证通过"})
		}
	})
	e.Run(":8000")
}
```


### request 链路请求使用
##### config/dev.json 配置
```
{
  "service": "user-center",
  "request":{
    "enable_log": true,
    "retry_count": 3,
    "retry_wait_time":"1s",
    "timeout": "10s",
    "request_msg": "http request",
    "response_msg": "http response"
  }
}
```

#### request配置解析

|参数|说明|类型|是否必填|默认值|
|-|-|-|-|-|
|enable_log|是否启用日志|bool|false|true|
|retry_count|请求失败重试次数|int|false|3|
|retry_wait_time|重试等待时长|duration|false|2s|
|max_retry_wait_time|最大重试等待时长|duration|false|4s|
|request_msg|请求数据标识|string|false|request|
|response_msg|返回数据标识|string|false|response|

#### 使用示范

```
package main

import "github.com/limeschool/gin"

func main() {
	e := gin.Default()
	e.GET("/test", func(ctx *gin.Context) {
		var data map[string]interface{}
		request := ctx.Http().Option(func(request *resty.Request) *resty.Request {
			request.Header.Set("Token", "123")
			return request
		}).Get("https://baidu.com")
		if err := request.Result(&data); err != nil {
			fmt.Println(err)
		}
	})
	e.Run(":8000")
}
```


### config 统一配置使用
##### config/dev.json 配置
```
{
  "service": "user-center",
  "tx_appid":"xxxx",
  "tx_key":"xxxx"
}
```

#### 使用示范

```
package main

import "github.com/limeschool/gin"

func main() {
	e := gin.Default()
	e.GET("/test", func(ctx *gin.Context) {
		appid := ctx.Config.Get("tx_appid")
		fmt.Println(appid)
	})
	e.Run(":8000")
}
```


### 链路日志使用
##### config/dev.json 配置
```
{
  "service": "user-center",
  "log": {
    "level": 0,
    "trace_key": "log-id",
    "service_key": "service"
  }
}
```

#### 配置解析
|参数|说明|类型|是否必填|默认值|
|-|-|-|-|-|
|level|日志等级|int|false|0|
|trace_key|链路id的标志符|string|false|trace-id|
|service_key|服务名标志符|string|false|service|

#### 使用示范

```
package main

import "github.com/limeschool/gin"

func main() {
	e := gin.Default()
	e.GET("/test", func(ctx *gin.Context) {
		ctx.Log.Info("hello")
		ctx.Log.Info("hello2")
	})
	e.Run(":8000")
}
```



### 系统配置使用
##### config/dev.json 配置
```
{
  "service": "user-center",
  "system": {
    "client_timeout":{
      "enable": false,
      "timeout": "10s"
    },
    "client_limit": {
      "enable": true,
      "threshold": 100
    },
    "cpu_threshold": {
      "enable": true,
      "threshold": 900
    }
  }
}
```

#### 
|参数|说明|类型|是否必填|默认值|
|-|-|-|-|-|
|client_timeout|请求超时设置|duration|false|-|
|client_limit|ip限流|int|false|-|
|cpu_threshold|cpu自适应降载 取值建议为100-1000 ，建议为900|int|false|-|

#### 使用示范
配置之后自动生效
