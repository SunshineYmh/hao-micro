package config

import (
	"errors"
	"fmt"
	"os"
	"path"

	"github.com/spf13/viper"
)

type Config struct {
	Service ServerConfig `json:"Service"`
	Consul  consulConfig `json:"consul"`
}

type ServerConfig struct {
	ServicePort int `json:"servicePort"`
	WebPort     int `json:"webPort"`
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

	fmt.Println("api网关端口： ", config.Service.ServicePort)
	fmt.Println("webApi 端口： ", config.Service.WebPort)
	fmt.Println("consul  服务地址： ", config.Consul.Address)
	return config, nil
}
