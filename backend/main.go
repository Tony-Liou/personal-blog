package main

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// --- 資料庫模型 (Database Models) ---

type User struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	Username     string    `gorm:"unique;not null;size:50" json:"username"`
	Email        string    `gorm:"unique;not null;size:255" json:"-"` // Don't expose email
	PasswordHash string    `gorm:"not null" json:"-"`                 // Don't expose password hash
	Bio          string    `json:"bio"`
	AvatarURL    string    `json:"avatar_url"`
	CreatedAt    time.Time `json:"created_at"`
	Posts        []Post    `gorm:"foreignKey:AuthorID" json:"-"`
}

type Post struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	Title         string    `gorm:"not null" json:"title"`
	Content       string    `gorm:"not null" json:"content"`
	CoverImageURL string    `json:"cover_image_url"`
	AuthorID      uint      `gorm:"not null" json:"author_id"`
	Author        User      `gorm:"foreignKey:AuthorID" json:"author"` // Preload author info
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// --- API Payload 結構 ---
type AuthPayload struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

type LoginPayload struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type PostPayload struct {
	Title         string `json:"title" binding:"required"`
	Content       string `json:"content" binding:"required"`
	CoverImageURL string `json:"cover_image_url"`
}

// --- 全域變數 ---
var db *gorm.DB
var err error
var jwtSecret []byte

// --- 初始化函式 ---

func init() {
	// 載入 .env 檔案
	if err := godotenv.Load(); err != nil {
		log.Println("找不到 .env 檔案，將使用預設環境變數")
	}

	// 設定 JWT Secret
	jwtSecret = []byte(os.Getenv("JWT_SECRET"))
	if len(jwtSecret) == 0 {
		log.Fatal("環境變數 JWT_SECRET 未設定")
	}

	// 連線到 PostgreSQL 資料庫
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		os.Getenv("DB_HOST"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"), os.Getenv("DB_PORT"), os.Getenv("DB_SSLMODE"),
	)
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("無法連線到資料庫: ", err)
	}

	// 自動遷移資料庫結構
	err = db.AutoMigrate(&User{}, &Post{})
	if err != nil {
		log.Fatal("資料庫遷移失敗: ", err)
	}
	log.Println("資料庫連線並遷移成功！")
}

// --- 主函式 ---

func main() {
	router := gin.Default()
	
	// 禁用自動重定向
	router.RedirectTrailingSlash = false
	
	// 設定 CORS
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))
	
	router.Static("/uploads", "./uploads")
	api := router.Group("/api/v1")
	{
		setupAuthRoutes(api)
		setupPostRoutes(api)
		setupAuthorRoutes(api)
		setupUploadRoutes(api)
	}
	port := os.Getenv("API_PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("伺服器正在 http://localhost:%s 上運行", port)
	err := router.Run(":" + port)
	if err != nil {
		return 
	}
}

// --- 路由設定 (Route Setup) ---

func setupAuthRoutes(group *gin.RouterGroup) {
	authRoutes := group.Group("/auth")
	{
		authRoutes.POST("/signup", handleSignUp)
		authRoutes.POST("/login", handleLogin)
	}
}

func setupPostRoutes(group *gin.RouterGroup) {
	postRoutes := group.Group("/posts")
	{
		// 公開路由
		postRoutes.GET("/", handleGetPosts)
		postRoutes.GET("/:id", handleGetPostByID)

		// 需要驗證的路由
		protected := postRoutes.Group("/")
		protected.Use(authMiddleware()) // 套用 JWT 中介層
		{
			protected.POST("/", handleCreatePost)
			protected.PUT("/:id", handleUpdatePost)
			protected.DELETE("/:id", handleDeletePost)
		}
	}
}

func setupAuthorRoutes(group *gin.RouterGroup) {
	group.GET("/author/:id", handleGetAuthor)
}

func setupUploadRoutes(group *gin.RouterGroup) {
	protected := group.Group("/upload")
	protected.Use(authMiddleware()) // 同樣保護檔案上傳
	{
		protected.POST("/", handleUploadFile)
	}
}

