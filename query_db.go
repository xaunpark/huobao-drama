package main

import (
	"fmt"
	"github.com/drama-generator/backend/domain/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	db, err := gorm.Open(sqlite.Open("huobao.db"), &gorm.Config{})
	if err != nil {
		fmt.Println("failed to connect database", err)
		return
	}
	var storyboards []models.Storyboard
	db.Where("episode_id = ?", 91).Find(&storyboards)
	for _, sb := range storyboards {
        imgPrompt := "NULL"
        if sb.ImagePrompt != nil { imgPrompt = *sb.ImagePrompt }
		fmt.Printf("Storyboard %d: ImagePrompt: '%s'\n", sb.ID, imgPrompt)
	}
}
