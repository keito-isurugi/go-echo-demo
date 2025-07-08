package repository

import (
	"database/sql"
	"time"

	"github.com/koukikitamura/go-echo-demo/internal/domain"
)

// refreshTokenRepository リフレッシュトークンリポジトリの実装
type refreshTokenRepository struct {
	db *sql.DB
}

// NewRefreshTokenRepository リフレッシュトークンリポジトリのコンストラクタ
func NewRefreshTokenRepository(db *sql.DB) domain.RefreshTokenRepository {
	return &refreshTokenRepository{db: db}
}

// Create リフレッシュトークンを保存
func (r *refreshTokenRepository) Create(token *domain.RefreshToken) error {
	query := `
		INSERT INTO refresh_tokens (
			user_id, token, access_token_jti, expires_at,
			device_info, ip_address, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id`

	err := r.db.QueryRow(
		query,
		token.UserID,
		token.Token,
		token.AccessTokenJTI,
		token.ExpiresAt,
		token.DeviceInfo,
		token.IPAddress,
		time.Now(),
		time.Now(),
	).Scan(&token.ID)

	return err
}

// GetByToken トークン文字列でリフレッシュトークンを取得
func (r *refreshTokenRepository) GetByToken(tokenString string) (*domain.RefreshToken, error) {
	query := `
		SELECT id, user_id, token, access_token_jti, expires_at,
			   created_at, updated_at, last_used_at, revoked, revoked_at,
			   device_info, ip_address
		FROM refresh_tokens
		WHERE token = $1 AND revoked = false`

	var token domain.RefreshToken
	err := r.db.QueryRow(query, tokenString).Scan(
		&token.ID,
		&token.UserID,
		&token.Token,
		&token.AccessTokenJTI,
		&token.ExpiresAt,
		&token.CreatedAt,
		&token.UpdatedAt,
		&token.LastUsedAt,
		&token.Revoked,
		&token.RevokedAt,
		&token.DeviceInfo,
		&token.IPAddress,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	return &token, err
}

// GetByUserID ユーザーIDでリフレッシュトークンを取得
func (r *refreshTokenRepository) GetByUserID(userID int) ([]*domain.RefreshToken, error) {
	query := `
		SELECT id, user_id, token, access_token_jti, expires_at,
			   created_at, updated_at, last_used_at, revoked, revoked_at,
			   device_info, ip_address
		FROM refresh_tokens
		WHERE user_id = $1 AND revoked = false
		ORDER BY created_at DESC`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tokens []*domain.RefreshToken
	for rows.Next() {
		var token domain.RefreshToken
		err := rows.Scan(
			&token.ID,
			&token.UserID,
			&token.Token,
			&token.AccessTokenJTI,
			&token.ExpiresAt,
			&token.CreatedAt,
			&token.UpdatedAt,
			&token.LastUsedAt,
			&token.Revoked,
			&token.RevokedAt,
			&token.DeviceInfo,
			&token.IPAddress,
		)
		if err != nil {
			return nil, err
		}
		tokens = append(tokens, &token)
	}

	return tokens, nil
}

// Update リフレッシュトークンを更新
func (r *refreshTokenRepository) Update(token *domain.RefreshToken) error {
	query := `
		UPDATE refresh_tokens
		SET access_token_jti = $2, last_used_at = $3, updated_at = $4
		WHERE id = $1`

	_, err := r.db.Exec(
		query,
		token.ID,
		token.AccessTokenJTI,
		time.Now(),
		time.Now(),
	)

	return err
}

// Revoke リフレッシュトークンを無効化
func (r *refreshTokenRepository) Revoke(tokenID int) error {
	query := `
		UPDATE refresh_tokens
		SET revoked = true, revoked_at = $2, updated_at = $3
		WHERE id = $1`

	_, err := r.db.Exec(query, tokenID, time.Now(), time.Now())
	return err
}

// RevokeAllByUserID ユーザーのすべてのリフレッシュトークンを無効化
func (r *refreshTokenRepository) RevokeAllByUserID(userID int) error {
	query := `
		UPDATE refresh_tokens
		SET revoked = true, revoked_at = $2, updated_at = $3
		WHERE user_id = $1 AND revoked = false`

	_, err := r.db.Exec(query, userID, time.Now(), time.Now())
	return err
}

// DeleteExpired 期限切れのトークンを削除
func (r *refreshTokenRepository) DeleteExpired() error {
	query := `
		DELETE FROM refresh_tokens
		WHERE expires_at < $1 OR (revoked = true AND revoked_at < $2)`

	// 期限切れまたは7日以上前に無効化されたトークンを削除
	sevenDaysAgo := time.Now().Add(-7 * 24 * time.Hour)
	_, err := r.db.Exec(query, time.Now(), sevenDaysAgo)
	return err
}