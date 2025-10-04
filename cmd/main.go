package main

import (
	"log"

	"testcase/cmd/http"
	"testcase/config"
	"testcase/internal/infrastructures/database"
)

func main() {
	cfg := config.LoadConfig()

	db, err := database.NewDatabase(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Printf("Error closing database: %v", err)
		}
	}()

	entityRegistry := database.NewEntityRegistry()
	if err := entityRegistry.RegisterAndMigrate(db); err != nil {
		log.Fatalf("Failed to run database migrations: %v", err)
	}

	if cfg.HttpServer.Env != "production" {
		entityRegistry.ListRegisteredEntities()
	}

	server := http.NewServer(cfg, db)

	log.Println("🎯 Starting POS application...")
	if err := server.Start(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
