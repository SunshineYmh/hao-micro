package filters

import (
	"hao-micro/hao-micro-gateway/gayproxy/handler"
	"hao-micro/hao-micro-gateway/gayproxy/haolog"
	"hao-micro/hao-micro-gateway/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func Filter() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now() //获取当前时间
		UUID := utils.GetUUID()
		path := c.Request.URL.Path
		c.Request.Header.Set("sys-session-id", UUID)
		loginfo := haolog.LogHelper(UUID, path, c, startTime, http.StatusOK, "")
		httpRecover := handler.HttpRecover{
			Ctx:       c,
			Req:       c.Request,
			UUID:      UUID,
			StartTime: startTime,
			Loginfo:   loginfo,
		}

		//获取处理路由接口
		var httpHandler handler.HttpHandler = httpRecover
		var err error
		httpRecover, err = httpHandler.ReqRouterCheck()
		if err != nil {
			haolog.ErrorLogHelpe(httpRecover)
			c.JSON(httpRecover.StatusCode, utils.NewErrorResult(httpRecover.StatusCode, err.Error()))
			c.Abort() // 中止后续处理器函数的执行
			return
		}

		//限流
		httpHandler = httpRecover
		httpRecover, err = httpHandler.LimitHandler()
		if err != nil {
			haolog.ErrorLogHelpe(httpRecover)
			c.JSON(httpRecover.StatusCode, utils.NewErrorResult(httpRecover.StatusCode, err.Error()))
			c.Abort() // 中止后续处理器函数的执行
			return
		}

		// 请求报文处理
		httpHandler = httpRecover
		httpRecover, err = httpHandler.ReqBody()
		if err != nil {
			haolog.ErrorLogHelpe(httpRecover)
			c.JSON(httpRecover.StatusCode, utils.NewErrorResult(httpRecover.StatusCode, err.Error()))
			c.Abort() // 中止后续处理器函数的执行
			return
		}

		// 设置请求头
		httpHandler = httpRecover
		httpRecover, err = httpHandler.ReqHeader()
		if err != nil {
			haolog.ErrorLogHelpe(httpRecover)
			c.JSON(httpRecover.StatusCode, utils.NewErrorResult(httpRecover.StatusCode, err.Error()))
			c.Abort() // 中止后续处理器函数的执行
			return
		}
		//打印请求数据
		haolog.LogHelpePrintln(httpRecover, "请求数据", false)

		//代理转发
		httpHandler = httpRecover
		httpRecover, err = httpHandler.HttpClient()
		if err != nil {
			haolog.ErrorLogHelpe(httpRecover)
			c.JSON(httpRecover.StatusCode, utils.NewErrorResult(httpRecover.StatusCode, err.Error()))
			c.Abort() // 中止后续处理器函数的执行
			return
		}
		// 获取响应数据
		httpHandler = httpRecover
		httpRecover, err = httpHandler.RespBody()
		if err != nil {
			haolog.ErrorLogHelpe(httpRecover)
			c.JSON(httpRecover.StatusCode, utils.NewErrorResult(httpRecover.StatusCode, err.Error()))
			c.Abort() // 中止后续处理器函数的执行
			return
		}

		// 获取响应数据
		httpHandler = httpRecover
		httpRecover, err = httpHandler.RespHeader()
		if err != nil {
			haolog.ErrorLogHelpe(httpRecover)
			c.JSON(httpRecover.StatusCode, utils.NewErrorResult(httpRecover.StatusCode, err.Error()))
			c.Abort() // 中止后续处理器函数的执行
			return
		}
		//打印响应数据
		haolog.LogHelpePrintln(httpRecover, "响应数据", true)
		c.Next()
	}
}

// func LimitHandler(maxConn int) gin.HandlerFunc {
//     tb := &limit.TokenBucket{
// 		Capacity:  100,
// 		Rate:      1.0,
// 		Tokens:    0,
// 		LastToken: time.Now(),
// 	}
//     return func(c *gin.Context) {
//         if !tb.Allow() {
//             c.String(503, "Too many request")
//             c.Abort()
//             return
//         }
//         c.Next()
//     }
// }
