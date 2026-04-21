package main

import (
	"encoding/json"
	"fmt"
	"log"
	"github.com/drama-generator/backend/domain/models"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func main() {
	db, err := gorm.Open(sqlite.Open("data/drama_generator.db"), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	var sb models.Storyboard
	if err := db.Where("narrative_part IS NOT NULL AND narrative_part != ''").Order("id DESC").First(&sb).Error; err != nil {
		log.Fatal("Could not find any storyboard with narrative_part:", err)
	}

	b, _ := json.MarshalIndent(sb, "", "  ")
	fmt.Println("Storyboard JSON representation:")
	fmt.Println(string(b))
}
