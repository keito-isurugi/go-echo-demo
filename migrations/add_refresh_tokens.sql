-- リフレッシュトークンテーブルの作成
-- このテーブルは、ユーザーのリフレッシュトークンを管理します
CREATE TABLE IF NOT EXISTS refresh_tokens (
    id SERIAL PRIMARY KEY,
    -- ユーザーID（外部キー）
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    -- リフレッシュトークン（ユニークで長い文字列）
    token VARCHAR(255) NOT NULL UNIQUE,
    -- アクセストークンのJTI（JWT ID）- トークンペアの追跡用
    access_token_jti VARCHAR(255),
    -- 有効期限
    expires_at TIMESTAMP NOT NULL,
    -- 作成日時
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    -- 更新日時
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    -- 最後に使用された日時
    last_used_at TIMESTAMP,
    -- 無効化されたかどうか
    revoked BOOLEAN DEFAULT FALSE,
    -- 無効化された日時
    revoked_at TIMESTAMP,
    -- デバイス情報（オプション）
    device_info VARCHAR(500),
    -- IPアドレス（オプション）
    ip_address VARCHAR(45)
);

-- インデックスの作成
-- トークンによる高速検索のため
CREATE INDEX idx_refresh_tokens_token ON refresh_tokens(token);
-- ユーザーIDによる検索のため
CREATE INDEX idx_refresh_tokens_user_id ON refresh_tokens(user_id);
-- 有効期限によるクリーンアップのため
CREATE INDEX idx_refresh_tokens_expires_at ON refresh_tokens(expires_at);
-- 無効化されていないトークンの検索のため
CREATE INDEX idx_refresh_tokens_revoked ON refresh_tokens(revoked);

-- 更新日時の自動更新トリガー
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_refresh_tokens_updated_at BEFORE UPDATE
    ON refresh_tokens FOR EACH ROW EXECUTE PROCEDURE update_updated_at_column();

-- コメントの追加
COMMENT ON TABLE refresh_tokens IS 'ユーザーのリフレッシュトークンを管理するテーブル';
COMMENT ON COLUMN refresh_tokens.user_id IS 'トークンを所有するユーザーのID';
COMMENT ON COLUMN refresh_tokens.token IS 'リフレッシュトークン本体（ハッシュ化して保存することを推奨）';
COMMENT ON COLUMN refresh_tokens.access_token_jti IS '関連するアクセストークンのJWT ID';
COMMENT ON COLUMN refresh_tokens.expires_at IS 'トークンの有効期限';
COMMENT ON COLUMN refresh_tokens.last_used_at IS '最後にトークンが使用された日時';
COMMENT ON COLUMN refresh_tokens.revoked IS 'トークンが無効化されているかどうか';
COMMENT ON COLUMN refresh_tokens.device_info IS 'トークンを生成したデバイスの情報';
COMMENT ON COLUMN refresh_tokens.ip_address IS 'トークンを生成したIPアドレス';