package jwt_utils

import (
	"github.com/dgrijalva/jwt-go"
	"errors"
)

var (
	TokenExpired     = errors.New("token is expired")
	TokenNotValidYet = errors.New("token not active yet")
	TokenMalformed   = errors.New("that's not even a token")
	TokenInvalid     = errors.New("couldn't handle this token")
)

type JWTBuilder struct {
	SigningKey []byte
}

// 创建身份令牌
func (j *JWTBuilder) CreateToken(claims jwt.Claims) (string, error) {
	// 选择HS256加密的HMAC
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// 使用服务端密钥对token签名
	return token.SignedString(j.SigningKey)
}

// 解析签名令牌
func (j *JWTBuilder) ParseToken(tokenString string, claims jwt.Claims) error {
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return j.SigningKey, nil
	})
	if err != nil {
		var ve *jwt.ValidationError
		if errors.As(err, &ve) {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return TokenMalformed
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				return TokenExpired
			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				return TokenNotValidYet
			} else {
				return TokenInvalid
			}
		}
	}
	if token == nil {
		return TokenInvalid
	}

	return nil
}
