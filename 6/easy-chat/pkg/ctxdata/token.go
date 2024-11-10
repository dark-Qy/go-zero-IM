package ctxdata

import "github.com/golang-jwt/jwt/v4"

const Identify = "imooc.com"

func GetJwtToken(secretKey string, iat, seconds int64, uid string) (string, error) {
	// 定义声明中的相关字段
	claims := make(jwt.MapClaims)
	claims["exp"] = iat + seconds
	claims["iat"] = iat
	claims[Identify] = uid

	// 生成token
	token := jwt.New(jwt.SigningMethodHS256)
	token.Claims = claims

	// 根据密钥进行加密
	return token.SignedString([]byte(secretKey))
}
