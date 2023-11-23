package config

import (
	"hao-micro/hao-micro-gateway/utils/common"
	"time"

	"github.com/gomodule/redigo/redis"
)

// redis配置
// @author Zhiqiang Guo
var Redis = map[string]string{
	"name":    "redis",
	"type":    "tcp",
	"address": "127.0.0.1:6379",
	"auth":    "",
}

// const (
// 	CACHE_HAO_MICRO_TOKEN = "user.token."
// 	config_BLACK_TOKEN    = "black.token."
// )

type RedisConfig struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	Address     string `json:"address"`
	Auth        string `json:"auth"`
	MaxIdle     int    `json:"maxIdle"`
	MaxActive   int    `json:"maxActive"`
	IdleTimeout int    `json:"idleTimeout"`
	Expires     int    `json:"expires"`
	RedisSuf    string `json:"redisSuf"`
	RedisBlack  string `json:"redisBlack"`
}

type RedisHandler interface {
	Init_redis_cil()
}

var RedisClient *redis.Pool
var RedisExpire = 86400 * 7
var RedisSuf = "hao.mircro.redis.suf"
var RedisBlack = "hao.mircro.redis.black"

func (rec RedisConfig) Init_redis_cil() {
	// 建立连接池
	RedisClient = &redis.Pool{
		// 从配置文件获取maxidle以及maxactive，取不到则用后面的默认值
		MaxIdle: rec.MaxIdle, //最初的连接数量 16
		// MaxActive:1000000,    //最大连接数量
		MaxActive:   rec.MaxActive,                                //连接池最大连接数量,不确定可以用0（0表示自动定义），按需分配
		IdleTimeout: time.Duration(rec.IdleTimeout) * time.Second, //连接关闭时间 300秒 （300秒不使用自动关闭）
		Dial: func() (redis.Conn, error) { //要连接的redis数据库
			c, err := redis.Dial(rec.Type, rec.Address)
			if err != nil {
				return nil, err
			}
			if rec.Auth != "" {
				if _, err := c.Do("AUTH", rec.Auth); err != nil {
					_ = c.Close()
					return nil, err
				}
			}
			return c, nil
		},
	}
	RedisExpire = common.SetDefaultInt(rec.Expires, 0, 300)
	RedisSuf = common.SetDefault(rec.RedisSuf, "", "hao.mircro.redis.suf")
	RedisBlack = common.SetDefault(rec.RedisBlack, "", "hao.mircro.redis.black")
}

func SetEXData(key, value string, expires int) (err error) {
	key = RedisSuf + RedisSuf + key
	// 从池里获取连接
	rc := RedisClient.Get()
	// 用完后将连接放回连接池
	defer rc.Close()
	_, err = rc.Do("Set", key, value, "EX", expires)
	if err != nil {
		return
	}
	return
}

func GetEXData(key string) (string, error) {
	key = RedisSuf + RedisSuf + key
	// 从池里获取连接
	rc := RedisClient.Get()
	// 用完后将连接放回连接池
	defer rc.Close()
	val, err := redis.String(rc.Do("GET", key))
	if err != nil {
		return "", err
	}
	return val, nil
}

func DelKey(key string) error {
	key = RedisSuf + RedisSuf + key
	// 从池里获取连接
	rc := RedisClient.Get()
	// 用完后将连接放回连接池
	defer rc.Close()
	_, err := rc.Do("DEL", key)
	if err != nil {
		return err
	}
	return nil
}

// 加入到黑名单
func AddBlack(key, token string, expires int) (err error) {
	key = RedisSuf + RedisBlack + key
	// 从池里获取连接
	rc := RedisClient.Get()
	// 用完后将连接放回连接池
	defer rc.Close()
	_, err = rc.Do("Set", key, token, "EX", expires)
	if err != nil {
		return
	}
	return
}

// 检查token是否在黑名单
func CheckBlack(key, token string) bool {
	key = RedisSuf + RedisBlack + key
	// 从池里获取连接
	rc := RedisClient.Get()
	// 用完后将连接放回连接池
	defer rc.Close()
	val, err := redis.String(rc.Do("GET", key))
	if err != nil || val != token {
		return false
	}
	return true
}
