package main

import (
	"fmt"
	"log"
	"strings"

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
	
	sqlDB, err := db.DB()
	if err == nil {
		defer sqlDB.Close()
	}

	var storyboards []models.Storyboard
	// Query for episode 1
	if err := db.Where("episode_id = ?", 1).Find(&storyboards).Error; err != nil {
		log.Fatalf("Failed to load storyboards: %v", err)
	}

	fmt.Printf("Found %d storyboards for episode 1\n", len(storyboards))
	
	updated := 0
	
	for _, sb := range storyboards {
		inferredMode := "key" // Default to key
		
		if sb.ImagePrompt != nil {
			promptText := strings.ToLower(*sb.ImagePrompt)
			// Check for typical action sequence keywords
			if strings.Contains(promptText, "1x3") || 
				strings.Contains(promptText, "horizontal animation strip") ||
				strings.Contains(promptText, "first panel") ||
				strings.Contains(promptText, "split screen") ||
				strings.Contains(promptText, "panel animation") {
				inferredMode = "action"
			}
		}

		if sb.GenerationMode == nil || *sb.GenerationMode != inferredMode {
			fmt.Printf("Shot %d: Changing from %v to %s (based on latest image generation)\n", sb.StoryboardNumber, sb.GenerationMode, inferredMode)
			modePtr := &inferredMode
			if err := db.Model(&sb).UpdateColumn("generation_mode", modePtr).Error; err != nil {
				log.Printf("Failed to update shot %d: %v", sb.StoryboardNumber, err)
			} else {
				updated++
			}
		} else {
		    fmt.Printf("Shot %d: Already correct (%s)\n", sb.StoryboardNumber, inferredMode)
		}
	}
	
	fmt.Printf("Finished! Recovered %d shots based on their images.\n", updated)
}
