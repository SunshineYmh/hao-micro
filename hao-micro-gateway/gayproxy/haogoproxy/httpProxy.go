package haogoproxy

import (
	"crypto/tls"
	"errors"
	"fmt"
	"hao-micro/hao-micro-gay/gayproxy/haoType"
	"hao-micro/hao-micro-gay/gayproxy/limit"
	"hao-micro/hao-micro-gay/utils"
	"net"
	"net/http"
	"time"
)

var HTTP_CIL_PROXY_MAP map[string]HttpCliServieProxyRouter

type HttpCliServieProxyRouter struct {
	HaorouterConfing haoType.HaorouterConfing
	Client           *http.Client
	TokenBucket      *limit.TokenBucket // 限流
	// Lb               *loadbalance.LeastConnectionBalancer // 负载均衡
}

func IntoHttpServieProxy() {
	// 初始化 map
	HTTP_CIL_PROXY_MAP = make(map[string]HttpCliServieProxyRouter)
}

func AddHttpCilServieProxy(haorouterConfing haoType.HaorouterConfing) utils.Result {
	client, err := HttpServieClientProxy(haorouterConfing)
	if err != nil {
		return utils.NewErrorResult(500, err.Error())
	}
	// // 创建负载均衡器实例
	// proxy_url := haorouterConfing.ProxyUrls
	// lb := loadbalance.NewLeastConnectionBalancer()
	// for i := range proxy_url {
	// 	lb.Add(proxy_url[i])
	// }

	HTTP_CIL_PROXY_MAP[haorouterConfing.Path] = HttpCliServieProxyRouter{
		HaorouterConfing: haorouterConfing,
		Client:           client,
		TokenBucket: &limit.TokenBucket{
			Capacity:  haorouterConfing.LimitCapacity,
			Rate:      1.0,
			Tokens:    0,
			LastToken: time.Now(),
		},
		// Lb: lb,
	}

	_, ok := HTTP_CIL_PROXY_MAP[haorouterConfing.Path]
	if !ok {
		return utils.NewErrorResult(500, "创建路由失败！")
	}
	return utils.NewSuccessResult(haorouterConfing.Path)
}

