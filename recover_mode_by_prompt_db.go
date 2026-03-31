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

	// For Project 53, Episode 1
	var episode models.Episode
	if err := db.Where("id = ? AND drama_id = ?", 1, 53).First(&episode).Error; err != nil {
		fmt.Printf("Warning: Could not strictly verify drama_id 53, falling back to episode_id 1... Error: %v\n", err)
	}

	var storyboards []models.Storyboard
	if err := db.Where("episode_id = ?", 1).Find(&storyboards).Error; err != nil {
		log.Fatalf("Failed to load storyboards: %v", err)
	}

	fmt.Printf("Found %d storyboards for episode 1 (Drama 53)\n", len(storyboards))
	
	updated := 0
	
	for _, sb := range storyboards {
		// Just check which gen mode the shot's image prompt is in
		// The most reliable way is to check the `frame_prompts` table for the latest generated prompt for this storyboard!
		var latestPrompt models.FramePrompt
		err := db.Where("storyboard_id = ?", sb.ID).
			Order("created_at desc").
			First(&latestPrompt).Error

		if err != nil {
			// No frame prompt explicitly generated for this shot yet
			continue
		}

		// The FrameType of the prompt is exactly the Gen Mode it was generated under!
		inferredMode := latestPrompt.FrameType
		
		// Note: The UI expects 'key' or 'action'
		if inferredMode != "key" && inferredMode != "action" {
			// If it's something else like 'first', map it to 'key' because they are both structurally single frames
			inferredMode = "key"
		}

		if sb.GenerationMode == nil || *sb.GenerationMode != inferredMode {
			fmt.Printf("Shot %d: Changing from %v to %s (recovered strictly from frame_prompts FrameType)\n", sb.StoryboardNumber, sb.GenerationMode, inferredMode)
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
	
	fmt.Printf("Finished! Recovered %d shots based on their ACTUAL image frame_prompts.\n", updated)
}
