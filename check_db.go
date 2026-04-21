package main

import (
	"fmt"
	"log"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func main() {
	db, err := gorm.Open(sqlite.Open("data/drama_generator.db"), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	var results []map[string]interface{}
	// query the columns of storyboards table
	db.Raw("PRAGMA table_info(storyboards);").Scan(&results)
	
	narrativePartExists := false
	for _, col := range results {
        name := col["name"].(string)
        if name == "narrative_part" {
            narrativePartExists = true
        }
	}
    
    if narrativePartExists {
        fmt.Println("YAY! narrative_part exists!")
        var count int64
        db.Table("storyboards").Where("narrative_part IS NOT NULL").Count(&count)
        fmt.Printf("There are %d shots with narrative_part populated.\n", count)
    } else {
        fmt.Println("OH NO! narrative_part does NOT exist!")
    }
}
