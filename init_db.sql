-- データベース初期化スクリプト

-- usersテーブルの作成（OAuthプロバイダー対応）
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    password VARCHAR(100),
    provider_id VARCHAR(100),
    provider_name VARCHAR(50),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (provider_id, provider_name)
);

-- rolesテーブルの作成
CREATE TABLE IF NOT EXISTS roles (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) UNIQUE NOT NULL,
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- permissionsテーブルの作成
CREATE TABLE IF NOT EXISTS permissions (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) UNIQUE NOT NULL,
    description TEXT,
    resource VARCHAR(100) NOT NULL,
    action VARCHAR(50) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- user_rolesテーブルの作成（ユーザーとロールの関連）
CREATE TABLE IF NOT EXISTS user_roles (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    role_id INTEGER NOT NULL REFERENCES roles(id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(user_id, role_id)
);

-- role_permissionsテーブルの作成（ロールと権限の関連）
CREATE TABLE IF NOT EXISTS role_permissions (
    id SERIAL PRIMARY KEY,
    role_id INTEGER NOT NULL REFERENCES roles(id) ON DELETE CASCADE,
    permission_id INTEGER NOT NULL REFERENCES permissions(id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(role_id, permission_id)
);

-- テストユーザーの追加
INSERT INTO users (name, email, password) VALUES 
    ('テストユーザー1', 'user1@example.com', 'password123'),
    ('テストユーザー2', 'user2@example.com', 'password456')
ON CONFLICT (email) DO NOTHING;

-- 基本ロールの追加
INSERT INTO roles (name, description) VALUES 
    ('admin', '管理者：すべての権限を持つ'),
    ('user', '一般ユーザー：基本的な権限を持つ'),
    ('guest', 'ゲスト：読み取り専用権限')
ON CONFLICT (name) DO NOTHING;

-- 基本権限の追加
INSERT INTO permissions (name, description, resource, action) VALUES 
    ('user:read', 'ユーザー情報の読み取り', 'user', 'read'),
    ('user:write', 'ユーザー情報の作成・更新', 'user', 'write'),
    ('user:delete', 'ユーザーの削除', 'user', 'delete'),
    ('admin:read', '管理者機能の読み取り', 'admin', 'read'),
    ('admin:write', '管理者機能の作成・更新', 'admin', 'write'),
    ('admin:delete', '管理者機能の削除', 'admin', 'delete'),
    ('content:read', 'コンテンツの読み取り', 'content', 'read'),
    ('content:write', 'コンテンツの作成・更新', 'content', 'write'),
    ('content:delete', 'コンテンツの削除', 'content', 'delete')
ON CONFLICT (name) DO NOTHING;

-- ロールと権限の関連付け
-- adminロール：すべての権限
INSERT INTO role_permissions (role_id, permission_id) 
SELECT r.id, p.id FROM roles r, permissions p 
WHERE r.name = 'admin'
ON CONFLICT (role_id, permission_id) DO NOTHING;

-- userロール：基本的な権限
INSERT INTO role_permissions (role_id, permission_id) 
SELECT r.id, p.id FROM roles r, permissions p 
WHERE r.name = 'user' AND p.name IN ('user:read', 'user:write', 'content:read', 'content:write')
ON CONFLICT (role_id, permission_id) DO NOTHING;

-- guestロール：読み取り専用権限
INSERT INTO role_permissions (role_id, permission_id) 
SELECT r.id, p.id FROM roles r, permissions p 
WHERE r.name = 'guest' AND p.name IN ('content:read')
ON CONFLICT (role_id, permission_id) DO NOTHING;

-- テストユーザーにロールを割り当て
INSERT INTO user_roles (user_id, role_id) 
SELECT u.id, r.id FROM users u, roles r 
WHERE u.email = 'user1@example.com' AND r.name = 'admin'
ON CONFLICT (user_id, role_id) DO NOTHING;

INSERT INTO user_roles (user_id, role_id) 
SELECT u.id, r.id FROM users u, roles r 
WHERE u.email = 'user2@example.com' AND r.name = 'user'
ON CONFLICT (user_id, role_id) DO NOTHING; 
