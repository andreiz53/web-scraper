package api

import (
	"fmt"

	"github.com/andreiz53/web-scraper/db"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Server struct {
	listenAddr string
	Store      *db.SqliteStorage
}

func NewServer(addr string, store *db.SqliteStorage) *Server {
	return &Server{
		listenAddr: addr,
		Store:      store,
	}
}

func (s Server) Start() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Static("/static/css", "./static/css")

	e.Use(middleware.CORS())
	e.Use(WithContextValue)
	e.Use(s.IsLoggedIn)

	e.GET("/", RenderIndex)
	e.GET("/login", RenderLogin)
	e.GET("/register", RenderRegister)

	e.POST("/register", s.RegisterUser)
	e.POST("/login", s.LoginUser)

	e.POST("/website", s.CheckWebsite)

	e.GET("/history", s.RenderHistory, s.RequireAuth)

	fmt.Println("Server listening at:", s.listenAddr)
	e.Logger.Fatal(e.Start(s.listenAddr))
}
