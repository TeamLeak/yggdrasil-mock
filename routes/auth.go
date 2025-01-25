package routes

import (
	"database/sql"
	"yggdrasil/services"

	"github.com/gin-gonic/gin"
)

func RegisterAuthRoutes(r *gin.Engine, db *sql.DB) {
	r.POST("/authserver/authenticate", func(c *gin.Context) {
		services.AuthenticateHandler(c, db)
	})
	r.POST("/authserver/refresh", func(c *gin.Context) {
		services.RefreshHandler(c, db)
	})
	r.POST("/authserver/validate", func(c *gin.Context) {
		services.ValidateHandler(c, db)
	})
	r.POST("/authserver/invalidate", func(c *gin.Context) {
		services.InvalidateHandler(c, db)
	})
	r.POST("/authserver/signout", func(c *gin.Context) {
		services.SignoutHandler(c, db)
	})
}
