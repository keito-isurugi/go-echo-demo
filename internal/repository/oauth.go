package repository

import (
	"database/sql"

	"go-echo-demo/internal/domain"
)

type OAuthRepository struct {
	db *sql.DB
}

func NewOAuthRepository(db *sql.DB) domain.OAuthRepository {
	return &OAuthRepository{db: db}
}

func (r *OAuthRepository) GetOrCreateUser(oauthUser *domain.OAuthUser) (*domain.User, error) {
	// 既存のユーザーを検索（プロバイダーIDまたはメールアドレスで）
	var user domain.User
	query := `SELECT id, name, email, password FROM users WHERE email = $1 OR provider_id = $2 AND provider_name = $3`

	err := r.db.QueryRow(query, oauthUser.Email, oauthUser.ProviderID, oauthUser.ProviderName).Scan(&user.ID, &user.Name, &user.Email, &user.Password)

	if err == sql.ErrNoRows {
		// ユーザーが存在しない場合は新規作成
		insertQuery := `INSERT INTO users (name, email, provider_id, provider_name, password) VALUES ($1, $2, $3, $4, $5) RETURNING id`
		err = r.db.QueryRow(insertQuery, oauthUser.Name, oauthUser.Email, oauthUser.ProviderID, oauthUser.ProviderName, "").Scan(&user.ID)
		if err != nil {
			return nil, err
		}

		user.Name = oauthUser.Name
		user.Email = oauthUser.Email
		user.Password = "" // OAuth認証ユーザーはパスワード不要

		return &user, nil
	}

	if err != nil {
		return nil, err
	}

	// 既存ユーザーの場合、プロバイダー情報を更新（まだ設定されていない場合）
	if user.Password == "" {
		// OAuth認証ユーザーの場合、プロバイダー情報を更新
		updateQuery := `UPDATE users SET provider_id = $1, provider_name = $2 WHERE id = $3`
		_, err = r.db.Exec(updateQuery, oauthUser.ProviderID, oauthUser.ProviderName, user.ID)
		if err != nil {
			return nil, err
		}
	}

	return &user, nil
}
