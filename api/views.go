package api

import (
	"net/http"

	"github.com/andreiz53/web-scraper/utils"
	pages "github.com/andreiz53/web-scraper/views"
	layout "github.com/andreiz53/web-scraper/views/layout"
	"github.com/labstack/echo/v4"
)

func RenderIndex(c echo.Context) error {
	return utils.Render(c, http.StatusOK, pages.Home())
}

func RenderLogin(c echo.Context) error {
	return utils.Render(c, http.StatusOK, pages.Login())
}

func RenderRegister(c echo.Context) error {
	return utils.Render(c, http.StatusOK, pages.Register())
}

func RenderHistory(c echo.Context) error {
	return utils.Render(c, http.StatusOK, pages.History())
}

func RenderLayout(c echo.Context) error {
	return utils.Render(c, http.StatusOK, layout.Page("Layout"))
}
