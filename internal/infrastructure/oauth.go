package infrastructure

import (
	"database/sql"
	"log"
	"strings"

	"go-echo-demo/internal/domain"
	"go-echo-demo/internal/repository"
	"go-echo-demo/internal/usecase"
)

func NewOAuthRepository(db *sql.DB) domain.OAuthRepository {
	return repository.NewOAuthRepository(db)
}

// OAuthプロバイダーのマップを作成
func NewOAuthProviders(oauthRepo domain.OAuthRepository, authUsecase domain.AuthUsecase, stateManager domain.StateManager) map[string]domain.OAuthUsecase {
	providers := make(map[string]domain.OAuthUsecase)

	// Google認証を追加
	googleClientID := getEnv("GOOGLE_CLIENT_ID", "")
	log.Printf("Google Client ID: %s", googleClientID)

	if googleClientID != "" {
		config := domain.GoogleAuthConfig{
			ClientID:     googleClientID,
			ClientSecret: getEnv("GOOGLE_CLIENT_SECRET", ""),
			RedirectURL:  getEnv("GOOGLE_REDIRECT_URL", "http://localhost:8080/auth/google/callback"),
			Scopes:       []string{"openid", "email", "profile"},
		}
		googleAuth := usecase.NewGoogleAuthUsecase(config, oauthRepo, authUsecase, stateManager)
		providers["google"] = googleAuth
		log.Printf("Google OAuth provider initialized")
	} else {
		log.Printf("Google Client ID not found, skipping Google OAuth initialization")
	}

	// LINE認証を追加
	lineChannelID := getEnv("LINE_CHANNEL_ID", "")
	log.Printf("LINE Channel ID: %s", lineChannelID)

	if lineChannelID != "" {
		config := domain.LineAuthConfig{
			ChannelID:     lineChannelID,
			ChannelSecret: getEnv("LINE_CHANNEL_SECRET", ""),
			RedirectURL:   getEnv("LINE_CALLBACK_URL", "http://localhost:8080/auth/line/callback"),
			Scopes:        strings.Split(getEnv("LINE_SCOPES", "profile"), ","),
		}
		lineAuth := usecase.NewLineAuthUsecase(config, oauthRepo, authUsecase, stateManager)
		providers["line"] = lineAuth
		log.Printf("LINE OAuth provider initialized")
	} else {
		log.Printf("LINE Channel ID not found, skipping LINE OAuth initialization")
	}

	log.Printf("Total OAuth providers initialized: %d", len(providers))
	for providerName := range providers {
		log.Printf("Provider: %s", providerName)
	}

	return providers
}
