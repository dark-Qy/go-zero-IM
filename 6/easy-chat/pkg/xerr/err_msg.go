package xerr

var codeText = map[int]string{
	SERVER_COMMON_ERROR: "服务器内部错误",
	REQUEST_PARAM_ERROR: "请求参数错误",
	TOKEN_EXPIRE_ERROR:  "token过期",
	DB_ERROR:            "数据库错误",
}

func ErrMsg(code int) string {
	if msg, ok := codeText[code]; ok {
		return msg
	}
	return codeText[SERVER_COMMON_ERROR]
}
