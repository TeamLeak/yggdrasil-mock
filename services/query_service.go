package services

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
)

func QueryProfilesHandler(c *gin.Context, db *sql.DB) {
	var names []string
	if err := c.ShouldBindJSON(&names); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	var profiles []map[string]string
	for _, name := range names {
		var character struct {
			ID     int    `json:"id"`
			UUID   string `json:"uuid"`
			Name   string `json:"name"`
			Model  string `json:"model"`
			UserID int    `json:"user_id"`
		}

		err := db.QueryRow(
			`SELECT id, uuid, name, model, user_id FROM characters WHERE name = ?`,
			name,
		).Scan(&character.ID, &character.UUID, &character.Name, &character.Model, &character.UserID)
		if err == nil {
			profiles = append(profiles, map[string]string{
				"id":   character.UUID,
				"name": character.Name,
			})
		}
	}

	c.JSON(http.StatusOK, profiles)
}
