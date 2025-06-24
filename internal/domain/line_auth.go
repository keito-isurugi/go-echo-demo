package domain

import (
	"golang.org/x/oauth2"
)

// LineUser LINEユーザー情報
type LineUser struct {
	UserID        string `json:"userId"`
	DisplayName   string `json:"displayName"`
	PictureURL    string `json:"pictureUrl"`
	StatusMessage string `json:"statusMessage"`
}

// LineAuthConfig LINE認証設定の構造体
type LineAuthConfig struct {
	ChannelID     string
	ChannelSecret string
	RedirectURL   string
	Scopes        []string
}

// LineTokenResponse LINEトークンレスポンスの構造体
type LineTokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
	Scope        string `json:"scope"`
	IDToken      string `json:"id_token"`
}

// LineUserProfile LINEユーザープロフィールの構造体
type LineUserProfile struct {
	UserID        string `json:"userId"`
	DisplayName   string `json:"displayName"`
	PictureURL    string `json:"pictureUrl"`
	StatusMessage string `json:"statusMessage"`
}

// LineAuthRepository LINE認証リポジトリのインターフェース
type LineAuthRepository interface {
	GetOrCreateUser(lineUser *LineUser) (*User, error)
}

// LineAuthUsecase LINE認証ユースケースのインターフェース
type LineAuthUsecase interface {
	GetProviderName() string
	GetAuthURL() string
	ExchangeCodeForToken(code string) (*oauth2.Token, error)
	GetUserInfo(token *oauth2.Token) (*OAuthUser, error)
	Authenticate(code string) (*AuthResponse, error)
}
