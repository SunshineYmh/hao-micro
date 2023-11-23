package service

import (
	"hao-micro/hao-micro-gateway/consul"
	"hao-micro/hao-micro-gateway/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

// 注册服务
func ConsulServiceRegister(c *gin.Context) {
	hst := consul.HaoServiceRegistration{}
	if err := c.ShouldBindJSON(&hst); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var consulService consul.ConsulService = hst
	err := consulService.ConsulServiceRegister()
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.NewErrorResult(http.StatusInternalServerError, err.Error()))
	} else {
		c.JSON(http.StatusOK, utils.NewErrorResult(http.StatusOK, "注册服务成功！"))
	}
}

// 注销服务
func ConsulServiceDeregister(c *gin.Context) {
	hst := consul.HaoServiceRegistration{}
	if err := c.ShouldBindJSON(&hst); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var consulService consul.ConsulService = hst
	err := consulService.ConsulServiceDeregister()
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.NewErrorResult(http.StatusInternalServerError, err.Error()))
	} else {
		c.JSON(http.StatusOK, utils.NewSuccessResult(hst.Id))
	}
}

// 服务获取
func ConsulServiceQuery(c *gin.Context) {
	hst := consul.HaoServiceRegistration{}
	if err := c.ShouldBindJSON(&hst); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var consulService consul.ConsulService = hst
	hsts, err := consulService.ConsulServiceQuery()
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.NewErrorResult(http.StatusInternalServerError, err.Error()))
	} else {
		c.JSON(http.StatusOK, utils.NewSuccessResult(hsts))
	}
}
