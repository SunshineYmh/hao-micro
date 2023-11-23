package app

import (
	"errors"
	"fmt"
	"hao-micro/hao-micro-gateway/config"

	"github.com/dgrijalva/jwt-go"
)

// const (
// 	SECRETKEY          = "42wqTE23123wffLU94342wgldgFs"
// 	MAXAGE             = 30
// 	config_BLACK_TOKEN = "black.token."
// )

type CustomClaims struct {
	UserId int64
	jwt.StandardClaims
}

// 产生token
func (cc *CustomClaims) MakeToken() (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, cc)
	return token.SignedString([]byte(config.JWTCfg.JwtKey))
}

// 解析token
func ParseToken(accessToken string) (*CustomClaims, error) {
	tokenString, err2 := config.GetEXData(accessToken)
	if err2 != nil {
		msg := fmt.Sprintf("access_token 已失效， %s ！", err2.Error())
		return nil, errors.New(msg)
	}
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			msg := fmt.Sprintf("解析token 失败: %v", token.Header["alg"])
			return nil, errors.New(msg)
		}
		return []byte(config.JWTCfg.JwtKey), nil
	})
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, err
	}
}
