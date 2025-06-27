package usecase

import (
	"errors"
	"time"

	"go-echo-demo/internal/domain"

	"github.com/golang-jwt/jwt/v5"
)

type AuthUsecase struct {
	authRepo  domain.AuthRepository
	jwtConfig domain.JWTConfig
}

func NewAuthUsecase(authRepo domain.AuthRepository, jwtConfig domain.JWTConfig) domain.AuthUsecase {
	return &AuthUsecase{
		authRepo:  authRepo,
		jwtConfig: jwtConfig,
	}
}

func (u *AuthUsecase) Login(email, password string) (*domain.AuthResponse, error) {
	user, err := u.authRepo.ValidateCredentials(email, password)
	if err != nil {
		return nil, err
	}

	token, err := u.GenerateToken(user)
	if err != nil {
		return nil, err
	}

	return &domain.AuthResponse{
		Token: token,
		User:  *user,
	}, nil
}

func (u *AuthUsecase) GenerateToken(user *domain.User) (string, error) {
	claims := domain.Claims{
		UserID: user.ID,
		Email:  user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
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
