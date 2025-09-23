-- 為確保原子性操作，使用事務
BEGIN;

-- 使用者資料表 (用於儲存管理員資訊)
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    bio TEXT,
    avatar_url VARCHAR(255),
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

-- 文章資料表
CREATE TABLE IF NOT EXISTS posts (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    content TEXT NOT NULL,
    cover_image_url VARCHAR(255),
    author_id INT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    
    -- 設定外鍵關聯
    CONSTRAINT fk_author
        FOREIGN KEY(author_id) 
        REFERENCES users(id)
        ON DELETE CASCADE -- 如果刪除使用者，其所有文章也會被刪除
);

-- 建立一個觸發器，在更新 post 時自動更新 updated_at 欄位
CREATE OR REPLACE FUNCTION update_modified_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_posts_modtime
BEFORE UPDATE ON posts
FOR EACH ROW
EXECUTE FUNCTION update_modified_column();

-- 插入一個預設管理員（請在應用程式啟動時使用安全的密碼雜湊）
-- 密碼應該在後端應用程式中進行雜湊處理後再存入
INSERT INTO users (username, email, password_hash, bio) VALUES
('admin', 'admin@example.com', '在此處填入安全的密碼雜湊值', '這是網站管理員的個人簡介。');

COMMIT;
