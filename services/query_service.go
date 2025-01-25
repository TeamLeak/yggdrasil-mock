package services

import (
	"database/sql"
	"net/http"
	"yggdrasil/database"

	"github.com/gin-gonic/gin"
)

func QueryProfilesHandler(c *gin.Context, db *sql.DB) {
	var names []string
	if err := c.ShouldBindJSON(&names); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	profiles, err := database.GetCharactersByNames(db, names)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve profiles"})
		return
	}

	c.JSON(http.StatusOK, profiles)
}
