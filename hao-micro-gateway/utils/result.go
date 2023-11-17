package utils

type Result struct {
	Success bool        `json:"success"`
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Count   int         `json:"count"`
}

// 工厂函数，创建一个成功的响应
func NewSuccessResult(data interface{}) Result {
	return Result{
		Success: true,
		Code:    200,
		Message: "OK",
		Data:    data,
		Count:   0,
	}
}

// 工厂函数，创建一个失败的响应
func NewErrorResult(code int, message string) Result {
	return Result{
		Success: false,
		Code:    code,
		Message: message,
		Data:    nil,
		Count:   0,
	}
}
