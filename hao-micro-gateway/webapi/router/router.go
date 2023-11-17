package router

import (
	"bytes"
	"fmt"
	webapiServie "hao-micro/hao-micro-gay/webapi/service"
	"io"

	"github.com/gin-gonic/gin"
)

type ResponseWriterWrapper struct {
	gin.ResponseWriter
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

// middleware app debug log, 作用是记录请求响应的信息
func AppDebugLog() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// record request info
		reqBody, _ := io.ReadAll(ctx.Request.Body)
		// log.Logger.Info("requestInfo",
		// 	log.String("method", ctx.Request.Method),
		// 	log.Any("body", string(reqBody)),
		// 	log.String("clientIP", ctx.Request.RemoteAddr),
		// 	log.String("url", ctx.Request.RequestURI))
		// record response info
		fmt.Println(fmt.Printf("请求报文：%s", string(reqBody)))
		ctx.Request.Body = io.NopCloser(bytes.NewReader(reqBody))

		blw := &ResponseWriterWrapper{Body: bytes.NewBufferString(""), ResponseWriter: ctx.Writer}
		ctx.Writer = blw
		ctx.Next()

		fmt.Println(fmt.Printf("响应报文：%s", blw.Body.String()))

	}
}

/**
* webApi 接口
 */
func IntoWebApi(web_addr string) error {
	web_routes := gin.Default()
	//web_routes.Use(AppDebugLog())
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
