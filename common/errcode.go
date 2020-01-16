package common

var messages = make(map[int]string)

var ErrCodeMap = map[int]string{
	200: "操作成功",
	400: "数据格式不正确",
	401: "未授权，需登录",
	403: "没有访问权限",
	500: "服务端出错",
}

func ErrCode(code int) *ResponseData {
	return &ResponseData{
		Code: code,
		Msg:  ErrCodeMap[code],
	}
}

func GetMessage(code int) string {
	msg := ErrCodeMap[code]
	return msg
}
