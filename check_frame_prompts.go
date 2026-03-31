package main

import (
	"fmt"
	"log"

	"github.com/drama-generator/backend/infrastructure/database"
	"github.com/drama-generator/backend/domain/models"
	"github.com/drama-generator/backend/pkg/config"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	db, err := database.NewDatabase(cfg.Database)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	var count int64
	db.Model(&models.FramePrompt{}).Count(&count)
	fmt.Printf("Total rows in frame_prompts table: %d\n", count)
	
	// Print a few rows
	var prompts []models.FramePrompt
	db.Limit(5).Find(&prompts)
	for _, p := range prompts {
		fmt.Printf("Prompt ID: %d, StoryboardID: %d, FrameType: %s\n", p.ID, p.StoryboardID, p.FrameType)
	}
}
