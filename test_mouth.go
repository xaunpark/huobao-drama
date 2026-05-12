package main

import (
	"fmt"
	"strings"
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
	NarratorScript   *string
}

func getMouthConstraint(sb Storyboard) string {
	dialogue := ""
	if sb.Dialogue != nil {
		dialogue = *sb.Dialogue
	}
	narrator := ""
	if sb.NarratorScript != nil {
		narrator = *sb.NarratorScript
	}
	audioMode := ""
	if sb.AudioMode != nil {
		audioMode = *sb.AudioMode
	}

	isVoiceoverMode := audioMode == "narrator_only" || (audioMode == "" && narrator != "" && dialogue == "")

	if narrator != "" && isVoiceoverMode {
		return fmt.Sprintf("Narration (voice-over): %s. The character's mouth is strictly closed, silent expression, purely visual acting, no speaking, voiceover scene. --no talking, speaking, moving lips, open mouth, chatting", narrator)
	} else if dialogue != "" {
		dialogueLower := strings.ToLower(strings.TrimSpace(dialogue))
		isVoiceover := strings.HasPrefix(dialogueLower, "(vo)") ||
			strings.HasPrefix(dialogueLower, "(monologue)") ||
			strings.Contains(dialogueLower, "voiceover") ||
			strings.HasPrefix(dialogueLower, "【旁白") ||
			strings.HasPrefix(dialogueLower, "[旁白") ||
			strings.HasPrefix(dialogueLower, "(narrator") ||
			strings.Contains(dialogueLower, "（旁白）")

		if isVoiceover {
			return fmt.Sprintf("Dialogue (voice-over): %s. The character's mouth is strictly closed, silent expression, purely visual acting, no speaking, voiceover scene. --no talking, speaking, moving lips, open mouth, chatting", dialogue)
		} else {
			return fmt.Sprintf("Dialogue: %s. The character is actively speaking, lip-syncing naturally to the dialog, mouth moving", dialogue)
		}
	} else {
		return "The character's mouth is completely closed, silent scene. --no talking, speaking, moving lips"
	}
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
		mc := getMouthConstraint(s)
		fmt.Printf("Shot %d | MouthConstraint: %s\n", s.StoryboardNumber, mc)
	}
}
