package app

import (
	"hao-micro/hao-micro-gateway/config"
	"hao-micro/hao-micro-gateway/utils/common"
	"hao-micro/hao-micro-gateway/webapi/models"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func DoLogin(c *gin.Context, user models.Users) (map[string]interface{}, error) {
	secure := IsHttps(c)

	data := make(map[string]interface{})
	if config.JWTCfg.OpenJwt { //返回jwt
		customClaims := &CustomClaims{
			UserId: user.Id,
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: time.Now().Add(time.Duration(config.JWTCfg.Expires) * time.Second).Unix(), // 过期时间，必须设置
			},
		}
		accessToken, err := customClaims.MakeToken()
		if err != nil {
			return data, err
		}
		refreshClaims := &CustomClaims{
			UserId: user.Id,
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: time.Now().Add(time.Duration(config.JWTCfg.Expires+1800) * time.Second).Unix(), // 过期时间，必须设置
			},
		}
		refreshToken, err := refreshClaims.MakeToken()
		if err != nil {
			return data, err
		}
		data[ACCESS_TOKEN] = common.GetMd5String(accessToken)
		data[REFRESH_TOKEN] = common.GetMd5String(refreshToken)
		data[EXPIRES_IN] = config.JWTCfg.Expires
		data["secure"] = secure
		data["UserId"] = user.Id

		config.SetEXData(data[ACCESS_TOKEN].(string), accessToken, config.JWTCfg.Expires)
		config.SetEXData(data[REFRESH_TOKEN].(string), refreshToken, config.JWTCfg.Expires+1800)

		c.Header(ACCESS_TOKEN, data[ACCESS_TOKEN].(string))
		c.Header(REFRESH_TOKEN, data[REFRESH_TOKEN].(string))
		c.SetCookie(ACCESS_TOKEN, data[ACCESS_TOKEN].(string), config.JWTCfg.Expires, "/", "", secure, true)
		c.SetCookie(REFRESH_TOKEN, data[REFRESH_TOKEN].(string), config.JWTCfg.Expires, "/", "", secure, true)
	}
	//claims,err:=ParseToken(accessToken)
	//if err!=nil {
	//	return err
	//}
	id := strconv.Itoa(int(user.Id))
	c.SetCookie(COOKIE_TOKEN, id, config.JWTCfg.Expires, "/", "", secure, true)

	return data, nil
}

// 判断是否https
func IsHttps(c *gin.Context) bool {
	if c.GetHeader(HEADER_FORWARDED_PROTO) == "https" || c.Request.TLS != nil {
		return true
	}
	return false
}
