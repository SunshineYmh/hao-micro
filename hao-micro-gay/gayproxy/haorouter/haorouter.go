package haorouter

import (
	"hao-micro/hao-micro-gay/gayproxy/filters"

	"github.com/gin-gonic/gin"
)

// 初始化路由
func IntoRouter(proxy_addr string) error {
	gayproxy_router := gin.Default()

	//gayproxy_router.Use(haoProxyRecover(proxy_addr))
	// gayproxy_router.Use(httpServiceProxy())
	gayproxy_router.Use(filters.Filter())

	err := gayproxy_router.Run(proxy_addr)
	return err
}
