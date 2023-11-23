package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"hao-micro/hao-micro-gateway/consul"
	"hao-micro/hao-micro-gateway/gayproxy/haoType"
	"hao-micro/hao-micro-gateway/gayproxy/haogoproxy"
	"hao-micro/hao-micro-gateway/utils"
	gaylog "hao-micro/hao-micro-gateway/utils"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type HttpHandler interface {
	ReqRouterCheck() (HttpRecover, error) //请求路由就检查
	LimitHandler() (HttpRecover, error)   //限流操作
	ReqHeader() (HttpRecover, error)      // 请求头处理
	ReqBody() (HttpRecover, error)        //请求报文处理
	HttpClient() (HttpRecover, error)
	RespBody() (HttpRecover, error)   // 响应报文处理
	RespHeader() (HttpRecover, error) // 响应头处理
}

type HttpRecover struct {
	Ctx                      *gin.Context
	Req                      *http.Request
	Resp                     *http.Response
	UUID                     string
	ConsulServiceAddress     string
	HttpCliServieProxyRouter haogoproxy.HttpCliServieProxyRouter
	StatusCode               int
	ErrorMassage             string
	Loginfo                  haoType.LogInfo
	StartTime                time.Time
	EndTime                  time.Time
}

// 获取路由、客户端
func (hrc HttpRecover) ReqRouterCheck() (HttpRecover, error) {
	path := hrc.Ctx.Request.URL.Path
	httpCliServieProxyRouter, ok := haogoproxy.HTTP_CIL_PROXY_MAP[path]
	gaylog.GayINFO(hrc.UUID, fmt.Sprintf("请求开始:%s , 地址核查:%t", path, ok))
	gaylog.GayINFO(hrc.UUID, fmt.Sprintf("请求地址, %s", hrc.Ctx.Request.RequestURI))
	if !ok {
		msg := "请求地址不存在！"
		hrc.ErrorMassage = msg
		hrc.StatusCode = http.StatusNotFound
		return hrc, errors.New(msg)
	}
	hrc.HttpCliServieProxyRouter = httpCliServieProxyRouter

	hsts, err := consul.GetLeastConnectionInstance(httpCliServieProxyRouter.HaorouterConfing.ConsulServiceName)
	if err != nil {
		msg := "请求consul注册服务不存在！"
		hrc.ErrorMassage = msg
		hrc.StatusCode = http.StatusNotFound
		return hrc, errors.New(msg)
	}
	//获取http 协议
	meta, ok := hsts.Meta["protocol"]
	if !ok {
		meta = "http"
	}
	hrc.ConsulServiceAddress = fmt.Sprintf("%s://%s:%d", meta, hsts.Address, hsts.Port)
	return hrc, nil
}

// 请求限流处理
func (hrc HttpRecover) LimitHandler() (HttpRecover, error) {
	// 如果限流总量大于0，则进行限流计算
	if hrc.HttpCliServieProxyRouter.HaorouterConfing.LimitCapacity > 0 {
		tb := hrc.HttpCliServieProxyRouter.TokenBucket
		if !tb.Allow() {
			msg := "系统繁忙请稍后再试！"
			hrc.ErrorMassage = msg
			hrc.StatusCode = http.StatusServiceUnavailable
			return hrc, errors.New(msg)
		}
	}
	return hrc, nil
}

// 请求头处理
func (hrc HttpRecover) ReqHeader() (HttpRecover, error) {
	reqHeaders := make(map[string]string)
	// 遍历请求头
	for key, value := range hrc.Ctx.Request.Header {
		reqHeaders[key] = strings.Join(value, ",")
	}
	hrc.Loginfo.ReqHeader = reqHeaders
	return hrc, nil
}

