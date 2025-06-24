package domain

import (
	"golang.org/x/oauth2"
)

// OAuthProvider OAuthプロバイダーの共通インターフェース
type OAuthProvider interface {
	GetProviderName() string
	GetAuthURL() string
	ExchangeCodeForToken(code string) (*oauth2.Token, error)
	GetUserInfo(token *oauth2.Token) (*OAuthUser, error)
	Authenticate(code string) (*AuthResponse, error)
}

// OAuthUser OAuthユーザー情報の共通構造体
type OAuthUser struct {
	ProviderID   string `json:"provider_id"`
	ProviderName string `json:"provider_name"`
	Email        string `json:"email"`
	Name         string `json:"name"`
	Picture      string `json:"picture"`
	Verified     bool   `json:"verified"`
}

// OAuthConfig OAuth設定の共通構造体
type OAuthConfig struct {
	ClientID     string
	ClientSecret string
	RedirectURL  string
	Scopes       []string
}

// OAuthRepository OAuthリポジトリの共通インターフェース
type OAuthRepository interface {
	GetOrCreateUser(oauthUser *OAuthUser) (*User, error)
}

// OAuthUsecase OAuthユースケースの共通インターフェース
type OAuthUsecase interface {
	GetProviderName() string
	GetAuthURL() string
	GetUserInfo(token interface{}) (*OAuthUser, error)
	Authenticate(code string, state string) (*AuthResponse, error)
}
