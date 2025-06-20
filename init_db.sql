-- データベース初期化スクリプト

-- usersテーブルの作成（OAuthプロバイダー対応）
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    password VARCHAR(100),
    provider_id VARCHAR(100),
    provider_name VARCHAR(50)
);

-- プロバイダーIDとプロバイダー名の複合ユニーク制約
ALTER TABLE users ADD CONSTRAINT unique_provider_user UNIQUE (provider_id, provider_name);

-- テストユーザーの追加
INSERT INTO users (name, email, password) VALUES 
    ('テストユーザー1', 'user1@example.com', 'password123'),
    ('テストユーザー2', 'user2@example.com', 'password456')
ON CONFLICT (email) DO NOTHING; 