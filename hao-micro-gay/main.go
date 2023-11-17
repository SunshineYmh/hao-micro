package main

import (
	"fmt"
	"hao-micro/hao-micro-gay/config"
	"hao-micro/hao-micro-gay/gayproxy/haogoproxy"
	"hao-micro/hao-micro-gay/gayproxy/haorouter"
	"hao-micro/hao-micro-gay/utils"
	webapi_router "hao-micro/hao-micro-gay/webapi/router"
	"io"
	"os"

	"hao-micro/hao-micro-gay/consul"

	"github.com/gin-gonic/gin"
)

func main() {
	config, err := config.SyConfig()
	if err != nil {
		// 将对象格式化为字符串
		message := fmt.Sprintf("启动失败: %s", err)
		panic(message)
	}

	//初始化consul 客户端
	cousul_err := consul.IntoConsulClient(config.Consul.Address, config.Consul.TimeTicker)
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
	// 初始化 代理路由服务
	go func() {
		servicePort := fmt.Sprintf(":%d", config.Service.ServicePort)
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
		webPort := fmt.Sprintf(":%d", config.Service.WebPort)
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
