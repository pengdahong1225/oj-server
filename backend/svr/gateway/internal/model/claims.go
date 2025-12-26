package model

import "github.com/dgrijalva/jwt-go"

type UserClaims struct {
	jwt.StandardClaims

	Uid       int64
	Mobile    string
	Authority int32
	Type      string // access or refresh
}
