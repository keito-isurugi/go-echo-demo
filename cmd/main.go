package main

import (
	"log"

	"go-echo-demo/internal/handler/api"
	"go-echo-demo/internal/handler/frontend"
	"go-echo-demo/internal/infrastructure"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	// 環境変数の読み込み
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found")
	}

	// DB初期化
	db := infrastructure.NewDB()
	defer db.Close()

	// ユースケース初期化
	userUsecase := infrastructure.NewUserUsecase(db)
	authRepo := infrastructure.NewAuthRepository(db)
	authUsecase := infrastructure.NewAuthUsecase(authRepo)
	oauthRepo := infrastructure.NewOAuthRepository(db)
	stateManager := infrastructure.NewStateManager()
	oauthProviders := infrastructure.NewOAuthProviders(oauthRepo, authUsecase, stateManager)

	// Echoインスタンス
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// ルート登録
	api.RegisterRoutes(e, userUsecase)
	api.RegisterHealthRoutes(e)
	api.RegisterAuthRoutes(e, authUsecase)
	api.RegisterOAuthRoutes(e, oauthProviders)

	frontend.RegisterTopRoutes(e)
	frontend.RegisterBasicAuthRoutes(e)
	frontend.RegisterDigestAuthRoutes(e)
	frontend.RegisterFrontend(e)
	frontend.RegisterAuthFrontendRoutes(e, authUsecase)

	log.Fatal(e.Start(":8080"))
}
