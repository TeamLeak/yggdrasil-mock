package services

import (
	"database/sql"
	"net/http"
	"strconv"
	"yggdrasil/database"

	"github.com/gin-gonic/gin"
)

func JoinServerHandler(c *gin.Context, db *sql.DB) {
	var req struct {
		AccessToken     string `json:"accessToken"`
		SelectedProfile string `json:"selectedProfile"`
		ServerID        string `json:"serverId"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	token, err := database.GetToken(db, req.AccessToken)
	if err != nil || token == nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "Invalid token"})
		return
	}

	if strconv.Itoa(token.CharacterID) != req.SelectedProfile {
		c.JSON(http.StatusForbidden, gin.H{"error": "Invalid profile"})
		return
	}

	c.Status(http.StatusNoContent)
}

func HasJoinedHandler(c *gin.Context, db *sql.DB) {
	serverID := c.Query("serverId")
	username := c.Query("username")
	if serverID == "" || username == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing required parameters"})
		return
	}

	character, err := database.FindCharacterByServerAndName(db, username, serverID)
	if err != nil {
		c.Status(http.StatusNoContent)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":   character.UUID,
		"name": character.Name,
	})
}
