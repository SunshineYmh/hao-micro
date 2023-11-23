package config

import (
	"errors"
	"fmt"
	"os"
	"path"

	"github.com/spf13/viper"
)

type Config struct {
	Service ServerConfig `json:"service"`
	Consul  consulConfig `json:"consul"`
	MySql   DbConfig     `json:"mysql"`
	Redis   RedisConfig  `json:"redis"`
	Jwt     JwtConfig    `json:"jwt"`
}

type ServerConfig struct {
	HaoMicro HaoMicro `json:"haomicro"`
	HaoWeb   HaoWeb   `json:"haoweb"`
}

type HaoMicro struct {
	Port int `json:"port"`
}

type HaoWeb struct {
	Port int `json:"port"`
}

// consul  服务配置
type consulConfig struct {
	Address    string `json:"address"`
	TimeTicker int    `json:"timeTicker"`
}

func SyConfig() (Config, error) {
	// 获取项目目录
	workDir, _ := os.Getwd()
	config := Config{}

	v := viper.New()
	v.SetConfigFile(path.Join(workDir, "config.yaml"))
	v.SetConfigType("yaml")
	if err := v.ReadInConfig(); err != nil {
		massage := fmt.Sprintf("配置文件读取失败: %s ", err.Error())
		return config, errors.New(massage)
	}

	if err := v.Unmarshal(&config); err != nil {
		massage := fmt.Sprintf("解析结构体失败: %s ", err.Error())
		return config, errors.New(massage)
	}

	fmt.Println("api网关端口： ", config.Service.HaoMicro.Port)
	fmt.Println("webApi 端口： ", config.Service.HaoWeb.Port)
	fmt.Println("consul  服务地址： ", config.Consul.Address)
	return config, nil
}
