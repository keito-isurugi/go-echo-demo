package domain

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Claims JWTクレームの構造体
type Claims struct {
	UserID int    `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

// AuthRequest 認証リクエストの構造体
type AuthRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// AuthResponse 認証レスポンスの構造体
type AuthResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}

// AuthRepository 認証リポジトリのインターフェース
type AuthRepository interface {
	ValidateCredentials(email, password string) (*User, error)
}

// AuthUsecase 認証ユースケースのインターフェース
type AuthUsecase interface {
	Login(email, password string) (*AuthResponse, error)
	ValidateToken(tokenString string) (*Claims, error)
	GenerateToken(user *User) (string, error)
}

// JWTConfig JWT設定の構造体
type JWTConfig struct {
	SecretKey string
	Duration  time.Duration
}

// StateManager stateパラメータの管理インターフェース
type StateManager interface {
	GenerateState() string
	ValidateState(state string) bool
}
