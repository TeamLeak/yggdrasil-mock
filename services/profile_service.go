package services

import (
	"database/sql"
	"net/http"
	"yggdrasil/database"

	"github.com/gin-gonic/gin"
)

func ProfileHandler(c *gin.Context, db *sql.DB) {
	uuid := c.Param("uuid")
	unsigned := c.Query("unsigned") == "false"

	character, err := database.GetCharacterByUUID(db, uuid)
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
