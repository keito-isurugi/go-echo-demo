package domain

import (
	"github.com/casbin/casbin/v2"
)

// CasbinRBACRepository Casbinを使用したRBACリポジトリインターフェース
type CasbinRBACRepository interface {
	// ポリシー管理
	AddPolicy(sub, obj, act string) error
	RemovePolicy(sub, obj, act string) error
	GetPolicies() ([][]string, error)

	// ロール管理
	AddRoleForUser(user, role string) error
	RemoveRoleForUser(user, role string) error
	GetRolesForUser(user string) ([]string, error)
	GetUsersForRole(role string) ([]string, error)

	// 権限チェック
	Enforce(sub, obj, act string) (bool, error)
	HasRoleForUser(user, role string) (bool, error)
}

// CasbinRBACUsecase Casbinを使用したRBACユースケースインターフェース
type CasbinRBACUsecase interface {
	// ポリシー管理
	AddPolicy(role, resource, action string) error
	RemovePolicy(role, resource, action string) error
	GetPolicies() ([][]string, error)

	// ロール管理
	AssignRoleToUser(user, role string) error
	RemoveRoleFromUser(user, role string) error
	GetUserRoles(user string) ([]string, error)
	GetRoleUsers(role string) ([]string, error)

	// 権限チェック
	CheckPermission(user, resource, action string) error
	HasRole(user, role string) (bool, error)

	// 管理機能（既存のDBベースRBACとの互換性のため）
	GetRoles() ([]Role, error)
	GetPermissions() ([]Permission, error)
	CreateRole(name, description string) (*Role, error)
	CreatePermission(name, description, resource, action string) (*Permission, error)
}

// CasbinEnforcer Casbinエンフォーサーのラッパー
type CasbinEnforcer struct {
	enforcer *casbin.Enforcer
}

// NewCasbinEnforcer Casbinエンフォーサーを作成
func NewCasbinEnforcer(enforcer *casbin.Enforcer) *CasbinEnforcer {
	return &CasbinEnforcer{enforcer: enforcer}
}

// AddPolicy ポリシーを追加
func (c *CasbinEnforcer) AddPolicy(sub, obj, act string) error {
	_, err := c.enforcer.AddPolicy(sub, obj, act)
	return err
}

// RemovePolicy ポリシーを削除
func (c *CasbinEnforcer) RemovePolicy(sub, obj, act string) error {
	_, err := c.enforcer.RemovePolicy(sub, obj, act)
	return err
}

// GetPolicies すべてのポリシーを取得
func (c *CasbinEnforcer) GetPolicies() ([][]string, error) {
	return c.enforcer.GetPolicy()
}

// AddRoleForUser ユーザーにロールを割り当て
func (c *CasbinEnforcer) AddRoleForUser(user, role string) error {
	_, err := c.enforcer.AddRoleForUser(user, role)
	return err
}

// RemoveRoleForUser ユーザーからロールを削除
func (c *CasbinEnforcer) RemoveRoleForUser(user, role string) error {
	_, err := c.enforcer.RemoveFilteredGroupingPolicy(0, user, role)
	return err
}

// GetRolesForUser ユーザーのロールを取得
func (c *CasbinEnforcer) GetRolesForUser(user string) ([]string, error) {
	return c.enforcer.GetRolesForUser(user)
}

// GetUsersForRole ロールのユーザーを取得
func (c *CasbinEnforcer) GetUsersForRole(role string) ([]string, error) {
	return c.enforcer.GetUsersForRole(role)
}

// Enforce 権限チェック
func (c *CasbinEnforcer) Enforce(sub, obj, act string) (bool, error) {
	return c.enforcer.Enforce(sub, obj, act)
}

// HasRoleForUser ユーザーが特定のロールを持っているかチェック
func (c *CasbinEnforcer) HasRoleForUser(user, role string) (bool, error) {
	return c.enforcer.HasRoleForUser(user, role)
}
