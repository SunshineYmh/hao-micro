package auth

import (
	"fmt"
	"hao-micro/hao-micro-gateway/config"
	"hao-micro/hao-micro-gateway/utils"
	"hao-micro/hao-micro-gateway/utils/common"
	"hao-micro/hao-micro-gateway/utils/handle"
	"hao-micro/hao-micro-gateway/utils/request"
	"hao-micro/hao-micro-gateway/webapi/modules/app"
	"net/http"
	"net/url"
	"strconv"

	"github.com/gin-gonic/gin"
)

func Load(jwt config.JwtConfig) {
	// c := config.JwtConfig{}
	// permitRoutes := strings.Split(permitRoute, ",")
	// c.Routes = []string{"/hao-web/aouth"}
	// for _, v := range permitRoutes {
	// 	fmt.Println("-->>>ssss>   " + v)
	// // }
	// c.Routes = permitRoutes
	// c.OpenJwt = true //开启jwt
	config.Set(jwt)
	//初始化数据验证
	handle.InitValidate()
}

func Auth(c *gin.Context) {
	u, err := url.Parse(c.Request.RequestURI)
	if err != nil {
		panic(err)
	}
	if common.InArrayStringHasPrefix(u.Path, &config.JWTCfg.Routes) {
		c.Next()
		return
	}
	//开启jwt
	if config.JWTCfg.OpenJwt {
		accessToken, has := request.GetParam(c, app.ACCESS_TOKEN)
		if !has {
			c.Abort() //组织调起其他函数
			msg := fmt.Sprintf("获取请求头[%s]为空！", app.ACCESS_TOKEN)
			c.JSON(http.StatusExpectationFailed, utils.NewErrorResult(http.StatusExpectationFailed, msg))
			c.Abort() // 中止后续处理器函数的执行
			return
		}
		ret, err := app.ParseToken(accessToken)
		if err != nil {
			c.Abort()
			c.JSON(http.StatusUnauthorized, utils.NewErrorResult(http.StatusUnauthorized, err.Error()))
			return
		}
		uid := strconv.FormatInt(ret.UserId, 10)
		has = config.CheckBlack(uid, accessToken)
		if has {
			c.Abort() //组织调起其他函数
			msg := fmt.Sprintf("请求非法，[%s]已禁止访问！", app.ACCESS_TOKEN)
			c.JSON(http.StatusForbidden, utils.NewErrorResult(http.StatusForbidden, msg))
			return
		}
		c.Set("uid", ret.UserId)
		c.Next()
		return
	}
	//cookie
	_, err = c.Cookie(app.COOKIE_TOKEN)
	if err != nil {
		c.Abort() //组织调起其他函数
		c.JSON(http.StatusUnauthorized, utils.NewErrorResult(http.StatusUnauthorized, "请求非法，Cookie 无效！"))
		return
	}
	c.Next()
}
