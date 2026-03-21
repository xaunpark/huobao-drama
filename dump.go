package main

import (
	"fmt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Drama struct {
	ID               uint
	Title            string
	Style            string
	PromptTemplateID *uint
}

func main() {
	db, err := gorm.Open(sqlite.Open("data/drama_generator.db"), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	var dramas []Drama
	db.Find(&dramas)
	for _, d := range dramas {
		var tplID uint
		if d.PromptTemplateID != nil {
			tplID = *d.PromptTemplateID
		}
		fmt.Printf("Drama ID: %d, Style: %s, TemplateID: %d\n", d.ID, d.Style, tplID)
	}
}
