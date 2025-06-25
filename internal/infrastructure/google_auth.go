package infrastructure

import (
	"database/sql"

	"go-echo-demo/internal/domain"
	"go-echo-demo/internal/repository"
	"go-echo-demo/internal/usecase"
)

// NewGoogleAuthRepository Google認証リポジトリを作成
func NewGoogleAuthRepository(db *sql.DB) domain.OAuthRepository {
	return repository.NewGoogleAuthRepository(db)
}

// NewGoogleAuthUsecase Google認証ユースケースを作成
func NewGoogleAuthUsecase(oauthRepo domain.OAuthRepository, authUsecase domain.AuthUsecase, stateManager domain.StateManager) domain.OAuthUsecase {
	config := domain.GoogleAuthConfig{
		ClientID:     getEnv("GOOGLE_CLIENT_ID", ""),
		ClientSecret: getEnv("GOOGLE_CLIENT_SECRET", ""),
		RedirectURL:  getEnv("GOOGLE_REDIRECT_URL", "http://localhost:8080/auth/google/callback"),
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile"},
	}

	return usecase.NewGoogleAuthUsecase(config, oauthRepo, authUsecase, stateManager)
}
