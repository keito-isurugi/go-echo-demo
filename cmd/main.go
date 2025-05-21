package main

import (
	"go-echo-demo/internal/handler/api"
	"go-echo-demo/internal/handler/frontend"
	"go-echo-demo/internal/infrastructure"
	"log"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	// DB初期化
	db := infrastructure.NewDB()
	defer db.Close()

	// ユースケース初期化
	userUsecase := infrastructure.NewUserUsecase(db)

	// Echoインスタンス
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	api.RegisterRoutes(e, userUsecase)
	api.RegisterHealthRoutes(e)
	frontend.RegisterTopRoutes(e)
	frontend.RegisterBasicAuthRoutes(e)
	frontend.RegisterDigestAuthRoutes(e)
	frontend.RegisterFrontend(e)

	log.Fatal(e.Start(":8080"))
}
