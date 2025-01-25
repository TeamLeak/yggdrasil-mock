package routes

import (
	"database/sql"
	"yggdrasil/services"

	"github.com/gin-gonic/gin"
)

func RegisterTextureRoutes(r *gin.Engine, db *sql.DB) {
	r.GET("/textures/:hash", func(c *gin.Context) {
		services.TextureHandler(c, db)
	})
	r.DELETE("/api/user/profile/:uuid/:textureType", func(c *gin.Context) {
		services.DeleteTextureHandler(c, db)
	})
	r.PUT("/api/user/profile/:uuid/:textureType", func(c *gin.Context) {
		services.UploadTextureHandler(c, db)
	})
}
