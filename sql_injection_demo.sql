-- SQLインジェクションデモ用のテーブル作成とデータ投入

-- productsテーブルの作成
CREATE TABLE IF NOT EXISTS products (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    price DECIMAL(10, 2) NOT NULL,
    category VARCHAR(50) NOT NULL,
    stock INT NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- テストデータの投入
INSERT INTO products (name, description, price, category, stock) VALUES
    ('ノートパソコン', '高性能な15インチノートPC', 120000.00, 'Electronics', 10),
    ('ワイヤレスマウス', 'Bluetooth対応の静音マウス', 3500.00, 'Electronics', 50),
    ('USBメモリ 32GB', '高速転送対応USBメモリ', 2000.00, 'Electronics', 100),
    ('オフィスチェア', '人間工学に基づいた快適な椅子', 45000.00, 'Furniture', 5),
    ('スタンディングデスク', '高さ調整可能なデスク', 80000.00, 'Furniture', 3),
    ('コーヒーメーカー', '全自動コーヒーメーカー', 25000.00, 'Appliances', 15),
    ('電気ケトル', '1.2L容量の電気ケトル', 5000.00, 'Appliances', 30),
    ('プログラミング入門書', 'Go言語によるWebアプリケーション開発', 3500.00, 'Books', 20),
    ('データベース設計の本', 'SQLとNoSQLの使い分け', 4200.00, 'Books', 15),
    ('セキュリティ対策ガイド', 'Webアプリケーションのセキュリティ', 3800.00, 'Books', 25);

-- 更新トリガーの作成
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_products_updated_at BEFORE UPDATE
    ON products FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();