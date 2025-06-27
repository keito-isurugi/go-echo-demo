package usecase

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"

	"go-echo-demo/internal/domain"
)

type LineAuthUsecase struct {
	config       domain.LineAuthConfig
	oauthRepo    domain.OAuthRepository
	authUsecase  domain.AuthUsecase
	stateManager domain.StateManager
}

func NewLineAuthUsecase(config domain.LineAuthConfig, oauthRepo domain.OAuthRepository, authUsecase domain.AuthUsecase, stateManager domain.StateManager) domain.OAuthUsecase {
	return &LineAuthUsecase{
		config:       config,
		oauthRepo:    oauthRepo,
		authUsecase:  authUsecase,
		stateManager: stateManager,
	}
}

func (u *LineAuthUsecase) GetProviderName() string {
	return "line"
}

func (u *LineAuthUsecase) GetAuthURL() string {
	state := u.stateManager.GenerateState()

	params := url.Values{}
	params.Add("response_type", "code")
	params.Add("client_id", u.config.ChannelID)
	params.Add("redirect_uri", u.config.RedirectURL)
	params.Add("state", state)
	params.Add("scope", strings.Join(u.config.Scopes, " "))

	return "https://access.line.me/oauth2/v2.1/authorize?" + params.Encode()
}

func (u *LineAuthUsecase) ExchangeCodeForToken(code string) (*domain.LineTokenResponse, error) {
	log.Printf("Exchanging LINE code for token...")

	data := url.Values{}
	data.Set("grant_type", "authorization_code")
	data.Set("code", code)
	data.Set("redirect_uri", u.config.RedirectURL)
	data.Set("client_id", u.config.ChannelID)
	data.Set("client_secret", u.config.ChannelSecret)

	req, err := http.NewRequest("POST", "https://api.line.me/oauth2/v2.1/token", strings.NewReader(data.Encode()))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange code for token: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("token exchange failed with status: %d", resp.StatusCode)
	}

	var tokenResp domain.LineTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return nil, fmt.Errorf("failed to decode token response: %w", err)
	}

	log.Printf("LINE token exchange successful")
	return &tokenResp, nil
}

func (u *LineAuthUsecase) GetUserInfo(token interface{}) (*domain.OAuthUser, error) {
	log.Printf("Getting user info from LINE...")

	// トークンをLineTokenResponseにキャスト
	lineToken, ok := token.(*domain.LineTokenResponse)
	if !ok {
		return nil, fmt.Errorf("invalid token type")
	}

	req, err := http.NewRequest("GET", "https://api.line.me/v2/profile", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+lineToken.AccessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("get user info failed with status: %d", resp.StatusCode)
	}

	var profile domain.LineUserProfile
	if err := json.NewDecoder(resp.Body).Decode(&profile); err != nil {
		return nil, fmt.Errorf("failed to decode user profile: %w", err)
	}

	log.Printf("LINE user info retrieved: %s (%s)", profile.DisplayName, profile.UserID)

	oauthUser := &domain.OAuthUser{
		ProviderID:   profile.UserID,
		ProviderName: "line",
		Email:        "", // LINEはメールアドレスを提供しない
		Name:         profile.DisplayName,
		Picture:      profile.PictureURL,
		Verified:     true, // LINEユーザーは基本的に認証済み
	}

	return oauthUser, nil
}

func (u *LineAuthUsecase) Authenticate(code string, state string) (*domain.AuthResponse, error) {
	log.Printf("Starting LINE authentication...")

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

	// LINEユーザー情報を取得
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

	log.Printf("LINE authentication completed successfully")
	return &domain.AuthResponse{
		Token: jwtToken,
		User:  *user,
	}, nil
}
