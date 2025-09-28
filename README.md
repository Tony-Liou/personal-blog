# 個人部落格專案計畫 (Go/Gin + Svelte)

這是一個使用 Go (Gin) 作為後端，Svelte (SvelteKit) 作為前端的個人部落格網站的初步開發計畫。

專案儲存庫: https://github.com/Tony-Liou

## 1. 技術棧 (Tech Stack)

- 後端 (Backend): Go + Gin Web Framework
- 前端 (Frontend): Svelte + SvelteKit
- 資料庫 (Database): PostgreSQL
- 身份驗證 (Authentication): JWT (JSON Web Tokens)
- 部署 (Deployment): Docker (建議)

## 2. 如何執行後端 (How to Run Backend)

1. 安裝 Go: 請確保你的系統已安裝 Go (建議版本 1.22 或以上)。
2. 設定 PostgreSQL: 建立一個新的 PostgreSQL 資料庫。
3. 設定環境變數:
   - 複製 `.env.example` 檔案并重新命名為 `.env`。
   - 在 `.env` 檔案中填入你的 PostgreSQL 資料庫連線資訊和 `JWT_SECRET`。`JWT_SECRET` 必須設定一個隨機的長字串！
4. 安裝依賴: 在專案根目錄下執行 `go mod tidy`。
5. 啟動伺服器: 執行 `go run main.go`。
6. 伺服器將會在你設定的 `API_PORT` (預設為 8080) 上啟動。

## 3. 如何執行前端 (How to Run Frontend)

1. 安裝 Node.js: 請確保你的系統已安裝 Node.js (建議版本 18 或以上)。
2. 切換到前端目錄: `cd frontend`
3. 安裝依賴: 執行 `npm install`。
4. 啟動開發伺服器: 執行 `npm run dev`。
5. SvelteKit 開發伺服器將會啟動，你可以在瀏覽器中開啟 http://localhost:5173 (或終端機中顯示的 URL) 來查看網站。

