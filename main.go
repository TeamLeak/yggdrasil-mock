package main

import (
	"log"
	"yggdrasil/database"
	"yggdrasil/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	db, err := database.Connect()
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}
	defer func() {
		if cerr := db.Close(); cerr != nil {
			log.Printf("Failed to close the database connection: %v", cerr)
		}
	}()

	if err := database.Migrate(db); err != nil {
		log.Fatalf("Failed to apply migrations: %v", err)
	}

	r := gin.Default()
	routes.RegisterRoutes(r, db)

	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to start the server: %v", err)
	}
}
