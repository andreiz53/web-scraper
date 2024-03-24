package api

import (
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

		token := ParseToken(cookie.Value)

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

func (s Server) IsLoggedIn(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		userSettings := types.NewUserContext()
		cookie, err := ctx.Cookie("x-jwt-token")
		if err != nil {
			ctx.Request().Context()
			ctx.Set(types.UserContextKey, userSettings)
			return next(ctx)
		}

		token := ParseToken(cookie.Value)

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
