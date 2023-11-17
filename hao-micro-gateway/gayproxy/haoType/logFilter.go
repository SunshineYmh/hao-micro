package haoType

type LogInfo struct {
	Id              string            `json:"id"`              //唯一id
	RouteId         string            `json:"routeId"`         //路由id
	StartDate       string            `json:"startDate"`       //请求系统开始时间
	EndDate         string            `json:"endDate"`         //请求系统开始时间
	Scheme          string            `json:"scheme"`          //请求协议，http/https
	Method          string            `json:"method"`          //请求方式，get/post
	Host            string            `json:"host"`            //请求地址
	Url             string            `json:"url"`             //请求地址
	ProxyUrl        string            `json:"proxyUrl"`        //代理地址
	ReqContentType  string            `json:"reqContentType"`  //请求数据类型
	RespContentType string            `json:"respContentType"` //响应数据类型
	Charset         string            `json:"charset"`         //字符集
	ReqHeader       map[string]string `json:"reqHeader"`       //请求头信息
	RespHeader      map[string]string `json:"respHeader"`      //响应头信息
	ReqBody         string            `json:"reqBody"`         //请求报文
	RespBody        string            `json:"respBody"`        //响应报文
	ReqBodySize     int64             `json:"reqBodySize"`     //请求报文大小
	RespBodySize    int64             `json:"respBodySize"`    //响应报文大小
	ClientIP        string            `json:"clientIP"`        //请求ip
	ExecuteTime     int64             `json:"executeTime"`     //执行时间
	Status          int64             `json:"status"`          //服务状态
	ErrorMsg        string            `json:"rrrorMsg"`        //错误信息
}
