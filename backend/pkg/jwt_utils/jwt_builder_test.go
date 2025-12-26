package jwt_utils

import (
	"testing"
	"github.com/dgrijalva/jwt-go"
	"time"
)

type userClaims struct {
	jwt.StandardClaims

	Uid       int64
	Mobile    string
	Authority int32
	Type      string // access or refresh
}

func TestJWTBuilder_CreateToken(t *testing.T) {
	var (
		singningKey = "dAjOs8FF6H^df$BXaY@xG!%#AWMgQRbq"
	)
	jwt_builder := JWTBuilder{
		SigningKey: []byte(singningKey),
	}
	user_claims := &userClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Unix() + 60,
		},
		Uid:       1,
		Mobile:    "10000",
		Authority: 1,
		Type:      "access",
	}
	token, err := jwt_builder.CreateToken(user_claims)
	if err != nil {
		t.Fatalf("create token error: %v", err)
	}
	t.Logf("token: %s", token)
}

func TestJWTBuilder_ParseToken(t *testing.T) {
	var (
		token       = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NjY3MTk1ODIsIlVpZCI6MSwiTW9iaWxlIjoiMTAwMDAiLCJBdXRob3JpdHkiOjEsIlR5cGUiOiJhY2Nlc3MifQ.ufd7NZDgMMKgoQYS6WjwVrIk4zubRyyV32BRvAgmA00"
		singningKey = "dAjOs8FF6H^df$BXaY@xG!%#AWMgQRbq"
	)

	jwt_builder := JWTBuilder{
		SigningKey: []byte(singningKey),
	}
	user_claims := &userClaims{}
	err := jwt_builder.ParseToken(token, user_claims)
	if err != nil {
		t.Fatalf("parse token error: %v", err)
	}

	t.Logf("user_claims: %+v", user_claims)
}
