package services

import (
	"database/sql"
	"net/http"
	"time"
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

	user, err := GetUserByEmailOrCharacter(db, req.Username)
	if err != nil || user.Password != req.Password {
		c.JSON(http.StatusForbidden, gin.H{"error": "Invalid credentials"})
		return
	}

	err = RevokeAllTokens(db, user.ID)
	if err != nil {
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

	err := DeleteToken(db, req.AccessToken)
	if err != nil {
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

	_, err := GetToken(db, req.AccessToken)
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

	oldToken, err := GetToken(db, req.AccessToken)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "Invalid token"})
		return
	}

	newToken := Token{
		AccessToken: utils.GenerateUUID(),
		ClientToken: req.ClientToken,
		CreatedAt:   time.Now(),
		UserID:      oldToken.UserID,
		CharacterID: oldToken.CharacterID,
	}

	err = InsertToken(db, &newToken)
	if err != nil {
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

	user, err := GetUserByEmailOrCharacter(db, req.Username)
	if err != nil || user.Password != req.Password {
		c.JSON(http.StatusForbidden, gin.H{"error": "Invalid credentials"})
		return
	}

	if req.ClientToken == "" {
		req.ClientToken = utils.GenerateUUID()
	}

	token := Token{
		AccessToken: utils.GenerateUUID(),
		ClientToken: req.ClientToken,
		CreatedAt:   time.Now(),
		UserID:      user.ID,
	}

	err = InsertToken(db, &token)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to create token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"accessToken": token.AccessToken,
		"clientToken": token.ClientToken,
	})
}
