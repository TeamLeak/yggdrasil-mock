package services

import (
	"database/sql"
	"net/http"
	"time"
	"yggdrasil/database"
	"yggdrasil/models"
	"yggdrasil/utils"

	"github.com/gin-gonic/gin"
)

func SignoutHandler(c *gin.Context, db *sql.DB) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	user, err := database.GetUserByEmailOrCharacter(db, req.Username)
	if err != nil || user.Password != req.Password {
		c.JSON(http.StatusForbidden, gin.H{"error": "Invalid credentials"})
		return
	}

	if err = database.RevokeAllTokens(db, user.ID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to revoke tokens"})
		return
	}

	c.Status(http.StatusNoContent)
}

func InvalidateHandler(c *gin.Context, db *sql.DB) {
	var req struct {
		AccessToken string `json:"accessToken"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	if err := database.DeleteToken(db, req.AccessToken); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to invalidate token"})
		return
	}

	c.Status(http.StatusNoContent)
}

func ValidateHandler(c *gin.Context, db *sql.DB) {
	var req struct {
		AccessToken string `json:"accessToken"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	_, err := database.GetToken(db, req.AccessToken)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "Invalid token"})
		return
	}

	c.Status(http.StatusNoContent)
}

func RefreshHandler(c *gin.Context, db *sql.DB) {
	var req struct {
		AccessToken string `json:"accessToken"`
		ClientToken string `json:"clientToken"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	oldToken, err := database.GetToken(db, req.AccessToken)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "Invalid token"})
		return
	}

	newToken := models.Token{
		AccessToken: utils.GenerateUUID(),
		ClientToken: req.ClientToken,
		CreatedAt:   time.Now(),
		UserID:      oldToken.UserID,
		CharacterID: oldToken.CharacterID,
	}

	if err = database.InsertToken(db, &newToken); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to refresh token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"accessToken": newToken.AccessToken,
		"clientToken": newToken.ClientToken,
	})
}

func AuthenticateHandler(c *gin.Context, db *sql.DB) {
	var req struct {
		Username    string `json:"username"`
		Password    string `json:"password"`
		ClientToken string `json:"clientToken"`
		RequestUser bool   `json:"requestUser"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	var user *models.User
	var character *models.Character

	loginWithCharacterName := true // Конфигурация
	if loginWithCharacterName {
		character, _ = database.FindCharacterByName(db, req.Username)
	}
	if character == nil {
		var err error
		user, err = database.GetUserByEmailOrCharacter(db, req.Username)
		if err != nil || user.Password != req.Password {
			c.JSON(http.StatusForbidden, gin.H{"error": "Invalid credentials"})
			return
		}
	} else {
		user, _ = database.GetUserByID(db, character.UserID)
		if user.Password != req.Password {
			c.JSON(http.StatusForbidden, gin.H{"error": "Invalid credentials"})
			return
		}
	}

	if req.ClientToken == "" {
		req.ClientToken = utils.GenerateUUID()
	}

	token := models.Token{
		AccessToken: utils.GenerateUUID(),
		ClientToken: req.ClientToken,
		CreatedAt:   time.Now(),
		UserID:      user.ID,
	}
	if character != nil {
		token.CharacterID = character.ID
	}

	if err := database.InsertToken(db, &token); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to create token"})
		return
	}

	response := gin.H{
		"accessToken":       token.AccessToken,
		"clientToken":       token.ClientToken,
		"availableProfiles": database.GetUserCharacters(db, user.ID),
	}

	if character != nil {
		response["selectedProfile"] = map[string]string{
			"id":   character.UUID,
			"name": character.Name,
		}
	}

	if req.RequestUser {
		response["user"] = map[string]interface{}{
			"id":    user.ID,
			"email": user.Email,
		}
	}

	c.JSON(http.StatusOK, response)
}
