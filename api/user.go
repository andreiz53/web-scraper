package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/andreiz53/web-scraper/types"
	"github.com/andreiz53/web-scraper/utils"
	partialsLogin "github.com/andreiz53/web-scraper/views/partials/login"
	partialsRegister "github.com/andreiz53/web-scraper/views/partials/register"
	partialsWebsite "github.com/andreiz53/web-scraper/views/partials/website"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

func (s Server) RegisterUser(ctx echo.Context) error {
	email := ctx.FormValue("email")
	password := ctx.FormValue("password")
	if email == "" || password == "" {
		return utils.Render(ctx, http.StatusBadRequest, partialsRegister.RegisterFail())
	}
	encryptedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return utils.Render(ctx, http.StatusBadRequest, partialsRegister.RegisterFail())
	}
	user := types.NewUser(email, string(encryptedPassword))
	err = s.Store.CreateUser(user)
	if err != nil {
		return utils.Render(ctx, http.StatusBadRequest, partialsRegister.RegisterFail())
	}

	return utils.Render(ctx, http.StatusOK, partialsRegister.RegisterSuccess())
}

func (s Server) LoginUser(ctx echo.Context) error {
	email := ctx.FormValue("email")
	password := ctx.FormValue("password")
	if email == "" || password == "" {
		return echo.ErrUnauthorized
	}
	user, err := s.Store.GetUserByEmail(email)
	if err != nil {
		utils.Render(ctx, http.StatusUnauthorized, partialsLogin.LoginFail())
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return utils.Render(ctx, http.StatusUnauthorized, partialsLogin.LoginFail())
	}

	claims := jwt.MapClaims{
		"sub": user.Email,
		"exp": time.Now().Add(time.Hour * 72).Unix(),
	}
	tokenString, err := GenerateToken(claims)
	if err != nil {
		return utils.Render(ctx, http.StatusInternalServerError, partialsLogin.LoginFail())
	}

	// ctx.Response().Header().Set("Authorization", tokenString)
	cookie := new(http.Cookie)
	cookie.Name = "x-jwt-token"
	cookie.Value = tokenString
	cookie.Expires = time.Now().Add(24 * time.Hour)
	cookie.HttpOnly = true
	ctx.SetCookie(cookie)

	return utils.Render(ctx, http.StatusOK, partialsLogin.LoginSuccess())
}

func (s Server) CheckWebsite(ctx echo.Context) error {
	website := ctx.FormValue("website")
	if website == "" {
		return echo.ErrBadRequest
	}
	scraper, err := NewScraper(website)
	if err != nil {
		return echo.ErrBadRequest
	}
	scraper.SetSEO()
	scraper.Collector.Visit(scraper.URL)
	fmt.Printf("SENDING DATA: %+v\n", scraper.Data.Robots)

	return utils.Render(ctx, http.StatusOK, partialsWebsite.Website(scraper.Data))
}
