package routes

import (
	"database/sql"
	"yggdrasil/services"

	"github.com/gin-gonic/gin"
)

func RegisterSessionRoutes(r *gin.Engine, db *sql.DB) {
	r.POST("/sessionserver/session/minecraft/join", func(c *gin.Context) {
		services.JoinServerHandler(c, db)
	})
	r.GET("/sessionserver/session/minecraft/hasJoined", func(c *gin.Context) {
		services.HasJoinedHandler(c, db)
	})
}
