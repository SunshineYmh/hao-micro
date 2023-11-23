package router

import (
	"hao-micro/hao-micro-gateway/auth"
	webapiServie "hao-micro/hao-micro-gateway/webapi/service"

	"github.com/gin-gonic/gin"
)

/**
* webApi 接口
 */
func IntoWebApi(web_addr string) error {

	web_routes := gin.Default()
	gin.SetMode(gin.DebugMode) //开发环境
	// gin.SetMode(gin.ReleaseMode) //线上环境
	web_routes.Use(auth.Auth)
	hao_web_auth_api := web_routes.Group("/hao-web/aouth")
	{
		hao_web_auth_api.POST("/login", webapiServie.Login)
		hao_web_auth_api.POST("/signup/mobile", webapiServie.SignupByMobile)
		hao_web_auth_api.POST("/renewal", webapiServie.Renewal)
	}
	hao_web_user_api := web_routes.Group("/hao-web/user")
	{
		hao_web_user_api.POST("/logout", webapiServie.Logout)
		hao_web_user_api.GET("/my/info", webapiServie.Info)
	}

	hao_web_api := web_routes.Group("/hao-web/gayproxy")
	{
		hao_web_api.POST("/AddRouter", webapiServie.AddRouter)
		hao_web_api.POST("/test", webapiServie.Httpcli)
	}
	hao_web_consul_api := web_routes.Group("/hao-web/consul")
	{
		hao_web_consul_api.POST("/ConsulServiceRegister", webapiServie.ConsulServiceRegister)
		hao_web_consul_api.POST("/ConsulServiceDeregister", webapiServie.ConsulServiceDeregister)
		hao_web_consul_api.POST("/ConsulServiceQuery", webapiServie.ConsulServiceQuery)
	}
	err := web_routes.Run(web_addr)
	return err
}