// --- API 處理函式 (Handlers) ---
func handleSignUp(c *gin.Context) {
	var payload AuthPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "無效的請求資料: " + err.Error()})
		return
	}

	// 檢查使用者名稱或 Email 是否已存在
	var existingUser User
	if err := db.Where("username = ? OR email = ?", payload.Username, payload.Email).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "使用者名稱或 Email 已經存在"})
		return
	}

	// 雜湊密碼
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(payload.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "無法處理密碼"})
		return
	}

	// 建立新使用者
	newUser := User{
		Username:     payload.Username,
		Email:        payload.Email,
		PasswordHash: string(hashedPassword),
	}

	result := db.Create(&newUser)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "無法建立使用者"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "使用者建立成功"})
}

// handleLogin 處理使用者登入
func handleLogin(c *gin.Context) {
	var payload LoginPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "無效的請求資料: " + err.Error()})
		return
	}

	var user User
	// 根據使用者名稱查詢使用者
	if err := db.Where("username = ?", payload.Username).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "無效的使用者名稱或密碼"})
		return
	}

	// 驗證密碼
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(payload.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "無效的使用者名稱或密碼"})
		return
	}

	// 產生 JWT
	token, err := generateJWT(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "無法產生 token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

// handleGetPosts 取得所有文章 (含分頁)
func handleGetPosts(c *gin.Context) {
	var posts []Post

	// 分頁邏輯
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset := (page - 1) * limit

	// Preload("Author") 會自動帶上作者資訊
	result := db.Preload("Author").Limit(limit).Offset(offset).Order("created_at desc").Find(&posts)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "讀取文章失敗"})
		return
	}

	c.JSON(http.StatusOK, posts)
}

// handleGetPostByID 根據 ID 取得單篇文章
func handleGetPostByID(c *gin.Context) {
	id := c.Param("id")
	var post Post

	if err := db.Preload("Author").First(&post, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "找不到文章"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "讀取文章失敗"})
		return
	}

	c.JSON(http.StatusOK, post)
}

// handleCreatePost 建立新文章
func handleCreatePost(c *gin.Context) {
	var payload PostPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 從中介層取得 userID
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授權的操作"})
		return
	}

	newPost := Post{
		Title:         payload.Title,
		Content:       payload.Content,
		CoverImageURL: payload.CoverImageURL,
		AuthorID:      userID.(uint),
	}

	if err := db.Create(&newPost).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "建立文章失敗"})
		return
	}

	// 重新查詢文章以包含作者資訊
	db.Preload("Author").First(&newPost, newPost.ID)

	c.JSON(http.StatusCreated, newPost)
}

// handleUpdatePost 更新文章
func handleUpdatePost(c *gin.Context) {
	id := c.Param("id")
	var post Post

	// 檢查文章是否存在
	if err := db.First(&post, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "找不到文章"})
		return
	}

	// 檢查是否為作者本人
	userID, _ := c.Get("userID")
	if post.AuthorID != userID.(uint) {
		c.JSON(http.StatusForbidden, gin.H{"error": "你沒有權限修改此文章"})
		return
	}

	var payload PostPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 更新文章
	post.Title = payload.Title
	post.Content = payload.Content
	post.CoverImageURL = payload.CoverImageURL
	db.Save(&post)

	db.Preload("Author").First(&post, post.ID)

	c.JSON(http.StatusOK, post)
}

// handleDeletePost 刪除文章
func handleDeletePost(c *gin.Context) {
	id := c.Param("id")
	var post Post

	// 檢查文章是否存在
	if err := db.First(&post, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "找不到文章"})
		return
	}

	// 檢查是否為作者本人
	userID, _ := c.Get("userID")
	if post.AuthorID != userID.(uint) {
		c.JSON(http.StatusForbidden, gin.H{"error": "你沒有權限刪除此文章"})
		return
	}

	if err := db.Delete(&post).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "刪除文章失敗"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "文章刪除成功"})
}

