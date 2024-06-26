package api

import (
	"net/http"
	"time"

	"github.com/andreiz53/web-scraper/types"
	partialsLogin "github.com/andreiz53/web-scraper/views/partials/login"
	partialsRegister "github.com/andreiz53/web-scraper/views/partials/register"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

func (s Server) RegisterUser(ctx echo.Context) error {
	email := ctx.FormValue("email")
	password := ctx.FormValue("password")
	if email == "" || password == "" {
		return Render(ctx, http.StatusBadRequest, partialsRegister.RegisterFail())
	}
	encryptedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return Render(ctx, http.StatusBadRequest, partialsRegister.RegisterFail())
	}
	user := types.NewUser(email, string(encryptedPassword))
	err = s.Store.CreateUser(user)
	if err != nil {
		return Render(ctx, http.StatusBadRequest, partialsRegister.RegisterFail())
	}

	return Render(ctx, http.StatusOK, partialsRegister.RegisterSuccess())
}

func (s Server) LoginUser(ctx echo.Context) error {
	email := ctx.FormValue("email")
	password := ctx.FormValue("password")
	if email == "" || password == "" {
		return echo.ErrUnauthorized
	}
	user, err := s.Store.GetUserByEmail(email)
	if err != nil {
		Render(ctx, http.StatusUnauthorized, partialsLogin.LoginFail())
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return Render(ctx, http.StatusUnauthorized, partialsLogin.LoginFail())
	}

	claims := jwt.MapClaims{
		"sub": user.Email,
		"exp": time.Now().Add(time.Hour * 72).Unix(),
	}
	tokenString, err := GenerateToken(claims)
	if err != nil {
		return Render(ctx, http.StatusInternalServerError, partialsLogin.LoginFail())
	}

	cookie := new(http.Cookie)
	cookie.Name = "x-jwt-token"
	cookie.Value = tokenString
	cookie.Expires = time.Now().Add(24 * time.Hour)
	cookie.HttpOnly = true
	ctx.SetCookie(cookie)

	ctx.Response().Header().Set("HX-Redirect", "/")

	return Render(ctx, http.StatusOK, partialsLogin.LoginSuccess())
}
