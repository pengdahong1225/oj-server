package middlewares

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/pengdahong1225/oj-server/backend/consts"
	"github.com/pengdahong1225/oj-server/backend/module/settings"
	"time"
)

type UserClaims struct {
	Uid       int64
	Mobile    int64
	Authority int32
	jwt.StandardClaims
}

var (
	TokenExpired     = errors.New("token is expired")
	TokenNotValidYet = errors.New("token not active yet")
	TokenMalformed   = errors.New("that's not even a token")
	TokenInvalid     = errors.New("couldn't handle this token")
)

type JWT struct {
	SigningKey []byte
}

func NewJWT() *JWT {
	return &JWT{
		SigningKey: []byte(settings.Instance().JwtConfig.SigningKey),
	}
}

func (receiver *JWT) CreateToken(claims *UserClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(receiver.SigningKey)
}

// ParseToken 解析签名令牌
func (receiver *JWT) ParseToken(tokenString string) (*UserClaims, error) {
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

func (receiver *JWT) RefreshToken(tokenString string) (string, error) {
	token, err := jwt.ParseWithClaims(tokenString, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		return receiver.SigningKey, nil
	})
	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(*UserClaims); ok && token.Valid {
		jwt.TimeFunc = time.Now
		claims.StandardClaims.ExpiresAt = time.Now().Unix() + consts.TokenTimeOut
		return receiver.CreateToken(claims)
	}
	return "", err
}
