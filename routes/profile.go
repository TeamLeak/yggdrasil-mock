package routes

import (
	"database/sql"
	"yggdrasil/services"

	"github.com/gin-gonic/gin"
)

func RegisterProfileRoutes(r *gin.Engine, db *sql.DB) {
	r.POST("/api/profiles/minecraft", func(c *gin.Context) {
		services.QueryProfilesHandler(c, db)
	})
	r.GET("/sessionserver/session/minecraft/profile/:uuid", func(c *gin.Context) {
		services.ProfileHandler(c, db)
	})
}
