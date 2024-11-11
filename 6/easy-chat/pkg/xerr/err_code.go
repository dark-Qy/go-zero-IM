package xerr

// 前三位为业务码，后三位为功能码
const (
	SERVER_COMMON_ERROR = 100001
	REQUEST_PARAM_ERROR = 100002
	TOKEN_EXPIRE_ERROR  = 100003
	DB_ERROR            = 100004
)
