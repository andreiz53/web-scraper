package api

import (
	"net/http"
	"strings"
	"time"

	"github.com/andreiz53/web-scraper/scraper"
	"github.com/andreiz53/web-scraper/types"
	partialsWebsite "github.com/andreiz53/web-scraper/views/partials/website"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

func (s Server) CheckWebsite(ctx echo.Context) error {
	website := ctx.FormValue("website")
	if website == "" {
		return Render(ctx, http.StatusOK, partialsWebsite.WebsiteFail())
	}
	scr, err := scraper.NewScraper(website)
	if err != nil {
		return Render(ctx, http.StatusOK, partialsWebsite.WebsiteFail())
	}
	scr.SetSEO()
	scr.Collector.Visit(scr.URL)

	// get the logged in user ID
	userCtx := ctx.Get(types.UserContextKey).(*types.UserContext)
	if userCtx.IsLoggedIn {
		cookie, err := ctx.Cookie("x-jwt-token")
		if err != nil {
			return Render(ctx, http.StatusOK, partialsWebsite.WebsiteFail())
		}

		token := ParseToken(cookie.Value)

		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			if time.Now().Unix() > int64(claims["exp"].(float64)) {
				return Render(ctx, http.StatusOK, partialsWebsite.WebsiteFail())
			}
			user, err := s.Store.GetUserByEmail(claims["sub"].(string))
			if err != nil {
				return Render(ctx, http.StatusOK, partialsWebsite.WebsiteFail())
			}

			if user.ID == 0 {
				return Render(ctx, http.StatusOK, partialsWebsite.WebsiteFail())
			}

			s.Store.CreateWebsite(&types.Website{
				URL:            scr.URL,
				UserID:         user.ID,
				Robots:         scr.Data.Robots,
				Keywords:       strings.Join(scr.Data.Keywords, ", "),
				Description:    scr.Data.Description,
				HasSitemap:     scr.Data.HasSitemap,
				HasOpenGraph:   scr.Data.HasOpenGraph,
				HasTwitterCard: scr.Data.HasTwitterCard,
				HasDublinCore:  scr.Data.HasDublinCore,
				HasJsonLd:      scr.Data.HasJsonLd,
				HasMetaPixel:   scr.Data.HasMetaPixel,
			})
		}
	}
	return Render(ctx, http.StatusOK, partialsWebsite.Website(scr.Data))
}