// 请求报文处理
func (hrc HttpRecover) ReqBody() (HttpRecover, error) {
	hrc.Loginfo.ReqBodySize = hrc.Ctx.Request.ContentLength
	method := strings.ToLower(hrc.Ctx.Request.Method)
	if hrc.Ctx.Request.ContentLength > 0 {
		var massage string
		switch method {
		case "get":
			// 获取所有的请求参数
			params := hrc.Ctx.Request.URL.Query()
			params_map := make(map[string]string)
			// 遍历参数
			for key, values := range params {
				// 如果参数有多个值，可以使用 range 遍历
				for _, value := range values {
					params_map[key] = value
					//设置请求参数
					params.Set(key, value)
				}
			}
			params.Set("ccc", "我测试下赋值")
			// 将参数设置回URL对象
			hrc.Ctx.Request.URL.RawQuery = params.Encode()
			hrc.Ctx.Request.URL.RawQuery = params.Encode()
			hrc.Ctx.Request.RequestURI = hrc.Ctx.Request.URL.String()
			// 将数组转换为 JSON
			jsonData, _ := json.Marshal(params_map)
			hrc.Loginfo.ReqBody = string(jsonData)
		case "post":
			// 当 expression 的值等于 value2 时执行的代码
			if strings.HasPrefix(strings.ToLower(hrc.Ctx.ContentType()), "multipart/form-data") {
				body, from_ContentType, form_data := FormData(hrc.Ctx)
				// 将数组转换为 JSON
				jsonData, err := json.Marshal(form_data)
				if err != nil {
					massage = fmt.Sprintf("from-data 请求数据处理异常: %s ", err.Error())
					hrc.ErrorMassage = massage
					hrc.StatusCode = http.StatusInternalServerError
					return hrc, errors.New(massage)
				}
				// 打印 JSON 数据
				// fmt.Println("from-data 请求数据：", string(jsonData))
				hrc.Loginfo.ReqContentType = from_ContentType
				hrc.Loginfo.ReqBody = string(jsonData)
				hrc.Ctx.Request.Header.Set("Content-Type", from_ContentType)
				hrc.Req.Header.Set("Content-Type", from_ContentType)
				hrc.Ctx.Request.Body = io.NopCloser(body)
				hrc.Req.Body = io.NopCloser(body)
			} else if strings.HasPrefix(strings.ToLower(hrc.Ctx.ContentType()), "application/x-www-form-urlencoded") {
				// params := hrc.Ctx.Request.PostForm
				params_map := make(map[string]string)
				// 构造请求参数
				reqparams := url.Values{}
				// 绑定请求参数到结构体
				if err := hrc.Ctx.ShouldBind(&params_map); err != nil {
					massage = fmt.Sprintf("  请求数据处理异常: %s ", err.Error())
					hrc.ErrorMassage = massage
					hrc.StatusCode = http.StatusBadRequest
					return hrc, errors.New(massage)
				}
				// 遍历参数
				for key, value := range params_map {
					reqparams.Set(key, value)
				}
				reqparams.Set("x-www-form", "我测试下x-www-form赋值")
				// 将参数设置回URL对象
				hrc.Ctx.Request.Body = io.NopCloser(strings.NewReader(reqparams.Encode()))
				hrc.Req.Body = io.NopCloser(strings.NewReader(reqparams.Encode()))
				// 将数组转换为 JSON
				jsonData, _ := json.Marshal(params_map)
				hrc.Loginfo.ReqBody = string(jsonData)
			} else {
				reqData, err := io.ReadAll(hrc.Ctx.Request.Body)
				if err != nil {
					massage = fmt.Sprintf("读取请求报文失败: %s ", err.Error())
					hrc.ErrorMassage = massage
					hrc.StatusCode = http.StatusInternalServerError
					return hrc, errors.New(massage)
				}
				hrc.Loginfo.ReqBody = string(reqData)
				hrc.Ctx.Request.Body = io.NopCloser(bytes.NewReader(reqData))
				hrc.Req.Body = io.NopCloser(bytes.NewReader(reqData))
			}
		default:
			// 当 expression 的值不等于任何一个 case 时执行的代码
			massage = fmt.Sprintf("请求方式未定义Method[%s]", method)
			hrc.ErrorMassage = massage
			hrc.StatusCode = http.StatusInternalServerError
			return hrc, errors.New(massage)
		}
	} else {
		fmt.Println("无请求数据》》》》》")
	}
	return hrc, nil
}

// 请求客户端转发
func (hrc HttpRecover) HttpClient() (HttpRecover, error) {
	// 创建 http 请求实例
	path := hrc.Ctx.Request.RequestURI
	// lb := hrc.HttpCliServieProxyRouter.Lb
	// proxy_url := lb.Next() + path
	proxy_url := fmt.Sprintf("%s%s", hrc.ConsulServiceAddress, path)
	fmt.Println("代理转发地址：", proxy_url)
	hrc.Loginfo.ProxyUrl = proxy_url
	httpReq, err := http.NewRequest(hrc.Req.Method, proxy_url, hrc.Ctx.Request.Body)
	if err != nil {
		massage := fmt.Sprintf("代理转发地址错误[" + proxy_url + "]")
		hrc.ErrorMassage = massage
		hrc.StatusCode = http.StatusInternalServerError

		return hrc, errors.New(massage)
	}
	httpReq.Header = hrc.Ctx.Request.Header
	// 发送 HTTP 请求
	resp, err := hrc.HttpCliServieProxyRouter.Client.Do(httpReq)
	if err != nil {
		massage := fmt.Sprintf("Error sending HTTP request: %v", err)
		hrc.ErrorMassage = massage
		if resp != nil {
			hrc.StatusCode = resp.StatusCode
		} else {
			hrc.StatusCode = http.StatusInternalServerError
		}

		return hrc, errors.New(massage)
	}
	hrc.Resp = resp
	return hrc, nil
}