func handleGetAuthor(c *gin.Context) {
	id := c.Param("id")
	var author User
	if err := db.First(&author, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "找不到作者"})
		return
	}

	c.JSON(http.StatusOK, author)
}

func handleUploadFile(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "無法讀取檔案: " + err.Error()})
		return
	}

	// 檢查檔案大小 (限制 10MB)
	const maxSize = 10 << 20 // 10MB
	if file.Size > maxSize {
		c.JSON(http.StatusBadRequest, gin.H{"error": "檔案大小不能超過 10MB"})
		return
	}

	// 檢查檔案類型
	allowedTypes := map[string]bool{
		"image/jpeg": true,
		"image/jpg":  true,
		"image/png":  true,
		"image/gif":  true,
		"image/webp": true,
	}

	// 開啟檔案來檢查 MIME 類型
	src, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "無法開啟檔案"})
		return
	}
	defer src.Close()

	buffer := make([]byte, 512)
	_, err = src.Read(buffer)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "無法讀取檔案"})
		return
	}

	contentType := http.DetectContentType(buffer)
	if !allowedTypes[contentType] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "不支援的檔案類型，僅支援 JPEG, PNG, GIF, WebP"})
		return
	}

	// 創建 uploads 目錄（如果不存在）
	uploadsDir := "./uploads"
	if _, err := os.Stat(uploadsDir); os.IsNotExist(err) {
		err = os.Mkdir(uploadsDir, 0755)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "無法創建上傳目錄"})
			return
		}
	}

	// 生成唯一檔案名
	ext := strings.ToLower(filepath.Ext(file.Filename))
	if ext == "" {
		// 根據 MIME 類型設定副檔名
		switch contentType {
		case "image/jpeg":
			ext = ".jpg"
		case "image/png":
			ext = ".png"
		case "image/gif":
			ext = ".gif"
		case "image/webp":
			ext = ".webp"
		}
	}

	filename := fmt.Sprintf("%d_%s%s", time.Now().Unix(), generateRandomString(8), ext)
	filepath := fmt.Sprintf("%s/%s", uploadsDir, filename)

	// 儲存檔案
	if err := c.SaveUploadedFile(file, filepath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "儲存檔案失敗: " + err.Error()})
		return
	}

	// 回傳檔案 URL
	fileURL := fmt.Sprintf("/uploads/%s", filename)
	c.JSON(http.StatusOK, gin.H{
		"message": "檔案上傳成功",
		"url":     fileURL,
		"filename": filename,
	})
}

// --- 中介層 (Middleware) ---

func authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "缺少 Authorization header"})
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header 格式錯誤"})
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("非預期的簽章方法: %v", token.Header["alg"])
			}
			return jwtSecret, nil
		})

		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "無效的 token: " + err.Error()})
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			// 從 claims 中取得 user ID
			userIDFloat, ok := claims["sub"].(float64)
			if !ok {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "token 中找不到使用者 ID"})
				return
			}
			// 將 userID 存入 context 中，方便後續 handler 使用
			c.Set("userID", uint(userIDFloat))
			c.Next()
		} else {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "無效的 token"})
		}
	}
}

// --- Helper 函式 ---

// generateJWT 產生 JWT
func generateJWT(user *User) (string, error) {
	// 設定 token 的 claims
	claims := jwt.MapClaims{
		"sub": user.ID,                               // Subject (使用者ID)
		"usn": user.Username,                         // Username
		"exp": time.Now().Add(time.Hour * 24).Unix(), // 有效期限 24 小時
		"iat": time.Now().Unix(),                     // 簽發時間
	}

	// 建立 token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 使用 secret 簽署 token 並取得完整的 token 字串
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// generateRandomString 產生指定長度的隨機字串
func generateRandomString(length int) string {
	bytes := make([]byte, length/2)
	if _, err := rand.Read(bytes); err != nil {
		return fmt.Sprintf("%d", time.Now().UnixNano())
	}
	return hex.EncodeToString(bytes)
}
