package haolog

import (
	"fmt"
	"hao-micro/hao-micro-gateway/gayproxy/handler"
	"hao-micro/hao-micro-gateway/gayproxy/haoType"
	"strings"
	"time"

	dateutils "hao-micro/hao-micro-gateway/utils"
	gaylog "hao-micro/hao-micro-gateway/utils"

	"github.com/gin-gonic/gin"
)

func LogHelper(UUID string, routeId string, c *gin.Context, currentTime time.Time, status int64, errMsg string) haoType.LogInfo {
	r := c.Request
	reqHeaders := make(map[string]string)
	// 遍历请求头
	for key, value := range r.Header {
		reqHeaders[key] = strings.Join(value, ",")
	}
	reqContentType := c.ContentType()
	loginfo := haoType.LogInfo{
		Id:              UUID,                                                      //唯一id
		RouteId:         routeId,                                                   //路由id
		StartDate:       dateutils.TimeFormat(dateutils.Format_YMDHS, currentTime), //请求系统开始时间
		EndDate:         dateutils.TimeFormat(dateutils.Format_YMDHS, currentTime), //请求系统结束时间
		Scheme:          r.URL.Scheme,                                              //请求协议，http/https
		Method:          r.Method,                                                  //请求方式，get/post
		Host:            r.Host,                                                    //请求地址
		Url:             r.RequestURI,                                              //请求地址
		ReqContentType:  reqContentType,                                            //请求数据类型
		RespContentType: "",                                                        //响应数据类型
		Charset:         "UTF-8",                                                   //字符集
		ReqHeader:       reqHeaders,                                                //请求头信息
		RespHeader:      make(map[string]string),                                   //响应头信息
		ReqBody:         "",                                                        //请求报文
		RespBody:        "",                                                        //响应报文
		ReqBodySize:     c.Request.ContentLength,                                   //请求报文大小
		RespBodySize:    0,                                                         //响应报文大小
		ClientIP:        c.ClientIP(),                                              //请求ip
		ExecuteTime:     0,                                                         //执行时间
		Status:          status,                                                    //服务状态
		ErrorMsg:        errMsg,                                                    //错误信息
	}
	gaylog.GayINFO(UUID, fmt.Sprintf("客户端请求开始: %#v", loginfo))
	return loginfo
}

func ErrorLogHelpe(hrc handler.HttpRecover) {
	hrc.Loginfo.Status = int64(hrc.StatusCode)
	hrc.Loginfo.ErrorMsg = hrc.ErrorMassage
	endTime := time.Now() //获取当前时间
	hrc.Loginfo.ExecuteTime = dateutils.DiffTime(hrc.StartTime, endTime)
	hrc.Loginfo.EndDate = dateutils.TimeFormat(dateutils.Format_YMDHS, endTime)
	gaylog.GayERROR(hrc.UUID, fmt.Sprintf("客户端异常信息：%#v", hrc.Loginfo))
}

func LogHelpePrintln(hrc handler.HttpRecover, logNmae string, isTask bool) {
	hrc.Loginfo.Status = int64(hrc.StatusCode)
	hrc.Loginfo.ErrorMsg = hrc.ErrorMassage
	endTime := time.Now() //获取当前时间
	hrc.Loginfo.ExecuteTime = dateutils.DiffTime(hrc.StartTime, endTime)
	hrc.Loginfo.EndDate = dateutils.TimeFormat(dateutils.Format_YMDHS, endTime)
	//打印响应数据
	gaylog.GayINFO(hrc.UUID, fmt.Sprintf("%s: %#v", logNmae, hrc.Loginfo))
	if isTask {
		//结束请求处理任务，日志记录数据库等
	}
}
