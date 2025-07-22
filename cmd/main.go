package main

import (
	"log"

	"go-echo-demo/internal/handler/api"
	"go-echo-demo/internal/handler/frontend"
	"go-echo-demo/internal/infrastructure"
	"go-echo-demo/internal/repository"
	"go-echo-demo/internal/usecase"

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
	
	// sqlx DB for product repository
	dbx := infrastructure.NewDBX()
	defer dbx.Close()

	// リポジトリ初期化
	authRepo := infrastructure.NewAuthRepository(db)
	refreshTokenRepo := infrastructure.NewRefreshTokenRepository(db)
	userRepo := infrastructure.NewUserRepository(db)
	oauthRepo := infrastructure.NewOAuthRepository(db)
	rbacRepo := repository.NewRBACRepository(db)
	productRepo := repository.NewProductRepository(dbx)

	// Casbin RBAC初期化
	casbinRepo, err := infrastructure.NewCasbinRBACRepository()
	if err != nil {
		log.Printf("Warning: Casbin初期化に失敗しました: %v", err)
	}

	// ユースケース初期化
	userUsecase := infrastructure.NewUserUsecase(db)
	authUsecase := infrastructure.NewAuthUsecase(authRepo, refreshTokenRepo, userRepo)
	rbacUsecase := usecase.NewRBACUsecase(rbacRepo)
	casbinUsecase := infrastructure.NewCasbinRBACUsecase(casbinRepo, rbacRepo)
	productUsecase := usecase.NewProductUsecase(productRepo)
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
	api.RegisterRBACRoutes(e, rbacUsecase)
	api.RegisterCasbinRBACRoutes(e, casbinUsecase)

	frontend.RegisterTopRoutes(e)
	frontend.RegisterBasicAuthRoutes(e)
	frontend.RegisterDigestAuthRoutes(e)
	frontend.RegisterFrontend(e)
	frontend.RegisterAuthFrontendRoutes(e, authUsecase)

	// SQLインジェクションデモルート
	frontend.RegisterSqlInjectionRoutes(e, productUsecase)
	api.RegisterSqlInjectionAPIRoutes(e, productUsecase)

	log.Fatal(e.Start(":8080"))
}
