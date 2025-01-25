package routes

import (
	"database/sql"
	"yggdrasil/services"

	"github.com/gin-gonic/gin"
)

func RegisterAPIRoutes(r *gin.Engine, db *sql.DB) {
	r.POST("/api/register", func(c *gin.Context) {
		services.RegisterUserHandler(c, db)
	})
	r.POST("/api/setTexture", func(c *gin.Context) {
		services.SetTextureHandler(c, db)
	})
	r.POST("/api/addCharacter", func(c *gin.Context) {
		services.AddCharacterHandler(c, db)
	})
}
