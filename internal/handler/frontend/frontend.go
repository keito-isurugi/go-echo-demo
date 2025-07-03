package frontend

import (
	"html/template"
	"io"

	"go-echo-demo/internal/domain"
	"go-echo-demo/internal/middleware"

	"github.com/labstack/echo/v4"
)

type TemplateRenderer struct {
	Templates *template.Template
}

func (t *TemplateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.Templates.ExecuteTemplate(w, name, data)
}

func RegisterFrontend(e *echo.Echo) {
	e.Renderer = &TemplateRenderer{
		Templates: template.Must(template.ParseFiles(
			"templates/top.html",
			"templates/basic.html",
			"templates/digest.html",
			"templates/login.html",
			"templates/protected.html",
			"templates/google_login.html",
			"templates/line_login.html",
			"templates/_header.html",
			"templates/_footer.html",
		)),
	}
	e.Static("/static", "static")
	e.File("/rbac", "templates/rbac_admin.html")
	e.File("/casbin", "templates/casbin_admin.html")
}

func RegisterAuthFrontendRoutes(e *echo.Echo, authUsecase domain.AuthUsecase) {
	// 認証不要のルート
	e.GET("/login", LoginPage)
	e.GET("/line-login", LineLoginPage)

	// 保護されたページ（JWT認証付き）
	protected := e.Group("/protected")
	protected.Use(middleware.JWTAuth(authUsecase))
	protected.GET("", ProtectedPage)

	// 認証が必要なAPIエンドポイント（ユーザー情報取得用）
	apiProtected := e.Group("/api/user")
	apiProtected.Use(middleware.JWTAuth(authUsecase))
	apiProtected.GET("/info", GetUserInfo)
}
