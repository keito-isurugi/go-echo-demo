package domain

import "time"

// Role ロール情報
type Role struct {
	ID          int       `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Description string    `json:"description" db:"description"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// Permission 権限情報
type Permission struct {
	ID          int       `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Description string    `json:"description" db:"description"`
	Resource    string    `json:"resource" db:"resource"`
	Action      string    `json:"action" db:"action"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// UserRole ユーザーロール関連
type UserRole struct {
	ID     int `json:"id" db:"id"`
	UserID int `json:"user_id" db:"user_id"`
	RoleID int `json:"role_id" db:"role_id"`
}

// RolePermission ロール権限関連
type RolePermission struct {
	ID           int `json:"id" db:"id"`
	RoleID       int `json:"role_id" db:"role_id"`
	PermissionID int `json:"permission_id" db:"permission_id"`
}

// UserWithRoles ロール情報を含むユーザー
type UserWithRoles struct {
	User  User   `json:"user"`
	Roles []Role `json:"roles"`
}

// RoleWithPermissions 権限情報を含むロール
type RoleWithPermissions struct {
	Role        Role         `json:"role"`
	Permissions []Permission `json:"permissions"`
}

// RBACRepository RBACリポジトリインターフェース
type RBACRepository interface {
	// ロール関連
	GetRoles() ([]Role, error)
	GetRoleByID(id int) (*Role, error)
	GetRoleByName(name string) (*Role, error)
	CreateRole(role *Role) error
	UpdateRole(role *Role) error
	DeleteRole(id int) error

	// 権限関連
	GetPermissions() ([]Permission, error)
	GetPermissionByID(id int) (*Permission, error)
	GetPermissionByName(name string) (*Permission, error)
	CreatePermission(permission *Permission) error
	UpdatePermission(permission *Permission) error
	DeletePermission(id int) error

	// ユーザーロール関連
	GetUserRoles(userID int) ([]Role, error)
	AssignRoleToUser(userID, roleID int) error
	RemoveRoleFromUser(userID, roleID int) error
	GetUsersByRole(roleID int) ([]User, error)

	// ロール権限関連
	GetRolePermissions(roleID int) ([]Permission, error)
	AssignPermissionToRole(roleID, permissionID int) error
	RemovePermissionFromRole(roleID, permissionID int) error

	// 権限チェック
	HasPermission(userID int, resource, action string) (bool, error)
	HasRole(userID int, roleName string) (bool, error)
}

// RBACUsecase RBACユースケースインターフェース
type RBACUsecase interface {
	// ロール管理
	GetRoles() ([]Role, error)
	GetRoleByID(id int) (*Role, error)
	CreateRole(name, description string) (*Role, error)
	UpdateRole(id int, name, description string) (*Role, error)
	DeleteRole(id int) error

	// 権限管理
	GetPermissions() ([]Permission, error)
	GetPermissionByID(id int) (*Permission, error)
	CreatePermission(name, description, resource, action string) (*Permission, error)
	UpdatePermission(id int, name, description, resource, action string) (*Permission, error)
	DeletePermission(id int) error

	// ユーザーロール管理
	GetUserRoles(userID int) ([]Role, error)
	AssignRoleToUser(userID int, roleName string) error
	RemoveRoleFromUser(userID int, roleName string) error

	// ロール権限管理
	GetRolePermissions(roleID int) ([]Permission, error)
	AssignPermissionToRole(roleName, permissionName string) error
	RemovePermissionFromRole(roleName, permissionName string) error

	// 権限チェック
	HasPermission(userID int, resource, action string) (bool, error)
	HasRole(userID int, roleName string) (bool, error)
	CheckPermission(userID int, resource, action string) error
}
