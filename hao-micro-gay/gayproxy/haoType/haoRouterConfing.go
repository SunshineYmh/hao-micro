package haoType

type HaorouterConfing struct {
	RouterId          string `json:"routerId"`          //路由、服务名称ID
	RouterName        string `json:"routerName"`        //路由、服务名称
	ConsulServiceName string `json:"consulServiceName"` //consul服务注册名称
	Path              string `json:"path"`              //服务接口路径（路由）
	Method            string `json:"method"`            //请求方式 GET、 POST
	//	ProxyUrls             []string `json:"proxyUrls"`             //代理地址
	StripPrefix           int    `json:"stripPrefix"`           //表示过滤1个路径，2表示两个，一次类推
	DiaTimeout            int    `json:"diaTimeout"`            //TCP 连接的超时时间
	DiaKeepAlive          int    `json:"diaKeepAlive"`          //空闲时间
	TLSHandshakeTimeout   int    `json:"tLSHandshakeTimeout"`   //TLS 握手超时时间
	RespHeaderTimeout     int    `json:"respHeaderTimeout"`     //http 响应头的超时时间
	ExpectContinueTimeout int    `json:"rxpectContinueTimeout"` //该超时参数用于控制等待客户端发送 Expect：
	ForceAttemptHTTP2     bool   `json:"forceAttemptHTTP2"`     //是否强制使用HTTP/2协议
	MaxIdleConns          int    `json:"maxIdleConns"`          //最大空闲连接数，多余的连接将会被关闭，默认值为100。
	MaxIdleConnsPerHost   int    `json:"maxIdleConnsPerHost"`   //用于限制每个主机的最大空闲连接数，默认值为2
	IdleConnTimeout       int    `json:"idleConnTimeout"`       //空闲连接的超时时间，默认值为0，表示没有超时限制
	LimitCapacity         int64  `json:"limitCapacity"`         //限流，每秒最大访问量，默认值0，表示不限流
	IsTls                 bool   `json:"isTls"`                 //是否启用证书访问
	CertPath              string `json:"certPath"`              //证书路径
	CertKey               string `json:"certKey"`               //证书key路径
}

/**
openssl pkcs12 -in 77928742.p12 -out 77928742_cli.crt -clcerts -nokeys -passin pass:clientatwasoft

使用openssl命令导出.key

openssl pkcs12 -in 77928742.p12 -out 77928742_cli.key -nocerts -nodes -passin pass:clientatwasoft

使用openssl命令导出.csr

openssl pkcs12 -in keystore.p12 -nokeys -out my_key_store.csr

如果需要将证书和私钥合并成一个 PEM 文件，可以使用以下命令：
cat certificate.crt private.key > certificate.pem
现在，你应该有一个包含证书和私钥的 PEM 文件 certificate.pem，可以在你的 Go 代码中使用它了。
*/
