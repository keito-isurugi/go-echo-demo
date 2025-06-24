package domain

import (
	"golang.org/x/oauth2"
)

// GoogleUser Googleユーザー情報（レガシー、削除予定）
type GoogleUser struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	Name          string `json:"name"`
	Picture       string `json:"picture"`
	VerifiedEmail bool   `json:"verified_email"`
}

// GoogleAuthConfig Google OAuth設定
type GoogleAuthConfig struct {
	ClientID     string
	ClientSecret string
	RedirectURL  string
	Scopes       []string
}

// GoogleAuthRepository Google認証リポジトリのインターフェース（OAuthRepositoryに統合）
type GoogleAuthRepository interface {
	GetOrCreateUser(oauthUser *OAuthUser) (*User, error)
}

// GoogleAuthUsecase Google認証ユースケースのインターフェース（OAuthUsecaseに統合）
type GoogleAuthUsecase interface {
	GetProviderName() string
	GetAuthURL() string
	ExchangeCodeForToken(code string) (*oauth2.Token, error)
	GetUserInfo(token *oauth2.Token) (*OAuthUser, error)
	Authenticate(code string) (*AuthResponse, error)
}
