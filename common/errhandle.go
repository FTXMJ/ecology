package common

type ResponseData struct {
	Code int         `json:"code""`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data,omitempty"`
}

// NewResponse 成功返回
func NewResponse(data interface{}) *ResponseData {
	return &ResponseData{
		Code: 200,
		Msg:  ErrCodeMap[200],
		Data: data,
	}
}

// NewErrorResponse 错误返回
func NewErrorResponse(code int, msg string, data interface{}) *ResponseData {
	return &ResponseData{
		Code: code,
		Msg:  msg,
		Data: data,
	}
}
