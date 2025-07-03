package api

import (
	"net/http"
	"strconv"

	"go-echo-demo/internal/domain"
	"go-echo-demo/internal/middleware"

	"github.com/labstack/echo/v4"
)

type CasbinRBACHandler struct {
	casbinUsecase domain.CasbinRBACUsecase
}

func NewCasbinRBACHandler(casbinUsecase domain.CasbinRBACUsecase) *CasbinRBACHandler {
	return &CasbinRBACHandler{casbinUsecase: casbinUsecase}
}

// RegisterCasbinRBACRoutes Casbin RBACルートを登録
func RegisterCasbinRBACRoutes(e *echo.Echo, casbinUsecase domain.CasbinRBACUsecase) {
	h := NewCasbinRBACHandler(casbinUsecase)

	// 管理者権限が必要なルートグループ
	adminGroup := e.Group("/admin/casbin")
	// 注意: 実際の使用時は authUsecase も必要になります
	// adminGroup.Use(middleware.CasbinJWTAuthWithRole(authUsecase, casbinUsecase, "admin"))

	// ポリシー管理API
	adminGroup.GET("/policies", h.GetPolicies)
	adminGroup.POST("/policies", h.AddPolicy)
	adminGroup.DELETE("/policies", h.RemovePolicy)

	// ロール管理API
	adminGroup.GET("/users/:user/roles", h.GetUserRoles)
	adminGroup.POST("/users/:user/roles", h.AssignRoleToUser)
	adminGroup.DELETE("/users/:user/roles/:role", h.RemoveRoleFromUser)
	adminGroup.GET("/roles/:role/users", h.GetRoleUsers)

	// 管理機能（既存のDBベースRBACとの互換性）
	adminGroup.GET("/roles", h.GetRoles)
	adminGroup.GET("/permissions", h.GetPermissions)
	adminGroup.POST("/roles", h.CreateRole)
	adminGroup.POST("/permissions", h.CreatePermission)

	// 一般ユーザー用API
	userGroup := e.Group("/api/casbin")
	// 注意: 実際の使用時は authUsecase も必要になります
	// userGroup.Use(middleware.CasbinJWTAuthWithRole(authUsecase, casbinUsecase, "user", "admin"))
	userGroup.GET("/my/roles", h.GetMyRoles)
}

// ポリシー管理API

// GetPolicies ポリシー一覧取得
func (h *CasbinRBACHandler) GetPolicies(c echo.Context) error {
	policies, err := h.casbinUsecase.GetPolicies()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "ポリシーの取得に失敗しました")
	}
	return c.JSON(http.StatusOK, policies)
}

// AddPolicy ポリシー追加
func (h *CasbinRBACHandler) AddPolicy(c echo.Context) error {
	var req struct {
		Role     string `json:"role" validate:"required"`
		Resource string `json:"resource" validate:"required"`
		Action   string `json:"action" validate:"required"`
	}

	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "リクエストの解析に失敗しました")
	}

	if err := c.Validate(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "バリデーションエラー")
	}

	err := h.casbinUsecase.AddPolicy(req.Role, req.Resource, req.Action)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "ポリシーの追加に失敗しました")
	}

	return c.JSON(http.StatusCreated, map[string]string{"message": "ポリシーが追加されました"})
}

// RemovePolicy ポリシー削除
func (h *CasbinRBACHandler) RemovePolicy(c echo.Context) error {
	var req struct {
		Role     string `json:"role" validate:"required"`
		Resource string `json:"resource" validate:"required"`
		Action   string `json:"action" validate:"required"`
	}

	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "リクエストの解析に失敗しました")
	}

	if err := c.Validate(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "バリデーションエラー")
	}

	err := h.casbinUsecase.RemovePolicy(req.Role, req.Resource, req.Action)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "ポリシーの削除に失敗しました")
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "ポリシーが削除されました"})
}

// ロール管理API

