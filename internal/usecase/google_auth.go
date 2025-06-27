package usecase

import (
	"context"
	"fmt"
	"log"

	"go-echo-demo/internal/domain"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	googleoauth2 "google.golang.org/api/oauth2/v2"
	"google.golang.org/api/option"
)

type GoogleAuthUsecase struct {
	config       *oauth2.Config
	oauthRepo    domain.OAuthRepository
	authUsecase  domain.AuthUsecase
	stateManager domain.StateManager
}

func NewGoogleAuthUsecase(config domain.GoogleAuthConfig, oauthRepo domain.OAuthRepository, authUsecase domain.AuthUsecase, stateManager domain.StateManager) domain.OAuthUsecase {
	oauthConfig := &oauth2.Config{
		ClientID:     config.ClientID,
		ClientSecret: config.ClientSecret,
		RedirectURL:  config.RedirectURL,
		Scopes:       config.Scopes,
		Endpoint:     google.Endpoint,
	}

	return &GoogleAuthUsecase{
		config:       oauthConfig,
		oauthRepo:    oauthRepo,
		authUsecase:  authUsecase,
		stateManager: stateManager,
	}
}

func (u *GoogleAuthUsecase) GetProviderName() string {
	return "google"
}

func (u *GoogleAuthUsecase) GetAuthURL() string {
	state := u.stateManager.GenerateState()
	return u.config.AuthCodeURL(state)
}

func (u *GoogleAuthUsecase) ExchangeCodeForToken(code string) (*oauth2.Token, error) {
	log.Printf("Exchanging code for token...")
	ctx := context.Background()
	token, err := u.config.Exchange(ctx, code)
	if err != nil {
		log.Printf("Failed to exchange code for token: %v", err)
		return nil, fmt.Errorf("failed to exchange code for token: %w", err)
	}
	log.Printf("Token exchange successful")
	return token, nil
}

func (u *GoogleAuthUsecase) GetUserInfo(token interface{}) (*domain.OAuthUser, error) {
	log.Printf("Getting user info from Google...")
	ctx := context.Background()

	// トークンをoauth2.Tokenにキャスト
	oauthToken, ok := token.(*oauth2.Token)
	if !ok {
		return nil, fmt.Errorf("invalid token type")
	}

	// OAuth2サービスを作成
	oauth2Service, err := googleoauth2.NewService(ctx, option.WithTokenSource(u.config.TokenSource(ctx, oauthToken)))
	if err != nil {
		log.Printf("Failed to create oauth2 service: %v", err)
		return nil, fmt.Errorf("failed to create oauth2 service: %w", err)
	}

	// ユーザー情報を取得
	userInfo, err := oauth2Service.Userinfo.Get().Do()
	if err != nil {
		log.Printf("Failed to get user info: %v", err)
		return nil, fmt.Errorf("failed to get user info: %w", err)
	}

	log.Printf("User info retrieved: %s (%s)", userInfo.Name, userInfo.Email)

	oauthUser := &domain.OAuthUser{
		ProviderID:   userInfo.Id,
		ProviderName: "google",
		Email:        userInfo.Email,
		Name:         userInfo.Name,
		Picture:      userInfo.Picture,
		Verified:     userInfo.VerifiedEmail != nil && *userInfo.VerifiedEmail,
	}

	return oauthUser, nil
}

func (u *GoogleAuthUsecase) Authenticate(code string, state string) (*domain.AuthResponse, error) {
	log.Printf("Starting Google authentication...")

	// stateパラメータを検証
	if !u.stateManager.ValidateState(state) {
		return nil, fmt.Errorf("invalid state parameter")
	}

	// 認証コードをトークンに交換
	token, err := u.ExchangeCodeForToken(code)
	if err != nil {
		log.Printf("Token exchange failed: %v", err)
		return nil, err
	}

	// Googleユーザー情報を取得
	oauthUser, err := u.GetUserInfo(token)
	if err != nil {
		log.Printf("Get user info failed: %v", err)
		return nil, err
	}

	log.Printf("Creating or getting user from database...")
	// データベースからユーザーを取得または作成
	user, err := u.oauthRepo.GetOrCreateUser(oauthUser)
	if err != nil {
		log.Printf("Get or create user failed: %v", err)
		return nil, err
	}

	log.Printf("Generating JWT token...")
	// JWTトークンを生成
	jwtToken, err := u.authUsecase.GenerateToken(user)
	if err != nil {
		log.Printf("Generate token failed: %v", err)
		return nil, err
	}

	log.Printf("Google authentication completed successfully")
	return &domain.AuthResponse{
		Token: jwtToken,
		User:  *user,
	}, nil
}