**重要提示**: 請確保後端 Go 伺服器正在運行，因為前端需要從後端 API (http://localhost:8080) 獲取資料。

### 測試使用者驗證 (可使用 Postman 或 curl)

A. 註冊新使用者
   - 方法: POST
   - URL: http://localhost:8080/api/v1/auth/signup
   - Body (raw, JSON):
   ```json
   {
      "username": "admin",
      "email": "admin@example.com",
      "password": "a_very_strong_password"
   }
   ```
B. 登入取得 Token
   - 方法: POST
   - URL: http://localhost:8080/api/v1/auth/login
   - Body (raw, JSON):
   ```json
   {
      "username": "admin",
      "password": "a_very_strong_password"
   }
   ```
   - 成功回應:
   ```json
   {
      "token": "eyJhbGciOiJI..."
   }
   ```
## 3. 核心功能 (Core Features)

### 使用者管理

- [x] 管理員登入/登出: 透過 JWT 進行身份驗證，保護後台管理介面。
- [x] 使用者角色: 簡化為單一「管理員」角色。

### 文章管理 (CRUD)

- [x] 建立文章: 提供一個富文本或 Markdown 編輯器來撰寫文章。
- [x] 讀取文章:
  - 在首頁顯示所有文章的列表（分頁）。
  - 點擊後可閱讀單篇文章的完整內容。
- [x] 更新文章: 已登入的管理員可以編輯已發布的文章。
- [x] 刪除文章: 已登入的管理員可以刪除文章。

### 檔案上傳

- [x] 圖片上傳: 在文章編輯器中，可以上傳圖片並將其插入文章內容。
- [x] 靜態檔案服務: 後端需要提供一個路由來存取上傳的檔案。

### 作者頁面

- [x] 作者資訊: 建立一個靜態頁面，用於顯示作者的個人簡介、聯絡方式和頭像。

## 4. 資料庫結構 (Database Schema)

主要需要兩個資料表：`users` 和 `posts`。

- `users` table: 儲存管理員資訊。
   - `id` (Primary Key)
   - `username` (UNIQUE, NOT NULL)
   - `password_hash` (NOT NULL)
   - `email` (UNIQUE)
   - `bio` (TEXT)
   - `avatar_url` (VARCHAR)
   - `created_at` (TIMESTAMP)
- `posts` table: 儲存部落格文章。
   - `id` (Primary Key)
   - `title` (VARCHAR, NOT NULL)
   - `content` (TEXT, NOT NULL)
   - `cover_image_url` (VARCHAR)
   - `author_id` (Foreign Key to users.id)
   - `created_at` (TIMESTAMP)
   - `updated_at` (TIMESTAMP)
   
## 5. API 端點設計 (API Endpoints)

所有 API 都以 /api/v1 為前綴。

| 方法 (Method) | 路徑 (Path)      | 描述                      | 需要驗證 (Auth?) |
| :------------ | :--------------- | :------------------------ | :--------------- |
| POST          | /auth/signup     | 建立新使用者              | 否 (No)          |
| POST          | /auth/login      | 管理員登入，取得 JWT      | 否 (No)          |
| GET           | /posts           | 取得所有文章列表（可分頁） | 否 (No)          |
| GET           | /posts/:id       | 取得單篇文章內容          | 否 (No)          |
| POST          | /posts           | 建立新文章                | 是 (Yes)         |
| PUT           | /posts/:id       | 更新指定文章              | 是 (Yes)         |
| DELETE        | /posts/:id       | 刪除指定文章              | 是 (Yes)         |
| POST          | /upload          | 上傳檔案（例如圖片）      | 是 (Yes)         |
| GET           | /author/:id      | 取得作者公開資訊          | 否 (No)          |

## 6. 前端結構 (Frontend Structure - SvelteKit)

```
src/
├── lib/
│   ├── components/
│   │   ├── Navbar.svelte       # 導覽列
│   │   └── PostCard.svelte     # 文章卡片元件
│   └── services/
│       └── api.js              # 封裝 fetch API
├── routes/
│   ├── +page.svelte            # 首頁 (文章列表)
│   ├── about/
│   │   └── +page.svelte        # 作者頁面
│   ├── posts/
│   │   └── [slug]/
│   │       └── +page.svelte    # 單篇文章頁面
│   ├── login/
│   │   └── +page.svelte        # 登入頁面
│   └── admin/
│       ├── +layout.svelte      # 後台佈局 (檢查登入狀態)
│       ├── +page.svelte        # 後台儀表板
│       ├── new/
│       │   └── +page.svelte    # 建立新文章
│       └── edit/
│           └── [slug]/
│               └── +page.svelte # 編輯文章
└── app.html
```

## 7. 開發步驟 (Roadmap)

1. 第一階段：後端基礎建設
   - [x] 初始化 Go 專案，安裝 Gin 和 PostgreSQL 驅動。
   - [x] 設計並建立資料庫 Schema (已透過 GORM AutoMigrate 實現)。
   - [x] 實作使用者登入邏輯，產生 JWT。
   - [x] 實作文章的 CRUD API 端點。
2. 第二階段：前端基礎建設
   - [x] 初始化 SvelteKit 專案。
   - [x] 建立基本頁面路由結構。
   - [x] 開發首頁，從後端 API 獲取並顯示文章列表。
   - [x] 開發文章內頁，顯示完整文章內容。
3. 第三階段：功能整合
   - [x] 後端：為需要權限的 API 加上 JWT 中介層 (Middleware)。
   - [x] 前端：開發登入頁面，登入成功後儲存 JWT。
   - [ ] 前端：建立後台管理介面，用於建立、編輯和刪除文章。
   - [ ] 前端：實作受保護的路由，只有登入的管理員才能進入後台。
4. 第四階段：進階功能與優化
   - [ ] 後端：實作檔案上傳 API，將檔案儲存至伺服器。
   - [ ] 前端：在文章編輯器中整合圖片上傳功能。
   - [ ] 開發作者頁面。
   - [ ] 前後端程式碼優化、錯誤處理和環境變數設定。
5. 第五階段：部署
   - [ ] 編寫 Dockerfile，將前後端打包成容器。
   - [ ] 使用 docker-compose 進行本地部署測試。
   - [ ] 將專案部署到雲端伺服器。
