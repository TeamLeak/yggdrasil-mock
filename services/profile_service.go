package services

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ProfileHandler(c *gin.Context, db *sql.DB) {
	uuid := c.Param("uuid")
	unsigned := c.Query("unsigned") == "false"

	var character struct {
		ID     int    `json:"id"`
		UUID   string `json:"uuid"`
		Name   string `json:"name"`
		Model  string `json:"model"`
		UserID int    `json:"user_id"`
	}

	err := db.QueryRow(
		`SELECT id, uuid, name, model, user_id FROM characters WHERE uuid = ?`,
		uuid,
	).Scan(&character.ID, &character.UUID, &character.Name, &character.Model, &character.UserID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Profile not found"})
		return
	}

	response := gin.H{
		"id":   character.UUID,
		"name": character.Name,
	}
	if !unsigned {
		response["model"] = character.Model
	}

	c.JSON(http.StatusOK, response)
}