// GetUserRoles ユーザーのロール一覧取得
func (h *CasbinRBACHandler) GetUserRoles(c echo.Context) error {
	user := c.Param("user")
	if user == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "ユーザー名が必要です")
	}

	roles, err := h.casbinUsecase.GetUserRoles(user)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "ユーザーロールの取得に失敗しました")
	}

	return c.JSON(http.StatusOK, roles)
}

// AssignRoleToUser ユーザーにロールを割り当て
func (h *CasbinRBACHandler) AssignRoleToUser(c echo.Context) error {
	user := c.Param("user")
	if user == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "ユーザー名が必要です")
	}

	var req struct {
		Role string `json:"role" validate:"required"`
	}

	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "リクエストの解析に失敗しました")
	}

	if err := c.Validate(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "バリデーションエラー")
	}

	err := h.casbinUsecase.AssignRoleToUser(user, req.Role)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "ロールの割り当てに失敗しました")
	}

	return c.JSON(http.StatusCreated, map[string]string{"message": "ロールが割り当てられました"})
}

// RemoveRoleFromUser ユーザーからロールを削除
func (h *CasbinRBACHandler) RemoveRoleFromUser(c echo.Context) error {
	user := c.Param("user")
	role := c.Param("role")

	if user == "" || role == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "ユーザー名とロール名が必要です")
	}

	err := h.casbinUsecase.RemoveRoleFromUser(user, role)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "ロールの削除に失敗しました")
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "ロールが削除されました"})
}

// GetRoleUsers ロールのユーザー一覧取得
func (h *CasbinRBACHandler) GetRoleUsers(c echo.Context) error {
	role := c.Param("role")
	if role == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "ロール名が必要です")
	}

	users, err := h.casbinUsecase.GetRoleUsers(role)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "ロールユーザーの取得に失敗しました")
	}

	return c.JSON(http.StatusOK, users)
}

// 管理機能（既存のDBベースRBACとの互換性）

// GetRoles ロール一覧取得
func (h *CasbinRBACHandler) GetRoles(c echo.Context) error {
	roles, err := h.casbinUsecase.GetRoles()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "ロールの取得に失敗しました")
	}
	return c.JSON(http.StatusOK, roles)
}

// GetPermissions 権限一覧取得
func (h *CasbinRBACHandler) GetPermissions(c echo.Context) error {
	permissions, err := h.casbinUsecase.GetPermissions()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "権限の取得に失敗しました")
	}
	return c.JSON(http.StatusOK, permissions)
}

// CreateRole ロール作成
func (h *CasbinRBACHandler) CreateRole(c echo.Context) error {
	var req struct {
		Name        string `json:"name" validate:"required"`
		Description string `json:"description"`
	}

	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "リクエストの解析に失敗しました")
	}

	if err := c.Validate(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "バリデーションエラー")
	}

	role, err := h.casbinUsecase.CreateRole(req.Name, req.Description)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "ロールの作成に失敗しました")
	}

	return c.JSON(http.StatusCreated, role)
}

// CreatePermission 権限作成
func (h *CasbinRBACHandler) CreatePermission(c echo.Context) error {
	var req struct {
		Name        string `json:"name" validate:"required"`
		Description string `json:"description"`
		Resource    string `json:"resource" validate:"required"`
		Action      string `json:"action" validate:"required"`
	}

	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "リクエストの解析に失敗しました")
	}

	if err := c.Validate(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "バリデーションエラー")
	}

	permission, err := h.casbinUsecase.CreatePermission(req.Name, req.Description, req.Resource, req.Action)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "権限の作成に失敗しました")
	}

	return c.JSON(http.StatusCreated, permission)
}

// GetMyRoles 自分のロール一覧取得
func (h *CasbinRBACHandler) GetMyRoles(c echo.Context) error {
	userID, err := middleware.GetUserIDFromContext(c)
	if err != nil {
		return err
	}

	// ユーザーIDを文字列に変換
	user := strconv.Itoa(userID)

	roles, err := h.casbinUsecase.GetUserRoles(user)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "ユーザーロールの取得に失敗しました")
	}

	return c.JSON(http.StatusOK, roles)
}
