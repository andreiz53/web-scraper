package api

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

type contextKey string

type ContextWithValue struct {
	echo.Context
}

func NewMiddlewareContextValue(fn echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		return fn(ContextWithValue{ctx})
	}
}

func (ctx ContextWithValue) Get(key string) interface{} {
	val := ctx.Context.Get(key)
	if val != nil {
		return val
	}
	return ctx.Request().Context().Value(key)
}
func (ctx ContextWithValue) Set(key string, val interface{}) {
	var newKey contextKey = contextKey(key)
	ctx.SetRequest(ctx.Request().WithContext(context.WithValue(ctx.Request().Context(), newKey, val)))
}

func ParseToken(tokenStr string) *jwt.Token {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if err != nil {
		log.Fatal(err)
	}

	return token
}

func GenerateToken(claims jwt.Claims) (tokenStr string, err error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err = token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return "", err
	}
	return tokenStr, nil
}
