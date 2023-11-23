package service

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"hao-micro/hao-micro-gateway/gayproxy/haoType"
	"hao-micro/hao-micro-gateway/gayproxy/haogoproxy"
	"net/http"
	"net/http/httputil"

	syslog "hao-micro/hao-micro-gateway/utils"

	"github.com/gin-gonic/gin"
)

func AddRouter(c *gin.Context) {
	haorouterConfing := haoType.HaorouterConfing{}
	if err := c.ShouldBindJSON(&haorouterConfing); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	syslog.SysINFO("", fmt.Sprintf("请求参数:%+v", haorouterConfing))
	// res := haogoproxy.AddHttpServieProxy(haorouterConfing)
	res := haogoproxy.AddHttpCilServieProxy(haorouterConfing)
	//res := haogoproxy.AddRouterProxy(haorouterConfing)
	syslog.SysINFO("", fmt.Sprintf("创建网关服务信息: %+v", res))
	c.JSON(200, res)
}

func Httpcli(c *gin.Context) {
	// 创建一个HTTP客户端并发送请求
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}

	client := &http.Client{
		Transport: tr,
	}

	data := map[string]string{
		"xingming": "贲没机",
		"blqd":     "app_12329",
		"appid":    "20211126000201",
		"sign":     "6d75ec9d81c44eaa305cf52e843991dd7576a9b2",
		"citybm":   "C65020",
		"zjhm":     "654123198502011502",
		"userId":   "sjgjj_123123123123",
	}

	jsonData, errjs := json.Marshal(data)
	if errjs != nil {
		fmt.Println("--req->errjserrjs>>", errjs)
	}

	fmt.Println("json>>>>>>", string(jsonData))

	// resp, errff := client.Get("https://dog.ceo/api/breeds/image/random")

	req, err := http.NewRequest("POST", "https://appcs.sjgjj.cn/app-web/public/zhcx/info.service", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("--req->>>", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, errff := client.Do(req)
	if errff != nil {
		fmt.Println("--req-ccc>>>", errff)
	}

	defer resp.Body.Close()

	respData, err := httputil.DumpResponse(resp, true)
	if err != nil {
		fmt.Println("Error dumping request:", err)
	}

	fmt.Printf("响应数据 %s \n", string(respData))

	// body, rrr := io.ReadAll(resp.Body)
	// if rrr != nil {
	// 	fmt.Println("--req-fff>>>", rrr)
	// }
	// fmt.Println("-->00>>>", resp.StatusCode)
	// fmt.Println("--rep>>>>", string(body))
	c.JSON(resp.StatusCode, respData)

}
