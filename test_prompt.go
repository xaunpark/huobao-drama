package main

import (
	"fmt"

	"github.com/drama-generator/backend/application/services"
	"github.com/drama-generator/backend/pkg/config"
	"github.com/drama-generator/backend/pkg/logger"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	cfg := &config.Config{
		App: config.AppConfig{Language: "zh"},
	}
	db, err := gorm.Open(sqlite.Open("data/drama_generator.db"), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	log := logger.NewLogger(cfg)
	
	// Create prompt template service
	pts := services.NewPromptTemplateService(db, log)
	
	// Print template prompt
	templateStylePrompt := pts.ResolvePromptIfCustom(19, "style_prompt")
	fmt.Println("--- ResolvePromptIfCustom ---")
	fmt.Println(templateStylePrompt)
	
	sceneExt := pts.ResolvePrompt(19, "scene_extraction")
	fmt.Println("--- ResolvePrompt (scene) ---")
	fmt.Println(sceneExt)

	promptI18n := services.NewPromptI18n(cfg)
	promptI18n.SetTemplateService(pts)

	res := promptI18n.WithDramaSceneExtractionPrompt(19, "ghibli", "")
	fmt.Println("--- WithDramaSceneExtractionPrompt ---")
	fmt.Println(res)
}