func HttpServieClientProxy(proxyinfo haoType.HaorouterConfing) (*http.Client, error) {
	if proxyinfo.DiaTimeout <= 0 {
		proxyinfo.DiaTimeout = 10
	}
	if proxyinfo.DiaKeepAlive <= 0 {
		proxyinfo.DiaKeepAlive = 30
	}
	if proxyinfo.TLSHandshakeTimeout <= 0 {
		proxyinfo.TLSHandshakeTimeout = 10
	}
	if proxyinfo.RespHeaderTimeout <= 0 {
		proxyinfo.RespHeaderTimeout = 10
	}
	if proxyinfo.ExpectContinueTimeout <= 0 {
		proxyinfo.ExpectContinueTimeout = 1
	}
	if proxyinfo.MaxIdleConns <= 0 {
		proxyinfo.MaxIdleConns = 100
	}

	if proxyinfo.MaxIdleConnsPerHost <= 0 {
		proxyinfo.MaxIdleConnsPerHost = 2
	}
	// proxyIndex := 0
	transport := &http.Transport{
		DialContext: (&net.Dialer{
			//该超时参数用于控制建立 TCP 连接的超时时间，即在调用 方法时，如果连接建立时间超过该超时时间，
			//则会返回一个 类型的错误。net.Dialer.DialContext()net.Error 10
			Timeout: time.Duration(proxyinfo.DiaTimeout) * time.Second,
			//保持连接的时间，即当连接处于空闲状态时，等待下一次读写操作的最长时间。
			KeepAlive: time.Duration(proxyinfo.DiaKeepAlive) * time.Second,
		}).DialContext,
		//TLS 握手超时时间，即在进行 TLS 握手时等待的最长时间。
		TLSHandshakeTimeout: time.Duration(proxyinfo.TLSHandshakeTimeout) * time.Second,
		//该超时参数用于控制读取 HTTP 响应头的超时时间， 响应头超时时间，即在读取响应头时等待的最长时间。
		//即在调用 方法时，如果读取响应头的时间超过该超时时间，则会返回一个 类型的错误。
		ResponseHeaderTimeout: time.Duration(proxyinfo.RespHeaderTimeout) * time.Second,
		// 该超时参数用于控制等待客户端发送 Expect：
		ExpectContinueTimeout: time.Duration(proxyinfo.ExpectContinueTimeout) * time.Second,

		//ForceAttemptHTTP2: 这个参数用于指定是否强制使用HTTP/2协议。
		//如果设置为true，则客户端将尝试使用HTTP/2协议与服务器进行通信。
		//如果服务器不支持HTTP/2，则会回退到HTTP/1.1协议。如果设置为false，
		//则客户端将根据服务器的响应来决定使用的协议，默认值为false。
		ForceAttemptHTTP2: proxyinfo.ForceAttemptHTTP2,
		//MaxIdleConns 属性表示连接池中最大的空闲连接数。
		//当连接池中的空闲连接数量达到MaxIdleConns时，多余的空闲连接将被关闭并从连接池中移除。默认值为100。
		MaxIdleConns: proxyinfo.MaxIdleConns,
		// MaxIdleConnsPerHost属性表示每个目标主机最大的空闲连接数。
		//当连接池中针对某个主机的空闲连接数量达到MaxIdleConnsPerHost时，
		//多余的空闲连接将被关闭并从连接池中移除。默认值为2
		MaxIdleConnsPerHost: proxyinfo.MaxIdleConnsPerHost,
		// IdleConnTimeout属性表示空闲连接的超时时间。
		// 当一个连接在一段时间内没有被使用，即处于空闲状态时，IdleConnTimeout定义了连接可以保持空闲的最长时间。
		// 超过这个时间，连接将被关闭并从连接池中移除。默认值为0，表示没有超时限制。
		IdleConnTimeout: time.Duration(proxyinfo.IdleConnTimeout) * time.Second,
		// Proxy: http.ProxyURL(proxyUrlParsed),
		// 实现负载均衡 RoundRobin 轮查法
		// Proxy: func(req *http.Request) (*url.URL, error) {
		// 	proxyUrl, err := url.Parse(proxyinfo.ProxyUrls[proxyIndex])
		// 	if err != nil {
		// 		return nil, err
		// 	}
		// 	proxyIndex = (proxyIndex + 1) % len(proxyinfo.ProxyUrls)
		// 	return proxyUrl, nil
		// },
		//在上述代码中，我们将 http.Transport 的 TLSClientConfig 设置为一个自定义的 tls.Config 对象，
		//并将 InsecureSkipVerify 设置为 true，以跳过证书验证。这样在与目标服务器建立 TLS 连接时，会使用代理服务器进行转发。

		//请注意，跳过证书验证可能会导致安全风险，请根据实际情况谨慎使用。
		//如果您的代理服务器具有有效的证书，可以将 InsecureSkipVerify 设置为 false，以启用证书验证。
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true, // 跳过证书验证
		},
	}

	if proxyinfo.IsTls {
		//"path/to/cert.pem"
		cert, err := tls.LoadX509KeyPair(proxyinfo.CertPath, proxyinfo.CertKey)
		if err != nil {
			msg := fmt.Sprintf("未能加载证书和密钥: %v", err)
			return nil, errors.New(msg)
		}
		transport.TLSClientConfig = &tls.Config{
			Certificates:       []tls.Certificate{cert},
			InsecureSkipVerify: true,
		}
	}

	//创建自定义客户端对象
	client := &http.Client{
		Transport: transport,
		Timeout:   time.Duration(proxyinfo.DiaTimeout) * time.Second, // 设置整个请求的超时时间为 10 秒
	}

	return client, nil

}
