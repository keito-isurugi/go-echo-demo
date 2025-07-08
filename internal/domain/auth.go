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
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token,omitempty"`
	User         User   `json:"user"`
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
	// リフレッシュトークンを使用して新しいトークンペアを生成
	RefreshToken(refreshToken string) (*TokenPair, error)
	// ログアウト時にリフレッシュトークンを無効化
	Logout(userID int) error
	// トークンペアを生成
	GenerateTokenPair(user *User, deviceInfo, ipAddress string) (*TokenPair, error)
}

// JWTConfig JWT設定の構造体
type JWTConfig struct {
	SecretKey              string
	Duration               time.Duration // アクセストークンの有効期限
	RefreshTokenDuration   time.Duration // リフレッシュトークンの有効期限
}

// StateManager stateパラメータの管理インターフェース
type StateManager interface {
	GenerateState() string
	ValidateState(state string) bool
}

// RefreshToken リフレッシュトークンのエンティティ
type RefreshToken struct {
	ID             int       `json:"id"`
	UserID         int       `json:"user_id"`
	Token          string    `json:"token"`
	AccessTokenJTI string    `json:"access_token_jti"`
	ExpiresAt      time.Time `json:"expires_at"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	LastUsedAt     *time.Time `json:"last_used_at"`
	Revoked        bool      `json:"revoked"`
	RevokedAt      *time.Time `json:"revoked_at"`
	DeviceInfo     string    `json:"device_info"`
	IPAddress      string    `json:"ip_address"`
}

// RefreshTokenRequest リフレッシュトークンリクエストの構造体
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token"`
}

// RefreshTokenResponse リフレッシュトークンレスポンスの構造体
type RefreshTokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// RefreshTokenRepository リフレッシュトークンリポジトリのインターフェース
type RefreshTokenRepository interface {
	// リフレッシュトークンを保存
	Create(token *RefreshToken) error
	// トークン文字列でリフレッシュトークンを取得
	GetByToken(token string) (*RefreshToken, error)
	// ユーザーIDでリフレッシュトークンを取得
	GetByUserID(userID int) ([]*RefreshToken, error)
	// リフレッシュトークンを更新
	Update(token *RefreshToken) error
	// リフレッシュトークンを無効化
	Revoke(tokenID int) error
	// ユーザーのすべてのリフレッシュトークンを無効化
	RevokeAllByUserID(userID int) error
	// 期限切れのトークンを削除
	DeleteExpired() error
}

// TokenPair アクセストークンとリフレッシュトークンのペア
type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"` // アクセストークンの有効期限（秒）
}
