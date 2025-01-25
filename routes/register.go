package routes

import (
	"database/sql"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine, db *sql.DB) {
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Yggdrasil Mock Server running"})
	})
	r.GET("/status", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Status handler placeholder"})
	})

	RegisterAuthRoutes(r, db)
	RegisterSessionRoutes(r, db)
	RegisterProfileRoutes(r, db)
	RegisterTextureRoutes(r, db)
	RegisterUserRoutes(r, db)
	RegisterAPIRoutes(r, db)
}
