package api

import (
	"net/http"
	"time"

	"github.com/andreiz53/web-scraper/types"
	pages "github.com/andreiz53/web-scraper/views"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

func RenderIndex(ctx echo.Context) error {
	return Render(ctx, http.StatusOK, pages.Home())
}

func RenderLogin(ctx echo.Context) error {
	return Render(ctx, http.StatusOK, pages.Login())
}

func RenderRegister(ctx echo.Context) error {
	return Render(ctx, http.StatusOK, pages.Register())
}

func (s Server) RenderHistory(ctx echo.Context) error {
	userCtx := ctx.Get(types.UserContextKey).(*types.UserContext)
	var websites []types.Website
	if userCtx.IsLoggedIn {
		cookie, err := ctx.Cookie("x-jwt-token")
		if err != nil {
			return Render(ctx, http.StatusOK, pages.History(websites))
		}

		token := ParseToken(cookie.Value)

		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			if time.Now().Unix() > int64(claims["exp"].(float64)) {
				return Render(ctx, http.StatusOK, pages.History(websites))
			}
			user, err := s.Store.GetUserByEmail(claims["sub"].(string))
			if err != nil {
				return Render(ctx, http.StatusOK, pages.History(websites))
			}

			if user.ID == 0 {
				return Render(ctx, http.StatusOK, pages.History(websites))
			}

			websites, err = s.Store.GetWebsitesByUserID(user.ID)
			if err != nil {
				return Render(ctx, http.StatusOK, pages.History(websites))
			}

			return Render(ctx, http.StatusOK, pages.History(websites))

		}
	}
	return Render(ctx, http.StatusOK, pages.History(websites))
}
