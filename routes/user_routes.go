package routes

import (
	"database/sql"
	"yggdrasil/services"

	"github.com/gin-gonic/gin"
)

func RegisterUserRoutes(r *gin.Engine, db *sql.DB) {
	r.GET("/users", func(c *gin.Context) {
		services.GetUsersHandler(c, db)
	})
	r.POST("/users", func(c *gin.Context) {
		services.CreateUserHandler(c, db)
	})
	r.PUT("/users/:id", func(c *gin.Context) {
		services.UpdateUserHandler(c, db)
	})
	r.DELETE("/users/:id", func(c *gin.Context) {
		services.DeleteUserHandler(c, db)
	})
}
