package main

import (
	"fmt"
	"hao-micro/hao-micro-gateway/auth"
	"hao-micro/hao-micro-gateway/config"
	"hao-micro/hao-micro-gateway/gayproxy/haogoproxy"
	"hao-micro/hao-micro-gateway/gayproxy/haorouter"
	"hao-micro/hao-micro-gateway/utils"
	webapi_router "hao-micro/hao-micro-gateway/webapi/router"
	"io"
	"os"

	"hao-micro/hao-micro-gateway/consul"

	"github.com/gin-gonic/gin"
)

func main() {
	sysconfig, err := config.SyConfig()
	if err != nil {
		// 将对象格式化为字符串
		message := fmt.Sprintf("启动失败: %s", err)
		panic(message)
	}

	//初始化 mysql
	var db config.DbHandler = sysconfig.MySql
	db.AddMySqlDB()
	//初始化 redis 客户端
	var redis config.RedisHandler = sysconfig.Redis
	redis.Init_redis_cil()

	//初始化consul 客户端
	cousul_err := consul.IntoConsulClient(sysconfig.Consul.Address, sysconfig.Consul.TimeTicker)
	if cousul_err != nil {
		// 将对象格式化为字符串
		message := fmt.Sprintf("初始化consul 客户端失败: %s", err)
		panic(message)
	}
	gin.DisableConsoleColor()
	// Logging to a file.
	f, _ := os.Create("service.log")
	//gin.DefaultWriter = io.MultiWriter(f)
	// 如果需要同时将日志写入文件和控制台，请使用以下代码。
	gin.DefaultWriter = io.MultiWriter(f, os.Stdout)
	// 初始化log
	utils.LogInto()
	//初始UUID
	utils.IntoSnowflake()
	// 初始化 map
	haogoproxy.IntoHttpServieProxy()
	//haogoproxy.IntoRoutersProxyMap()

	//初始化 web 服务，设置路径不拦截
	auth.Load(sysconfig.Jwt)

	// 初始化 代理路由服务
	go func() {
		servicePort := fmt.Sprintf(":%d", sysconfig.Service.HaoMicro.Port)
		err := haorouter.IntoRouter(servicePort)
		if err != nil {
			// 将对象格式化为字符串
			message := fmt.Sprintf("启动网关服务异常: %s", err)
			panic(message)
		} else {
			utils.GayINFO("000000000000", "启动网关服务成功》》》》》")
		}
	}()

	// 启动web-api 服务
	go func() {
		webPort := fmt.Sprintf(":%d", sysconfig.Service.HaoWeb.Port)
		err := webapi_router.IntoWebApi(webPort)
		if err != nil {
			message := fmt.Sprintf("启动webAip服务异常: %s", err)
			panic(message)
		} else {
			utils.GayINFO("000000000001", "启动webApi服务成功》》》》》")
		}
	}()

	// 主goroutine不退出，保持服务器运行
	select {}

}
