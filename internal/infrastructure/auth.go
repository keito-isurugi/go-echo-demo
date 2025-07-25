package infrastructure

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"strconv"
	"sync"
	"time"

	"go-echo-demo/internal/domain"
	"go-echo-demo/internal/repository"
	"go-echo-demo/internal/usecase"
)

func NewAuthRepository(db *sql.DB) domain.AuthRepository {
	return repository.NewAuthRepository(db)
}

func NewRefreshTokenRepository(db *sql.DB) domain.RefreshTokenRepository {
	return repository.NewRefreshTokenRepository(db)
}

func NewAuthUsecase(authRepo domain.AuthRepository, refreshTokenRepo domain.RefreshTokenRepository, userRepo domain.UserRepository) domain.AuthUsecase {
	// JWT有効期限の設定（デフォルト: 15分）
	jwtDurationStr := getEnv("JWT_DURATION_MINUTES", "15")
	jwtDurationMinutes, _ := strconv.Atoi(jwtDurationStr)
	
	// リフレッシュトークン有効期限の設定（デフォルト: 7日）
	refreshDurationStr := getEnv("REFRESH_TOKEN_DURATION_DAYS", "7")
	refreshDurationDays, _ := strconv.Atoi(refreshDurationStr)
	
	jwtConfig := domain.JWTConfig{
		SecretKey:            getEnv("JWT_SECRET_KEY", "your-secret-key"),
		Duration:             time.Duration(jwtDurationMinutes) * time.Minute,
		RefreshTokenDuration: time.Duration(refreshDurationDays) * 24 * time.Hour,
	}

	return usecase.NewAuthUsecase(authRepo, refreshTokenRepo, userRepo, jwtConfig)
}

// StateManagerImpl stateパラメータの実装
type StateManagerImpl struct {
	states map[string]time.Time
	mutex  sync.RWMutex
}

func NewStateManager() domain.StateManager {
	sm := &StateManagerImpl{
		states: make(map[string]time.Time),
	}

	// 古いstateを定期的にクリーンアップ
	go sm.cleanupExpiredStates()

	return sm
}

func (sm *StateManagerImpl) GenerateState() string {
	// 32バイトのランダムな値を生成
	bytes := make([]byte, 32)
	rand.Read(bytes)
	state := hex.EncodeToString(bytes)

	sm.mutex.Lock()
	sm.states[state] = time.Now()
	sm.mutex.Unlock()

	return state
}

func (sm *StateManagerImpl) ValidateState(state string) bool {
	sm.mutex.RLock()
	timestamp, exists := sm.states[state]
	sm.mutex.RUnlock()

	if !exists {
		return false
	}

	// 10分以内のstateのみ有効
	if time.Since(timestamp) > 10*time.Minute {
		sm.mutex.Lock()
		delete(sm.states, state)
		sm.mutex.Unlock()
		return false
	}

	// 使用済みのstateを削除
	sm.mutex.Lock()
	delete(sm.states, state)
	sm.mutex.Unlock()

	return true
}

func (sm *StateManagerImpl) cleanupExpiredStates() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		sm.mutex.Lock()
		now := time.Now()
		for state, timestamp := range sm.states {
			if now.Sub(timestamp) > 10*time.Minute {
				delete(sm.states, state)
			}
		}
		sm.mutex.Unlock()
	}
}
