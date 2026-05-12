package main

import (
	"fmt"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"log"
)

type Storyboard struct {
	ID               uint
	StoryboardNumber int
	Dialogue         *string
	AudioMode        *string
	Action           *string
	VideoPrompt      *string
	VideoPromptDistilled *string
}

func main() {
	db, err := gorm.Open(sqlite.Open("data/drama_generator.db"), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	var shots []Storyboard
	if err := db.Order("id desc").Limit(5).Find(&shots).Error; err != nil {
		log.Fatal(err)
	}

	for _, s := range shots {
		dlg := "NULL"
		if s.Dialogue != nil {
			dlg = *s.Dialogue
		}
		vp := "NULL"
		if s.VideoPromptDistilled != nil {
			vp = *s.VideoPromptDistilled
		}
		fmt.Printf("Shot %d | Dialogue: %s | Distilled: %s\n", s.StoryboardNumber, dlg, vp)
	}
}
