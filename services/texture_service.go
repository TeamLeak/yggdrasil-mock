package services

import (
	"database/sql"
	"io"
	"net/http"
	"time"

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

	hash := uuid + "-" + textureType

	_, err = db.Exec(
		`INSERT INTO textures (hash, data, uploaded_at) VALUES (?, ?, ?)
		ON CONFLICT(hash) DO UPDATE SET data = excluded.data, uploaded_at = excluded.uploaded_at`,
		hash, data, time.Now(),
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload texture"})
		return
	}

	c.Status(http.StatusNoContent)
}

func DeleteTextureHandler(c *gin.Context, db *sql.DB) {
	uuid := c.Param("uuid")
	textureType := c.Param("textureType")

	hash := uuid + "-" + textureType

	_, err := db.Exec(`DELETE FROM textures WHERE hash = ?`, hash)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete texture"})
		return
	}

	c.Status(http.StatusNoContent)
}

func TextureHandler(c *gin.Context, db *sql.DB) {
	hash := c.Param("hash")

	var texture struct {
		Hash       string `json:"hash"`
		Data       []byte `json:"data"`
		UploadedAt string `json:"uploaded_at"`
	}

	err := db.QueryRow(
		`SELECT hash, data, uploaded_at FROM textures WHERE hash = ?`,
		hash,
	).Scan(&texture.Hash, &texture.Data, &texture.UploadedAt)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Texture not found"})
		return
	}

	c.Data(http.StatusOK, "image/png", texture.Data)
}
