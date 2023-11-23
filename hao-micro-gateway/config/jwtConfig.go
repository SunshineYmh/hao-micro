package config

import (
	"hao-micro/hao-micro-gateway/utils/common"
	"sync"
)

type JwtConfig struct {
	Language string   `json:"language"`
	Token    string   `json:"token"`
	Super    string   `json:"super"`
	RedisPre string   `json:"redisPre"`
	Host     string   `json:"host"`
	OpenJwt  bool     `json:"openJwt"`
	Routes   []string `json:"routes"`
	Expires  int      `json:"expires"`
	JwtKey   string   `json:"jwtKey"`
	// BlackTokenPrefix string   `json:"blackTokenPrefix"`
}

var (
	JWTCfg  JwtConfig
	mutex   sync.Mutex
	declare sync.Once
)

func Set(Cfg JwtConfig) {
	mutex.Lock()
	JWTCfg.RedisPre = common.SetDefault(Cfg.RedisPre, "", "hao.sso.redis")
	JWTCfg.Language = common.SetDefault(Cfg.Language, "", "cn")
	JWTCfg.Token = common.SetDefault(Cfg.Token, "", "token")
	JWTCfg.Super = common.SetDefault(Cfg.Super, "", "admin")               //超级账户
	JWTCfg.Host = common.SetDefault(Cfg.Host, "", "http://localhost:8002") //域名
	JWTCfg.Expires = common.SetDefaultInt(Cfg.Expires, 0, 300)
	JWTCfg.JwtKey = common.SetDefault(Cfg.JwtKey, "", "42wqTE23123wffLU94342wgldgFs")
	// JWTCfg.BlackTokenPrefix = setDefault(Cfg.BlackTokenPrefix, "", "hao.auth.black.token.")
	JWTCfg.Routes = Cfg.Routes
	JWTCfg.OpenJwt = Cfg.OpenJwt
	mutex.Unlock()
}
