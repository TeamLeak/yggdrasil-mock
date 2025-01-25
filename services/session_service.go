package services

import (
	"database/sql"
	"net/http"
	"strconv"

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

	token, err := GetToken(db, req.AccessToken)
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

	var character struct {
		ID     int    `json:"id"`
		UUID   string `json:"uuid"`
		Name   string `json:"name"`
		Model  string `json:"model"`
		UserID int    `json:"user_id"`
	}

	err := db.QueryRow(
		`SELECT c.id, c.uuid, c.name, c.model, c.user_id
		 FROM characters c
		 JOIN tokens t ON c.id = t.character_id
		 WHERE c.name = ? AND t.client_token = ?`,
		username, serverID,
	).Scan(&character.ID, &character.UUID, &character.Name, &character.Model, &character.UserID)
	if err != nil {
		c.Status(http.StatusNoContent)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":   character.UUID,
		"name": character.Name,
	})
}
