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
		// Find the latest completed image generation EXACTLY to know what mode the user generated!
		var latestImage models.ImageGeneration
		err := db.Where("storyboard_id = ? AND status = 'completed'", sb.ID).
			Order("created_at desc").
			First(&latestImage).Error

		if err != nil {
			// If no image, fallback to "key" or leave it. 
			continue
		}

		if latestImage.FrameType != nil && *latestImage.FrameType != "" {
			inferredMode := *latestImage.FrameType
			if sb.GenerationMode == nil || *sb.GenerationMode != inferredMode {
				fmt.Printf("Shot %d: Changing from %v to %s (recovered from actual generated Image frame_type)\n", sb.StoryboardNumber, sb.GenerationMode, inferredMode)
				modePtr := &inferredMode
				if err := db.Model(&sb).UpdateColumn("generation_mode", modePtr).Error; err != nil {
					log.Printf("Failed to update shot %d: %v", sb.StoryboardNumber, err)
				} else {
					updated++
				}
			}
		}
	}
	
	fmt.Printf("Finished! Recovered %d shots based on their ACTUAL image generations.\n", updated)
}
