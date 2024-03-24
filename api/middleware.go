package api

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/andreiz53/web-scraper/types"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

func (s Server) RequireAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		cookie, err := ctx.Cookie("x-jwt-token")
		if err != nil {
			return ctx.String(401, "Unauthorized")
		}

		token := parseToken(cookie.Value)

		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			if time.Now().Unix() > int64(claims["exp"].(float64)) {
				return ctx.String(401, "Unauthorized")
			}
			user, err := s.Store.GetUserByEmail(claims["sub"].(string))
			if err != nil {
				return ctx.String(401, "Unauthorized")
			}

			if user.ID == 0 {
				return ctx.String(401, "Unauthorized")
			}

			ctx.Set("sub", user.Email)
		}
		return next(ctx)
	}

}

type contextWithValue struct {
	echo.Context
}

func NewMiddlewareContextValue(fn echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		return fn(contextWithValue{ctx})
	}
}

func (ctx contextWithValue) Get(key string) interface{} {
	val := ctx.Context.Get(key)
	if val != nil {
		return val
	}
	return ctx.Request().Context().Value(key)
}
func (ctx contextWithValue) Set(key string, val interface{}) {
	ctx.SetRequest(ctx.Request().WithContext(context.WithValue(ctx.Request().Context(), key, val)))
}

func (s Server) IsLoggedIn(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		userSettings := types.NewUserContext()
		cookie, err := ctx.Cookie("x-jwt-token")
		if err != nil {
			ctx.Request().Context()
			ctx.Set(types.UserContextKey, userSettings)
			return next(ctx)
		}

		token := parseToken(cookie.Value)

		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			if time.Now().Unix() > int64(claims["exp"].(float64)) {
				ctx.Set(types.UserContextKey, userSettings)
				return next(ctx)
			}
			user, err := s.Store.GetUserByEmail(claims["sub"].(string))
			if err != nil {
				ctx.Set(types.UserContextKey, userSettings)
				return next(ctx)
			}

			if user.ID == 0 {
				ctx.Set(types.UserContextKey, userSettings)
				return next(ctx)
			}

		}
		userSettings.IsLoggedIn = true
		ctx.Set(types.UserContextKey, userSettings)
		return next(ctx)
	}
}

func parseToken(tokenStr string) *jwt.Token {
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