func (hrc HttpRecover) RespBody() (HttpRecover, error) {
	resp := hrc.Resp
	defer resp.Body.Close()
	respContentType := resp.Header.Get("Content-Type")
	ContentLength := resp.ContentLength
	// fmt.Println("响应报文长度：", ContentLength)
	// fmt.Println("响应状态：", resp.StatusCode)
	// fmt.Println("响应状态：", respContentType)
	hrc.Loginfo.RespBodySize = ContentLength
	hrc.Loginfo.RespContentType = respContentType
	hrc.Loginfo.Status = int64(resp.StatusCode)
	hrc.Ctx.Status(hrc.Resp.StatusCode)
	if resp.StatusCode == http.StatusOK {
		respBody, err := io.ReadAll(resp.Body)
		if err != nil {
			// 处理错误
			massage := fmt.Sprintf("读取数据失败: %v", err)
			hrc.ErrorMassage = massage
			hrc.StatusCode = http.StatusInternalServerError

			return hrc, errors.New(massage)
		}

		//设置返回数据
		clonedBody := bytes.NewReader(respBody)

		respHeaders := make(map[string]string)
		for key, values := range resp.Header {
			respHeaders[key] = strings.Join(values, ",")
			for _, value := range values {
				hrc.Ctx.Writer.Header().Add(key, value)
			}
		}
		isFile := utils.IsFIle(respContentType)
		if isFile {
			filedata := map[string]interface{}{
				"size":        ContentLength,
				"ContentType": respContentType,
			}
			// 将数组转换为 JSON
			jsonData, _ := json.Marshal(filedata)
			hrc.Loginfo.RespBody = string(jsonData)
		} else {
			hrc.Loginfo.RespBody = string(respBody)
		}
		hrc.Loginfo.RespBodySize = int64(clonedBody.Len())
		// 将响应内容复制到c.Writer中
		_, err = io.Copy(hrc.Ctx.Writer, clonedBody)
		if err != nil {
			massage := fmt.Sprintf("Error copying response body: %v", err)
			fmt.Println(massage)
			hrc.ErrorMassage = massage
			hrc.StatusCode = http.StatusInternalServerError

			return hrc, errors.New(massage)
		}
		return hrc, nil
	} else {
		massage := "响应异常"
		hrc.ErrorMassage = massage
		hrc.StatusCode = resp.StatusCode

		return hrc, errors.New(massage)
	}
}

// 响应头处理
func (hrc HttpRecover) RespHeader() (HttpRecover, error) {
	resp := hrc.Resp
	respHeaders := make(map[string]string)
	// 遍历请求头
	for key, value := range resp.Header {
		respHeaders[key] = strings.Join(value, ",")
	}
	hrc.Loginfo.RespHeader = respHeaders
	return hrc, nil
}

func FormData(c *gin.Context) (*bytes.Buffer, string, map[string]interface{}) {
	form_data := make(map[string]interface{})
	// 创建一个缓冲区来构建请求体数据
	body := new(bytes.Buffer)
	var ContentType string
	if c.Request.ContentLength > 0 {

		writer := multipart.NewWriter(body)

		form, _ := c.MultipartForm() // 获取multipart/form-data请求体
		// files := form.File["file"]   // 获取上传的文件切片
		files := make(map[string]interface{})
		for field, fileHeaders := range form.File {
			filesArray := make([]map[string]interface{}, 0)
			for i, file := range fileHeaders {
				fileName := file.Filename
				size := file.Size
				fmt.Println(i, "文件名：", fileName, ";大小：", size)
				filedata := map[string]interface{}{
					"fileName":    fileName,
					"size":        size,
					"ContentType": file.Header.Get("Content-Type"),
				}
				filesArray = append(filesArray, filedata)

				// 创建一个文件表单字段
				fileField, _ := writer.CreateFormFile(field, fileName)
				// 将文件内容拷贝到表单字段中
				file2, _ := file.Open()
				defer file2.Close()
				_, _ = io.Copy(fileField, file2)
			}
			files[field] = filesArray
		}

		form_data["files"] = files
		for k, v := range c.Request.PostForm {
			value := strings.Join(v, ",")
			form_data[k] = strings.Join(v, ",")
			// 创建一个文本表单字段
			textField, _ := writer.CreateFormField(k)
			// 设置文本表单字段的值
			textField.Write([]byte(value))
			writer.WriteField(k, value)
		}
		// 关闭multipart.Writer，以便写入结束标记
		writer.Close()
		ContentType = writer.FormDataContentType()
		form_data["ContentType"] = ContentType
		form_data["size"] = body.Len()
	}
	return body, ContentType, form_data
}
