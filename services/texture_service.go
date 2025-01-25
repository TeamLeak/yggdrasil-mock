package services

import (
	"database/sql"
	"io"
	"net/http"
	_ "time"
	"yggdrasil/database"

	"github.com/gin-gonic/gin"
)

func UploadTextureHandler(c *gin.Context, db *sql.DB) {
	uuid := c.Param("uuid")
	textureType := c.Param("textureType")

	file, _, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file"})
		return
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read file"})
		return
	}

	if err := database.InsertOrUpdateTexture(db, uuid, textureType, data); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload texture"})
		return
	}

	c.Status(http.StatusNoContent)
}

func DeleteTextureHandler(c *gin.Context, db *sql.DB) {
	uuid := c.Param("uuid")
	textureType := c.Param("textureType")

	if err := database.DeleteTexture(db, uuid, textureType); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete texture"})
		return
	}

	c.Status(http.StatusNoContent)
}

func TextureHandler(c *gin.Context, db *sql.DB) {
	hash := c.Param("hash")

	texture, err := database.GetTextureByHash(db, hash)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Texture not found"})
		return
	}

	c.Data(http.StatusOK, "image/png", texture.Data)
}
