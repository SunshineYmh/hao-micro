package service

import (
	"fmt"
	"hao-micro/hao-micro-gateway/config"
	"hao-micro/hao-micro-gateway/utils"
	"hao-micro/hao-micro-gateway/utils/common"
	"hao-micro/hao-micro-gateway/utils/handle"
	"hao-micro/hao-micro-gateway/utils/request"
	"hao-micro/hao-micro-gateway/webapi/models"
	"hao-micro/hao-micro-gateway/webapi/modules/app"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type UserMobile struct {
	Mobile string `form:"mobile" json:"mobile" binding:"required"`
	Passwd string `form:"passwd" json:"passwd" binding:"required,max=20,min=6"`
	Code   string `form:"code" json:"code" binding:"required,len=6"`
}

type UserMobilePasswd struct {
	Mobile string `form:"mobile" json:"mobile" binding:"required"`
	Passwd string `form:"passwd" json:"passwd" binding:"required,max=20,min=6"`
}

var UserMobileTrans = map[string]string{"Mobile": "手机号", "Passwd": "密码", "Code": "验证码"}

// 手机密码
func Login(c *gin.Context) {
	var userMobile UserMobilePasswd
	if err := c.BindJSON(&userMobile); err != nil {
		msg := handle.TransTagName(&UserMobileTrans, err)
		c.JSON(http.StatusBadRequest, utils.NewErrorResult(http.StatusBadRequest, msg))
		return
	}
	model := models.Users{Mobile: userMobile.Mobile}
	if has := model.GetRow(); !has {
		c.JSON(http.StatusBadRequest, utils.NewErrorResult(http.StatusBadRequest, "mobile_not_exists"))
		return
	}
	if common.Sha1En(userMobile.Passwd+model.Salt) != model.Passwd {
		c.JSON(http.StatusBadRequest, utils.NewErrorResult(http.StatusBadRequest, "login_error"))
		return
	}
	token, err := app.DoLogin(c, model)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.NewErrorResult(http.StatusBadRequest, "fail"))
		return
	}
	c.JSON(http.StatusOK, utils.NewSuccessResult(token))
	// return
}

// 注销登录
func Logout(c *gin.Context) {
	secure := app.IsHttps(c)
	//access_token  refresh_token 加黑名单
	accessToken, has := request.GetParam(c, app.ACCESS_TOKEN)
	if has {
		uid := strconv.FormatInt(c.MustGet("uid").(int64), 10)
		config.AddBlack(uid, accessToken, config.JWTCfg.Expires)
		config.DelKey(accessToken)
	}
	c.SetCookie(app.COOKIE_TOKEN, "", -1, "/", "", secure, true)
	c.SetCookie(app.ACCESS_TOKEN, "", -1, "/", "", secure, true)
	c.SetCookie(app.REFRESH_TOKEN, "", -1, "/", "", secure, true)
	c.JSON(http.StatusOK, utils.NewSuccessResult("success"))
	// return
}

// 手机号注册
func SignupByMobile(c *gin.Context) {
	var userMobile UserMobile
	if err := c.BindJSON(&userMobile); err != nil {
		msg := handle.TransTagName(&UserMobileTrans, err)
		fmt.Println(msg)
		c.JSON(http.StatusBadRequest, utils.NewErrorResult(http.StatusBadRequest, msg))
		return
	}
	model := models.Users{Mobile: userMobile.Mobile}
	if has := model.GetRow(); has {
		c.JSON(http.StatusBadRequest, utils.NewErrorResult(http.StatusBadRequest, "mobile_exists"))
		return
	}
	//验证code
	//if sms.SmsCheck("code"+userMobile.Mobile,userMobile.Code) {
	//	response.ShowError(c, "code_error")
	//	return
	//}

	model.Salt = common.GetRandomBoth(4)
	model.Passwd = common.Sha1En(userMobile.Passwd + model.Salt)
	model.Ctime = int(time.Now().Unix())
	model.Status = models.UsersStatusOk
	model.Mtime = time.Now()

	traceModel := models.Trace{Ctime: model.Ctime}
	traceModel.Ip = common.IpStringToInt(request.GetClientIp(c))
	traceModel.Type = models.TraceTypeReg

	deviceModel := models.Device{Ctime: model.Ctime, Ip: traceModel.Ip, Client: c.GetHeader("User-Agent")}
	_, err := model.Add(&traceModel, &deviceModel)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, utils.NewErrorResult(http.StatusBadRequest, "fail"))
		return
	}
	c.JSON(http.StatusOK, utils.NewSuccessResult("success"))
	// return
}

// access token 续期
func Renewal(c *gin.Context) {
	_, has := request.GetParam(c, app.ACCESS_TOKEN)
	if !has {
		c.JSON(http.StatusBadRequest, utils.NewErrorResult(http.StatusBadRequest, "access token not found"))
		return
	}
	refreshToken, has := request.GetParam(c, app.REFRESH_TOKEN)
	if !has {
		c.JSON(http.StatusBadRequest, utils.NewErrorResult(http.StatusBadRequest, "refresh_token"))

		return
	}
	ret, err := app.ParseToken(refreshToken)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.NewErrorResult(http.StatusBadRequest, "refresh_token"))
		return
	}
	// uid := strconv.FormatInt(ret.UserId, 10)
	model := models.Users{Id: ret.UserId}
	token, err := app.DoLogin(c, model)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.NewErrorResult(http.StatusBadRequest, "fail"))
		return
	}
	c.JSON(http.StatusOK, utils.NewSuccessResult(token))
}

func Info(c *gin.Context) {
	uid := c.MustGet("uid").(int64)
	fmt.Println(uid)
	model := models.Users{}
	model.Id = uid
	row, err := model.GetRowById()
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, utils.NewErrorResult(http.StatusBadRequest, err.Error()))
		return
	}
	fmt.Println(row)
	fmt.Println(row.Name)
	//隐藏手机号中间数字
	s := row.Mobile
	row.Mobile = string([]byte(s)[0:3]) + "****" + string([]byte(s)[6:])
	c.JSON(http.StatusOK, utils.NewSuccessResult(row))
	return
}
