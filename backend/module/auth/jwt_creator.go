package auth

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
)

var (
	TokenExpired     = errors.New("token is expired")
	TokenNotValidYet = errors.New("token not active yet")
	TokenMalformed   = errors.New("that's not even a token")
	TokenInvalid     = errors.New("couldn't handle this token")
)

type UserClaims struct {
	Uid       int64
	Mobile    string
	Authority int32
	Type      string // access or refresh
	jwt.StandardClaims
}

type JWTCreator struct {
	SigningKey []byte
}

// CreateToken 创建身份令牌
func (receiver *JWTCreator) CreateToken(claims *UserClaims) (string, error) {
	// 选择HS256加密的HMAC
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// 使用服务端密钥对token签名
	return token.SignedString(receiver.SigningKey)
}

// ParseToken 解析签名令牌
func (receiver *JWTCreator) ParseToken(tokenString string) (*UserClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		return receiver.SigningKey, nil
	})
	if err != nil {
		var ve *jwt.ValidationError
		if errors.As(err, &ve) {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, TokenMalformed
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				return nil, TokenExpired
			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				return nil, TokenNotValidYet
			} else {
				return nil, TokenInvalid
			}
		}
	}
	if token == nil {
		return nil, TokenInvalid
	}
	if claims, ok := token.Claims.(*UserClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, TokenInvalid
}
