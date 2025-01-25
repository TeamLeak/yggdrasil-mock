package main

import (
	"flag"
	"fmt"
	"log"
	_ "os"
	"yggdrasil/config"
	"yggdrasil/database"
	"yggdrasil/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	generateConfig := flag.Bool("generate-config", false, "Generate a default config.yaml file")
	flag.Parse()

	if *generateConfig {
		if err := config.GenerateDefaultConfig(); err != nil {
			log.Fatalf("Failed to generate default config: %v", err)
		}
		fmt.Println("Default config.yaml generated successfully.")
		return
	}

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	db, err := database.Connect(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}
	defer func() {
		if cerr := db.Close(); cerr != nil {
			log.Printf("Failed to close the database connection: %v", cerr)
		}
	}()

	if err := database.Migrate(db, cfg); err != nil {
		log.Fatalf("Failed to apply migrations: %v", err)
	}

	r := gin.Default()
	routes.RegisterRoutes(r, db)

	serverAddress := ":" + cfg.App.Port
	if err := r.Run(serverAddress); err != nil {
		log.Fatalf("Failed to start the server: %v", err)
	}
}
