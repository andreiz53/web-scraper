package api

import (
	"fmt"
	"net/http"

	"github.com/andreiz53/web-scraper/scraper"
	"github.com/andreiz53/web-scraper/utils"
	partialsWebsite "github.com/andreiz53/web-scraper/views/partials/website"
	"github.com/labstack/echo/v4"
)

func (s Server) CheckWebsite(ctx echo.Context) error {
	website := ctx.FormValue("website")
	if website == "" {
		return echo.ErrBadRequest
	}
	scr, err := scraper.NewScraper(website)
	if err != nil {
		return echo.ErrBadRequest
	}
	scr.SetSEO()
	scr.Collector.Visit(scr.URL)
	fmt.Printf("SENDING DATA: %+v\n", scr.Data.Robots)

	return utils.Render(ctx, http.StatusOK, partialsWebsite.Website(scr.Data))
}
