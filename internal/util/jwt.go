package util

import (
	"errors"
	"fmt"
	"github.com/SyaibanAhmadRamadhan/go-collection"
	"github.com/golang-jwt/jwt/v5"
	"order-service/generated/proto/jwt_claims_proto"
)

var ErrInvalidParseToken = errors.New("invalid parse token")

type AuthAccessTokenClaims struct {
	*jwt_claims_proto.JwtAuthAccessTokenClaims
	jwt.RegisteredClaims
}

func (a *AuthAccessTokenClaims) ClaimsHS256(tokenStr string, key string) (err error) {
	tokenParse, err := jwt.ParseWithClaims(tokenStr, a, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(key), nil
	})
	if err != nil {
		return collection.Err(err)
	}

	if !tokenParse.Valid {
		return collection.Err(ErrInvalidParseToken)
	}

	return
}
