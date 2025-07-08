package usecase

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"time"

	"go-echo-demo/internal/domain"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type AuthUsecase struct {
	authRepo         domain.AuthRepository
	refreshTokenRepo domain.RefreshTokenRepository
	userRepo         domain.UserRepository
	jwtConfig        domain.JWTConfig
}

func NewAuthUsecase(
	authRepo domain.AuthRepository,
	refreshTokenRepo domain.RefreshTokenRepository,
	userRepo domain.UserRepository,
	jwtConfig domain.JWTConfig,
) domain.AuthUsecase {
	return &AuthUsecase{
		authRepo:         authRepo,
		refreshTokenRepo: refreshTokenRepo,
		userRepo:         userRepo,
		jwtConfig:        jwtConfig,
	}
}

func (u *AuthUsecase) Login(email, password string) (*domain.AuthResponse, error) {
	user, err := u.authRepo.ValidateCredentials(email, password)
	if err != nil {
		return nil, err
	}

	// トークンペアを生成（デバイス情報とIPアドレスは後でハンドラーから渡すことも可能）
	tokenPair, err := u.GenerateTokenPair(user, "", "")
	if err != nil {
		return nil, err
	}

	return &domain.AuthResponse{
		Token:        tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		User:         *user,
	}, nil
}

func (u *AuthUsecase) GenerateToken(user *domain.User) (string, error) {
	// JWT IDを生成
	jti := uuid.New().String()
	
	claims := domain.Claims{
		UserID: user.ID,
		Email:  user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        jti, // JWT IDを追加
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(u.jwtConfig.Duration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(u.jwtConfig.SecretKey))
}

func (u *AuthUsecase) ValidateToken(tokenString string) (*domain.Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &domain.Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(u.jwtConfig.SecretKey), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*domain.Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

// GenerateTokenPair アクセストークンとリフレッシュトークンのペアを生成
func (u *AuthUsecase) GenerateTokenPair(user *domain.User, deviceInfo, ipAddress string) (*domain.TokenPair, error) {
	// アクセストークンを生成
	accessToken, err := u.GenerateToken(user)
	if err != nil {
		return nil, err
	}

	// アクセストークンからJWTIDを取得
	claims, err := u.ValidateToken(accessToken)
	if err != nil {
		return nil, err
	}

	// リフレッシュトークンを生成
	refreshTokenString, err := u.generateRefreshTokenString()
	if err != nil {
		return nil, err
	}

	// リフレッシュトークンをデータベースに保存
	refreshToken := &domain.RefreshToken{
		UserID:         user.ID,
		Token:          refreshTokenString,
		AccessTokenJTI: claims.ID,
		ExpiresAt:      time.Now().Add(u.jwtConfig.RefreshTokenDuration),
		DeviceInfo:     deviceInfo,
		IPAddress:      ipAddress,
	}

	if err := u.refreshTokenRepo.Create(refreshToken); err != nil {
		return nil, err
	}

	return &domain.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshTokenString,
		ExpiresIn:    int(u.jwtConfig.Duration.Seconds()),
	}, nil
}

// RefreshToken リフレッシュトークンを使用して新しいトークンペアを生成
func (u *AuthUsecase) RefreshToken(refreshTokenString string) (*domain.TokenPair, error) {
	// リフレッシュトークンを検証
	refreshToken, err := u.refreshTokenRepo.GetByToken(refreshTokenString)
	if err != nil {
		return nil, err
	}

	if refreshToken == nil {
		return nil, errors.New("invalid refresh token")
	}

	// 有効期限をチェック
	if time.Now().After(refreshToken.ExpiresAt) {
		return nil, errors.New("refresh token expired")
	}

	// ユーザー情報を取得
	user, err := u.userRepo.GetByID(refreshToken.UserID)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, errors.New("user not found")
	}

	// 新しいアクセストークンを生成
	newAccessToken, err := u.GenerateToken(user)
	if err != nil {
		return nil, err
	}

	// 新しいアクセストークンからJWTIDを取得
	claims, err := u.ValidateToken(newAccessToken)
	if err != nil {
		return nil, err
	}

	// リフレッシュトークンのJTIを更新
	refreshToken.AccessTokenJTI = claims.ID
	refreshToken.LastUsedAt = &[]time.Time{time.Now()}[0]
	
	if err := u.refreshTokenRepo.Update(refreshToken); err != nil {
		return nil, err
	}

	return &domain.TokenPair{
		AccessToken:  newAccessToken,
		RefreshToken: refreshTokenString, // 同じリフレッシュトークンを返す
		ExpiresIn:    int(u.jwtConfig.Duration.Seconds()),
	}, nil
}

// Logout ログアウト時にリフレッシュトークンを無効化
func (u *AuthUsecase) Logout(userID int) error {
	return u.refreshTokenRepo.RevokeAllByUserID(userID)
}

// generateRefreshTokenString セキュアなランダムトークンを生成
func (u *AuthUsecase) generateRefreshTokenString() (string, error) {
	b := make([]byte, 32) // 256ビットのランダム値
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}
