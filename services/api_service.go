package services

import (
	"database/sql"
	"net/http"
	"yggdrasil/database"
	"yggdrasil/models"
	"yggdrasil/utils"

	"github.com/gin-gonic/gin"
)

func RegisterUserHandler(c *gin.Context, db *sql.DB) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&req); err != nil || req.Email == "" || req.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid request"})
		return
	}

	// Хэшируем пароль
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to hash password"})
		return
	}

	// Создаем пользователя
	user := models.User{
		Email:    req.Email,
		Password: hashedPassword,
	}

	if err := database.InsertUser(db, &user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to register user"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"status": "success", "message": "User registered successfully"})
}

func SetTextureHandler(c *gin.Context, db *sql.DB) {
	var req struct {
		UUID        string `json:"uuid"`
		TextureHash string `json:"texture_hash"`
	}
	if err := c.ShouldBindJSON(&req); err != nil || req.UUID == "" || req.TextureHash == "" {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid request"})
		return
	}

	// Устанавливаем текстуру персонажу
	if err := database.SetCharacterTexture(db, req.UUID, req.TextureHash); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to set texture"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Texture set successfully"})
}

func AddCharacterHandler(c *gin.Context, db *sql.DB) {
	var req struct {
		UUID   string `json:"uuid"`
		Name   string `json:"name"`
		Model  string `json:"model"`
		UserID int    `json:"user_id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil || req.UUID == "" || req.Name == "" || req.UserID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid request"})
		return
	}

	// Создаем персонажа
	character := models.Character{
		UUID:   req.UUID,
		Name:   req.Name,
		Model:  req.Model,
		UserID: req.UserID,
	}

	if err := database.InsertCharacter(db, &character); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to add character"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"status": "success", "message": "Character added successfully"})
}
