package main

import (
	"log"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func main() {
	db, err := gorm.Open(sqlite.Open("data/drama_generator.db"), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	result := db.Exec("UPDATE storyboards SET script_segment = lyrics_anchor WHERE narrative_part IS NOT NULL AND lyrics_anchor IS NOT NULL AND lyrics_anchor != ''")
	if result.Error != nil {
		log.Fatal(result.Error)
	}
	
	log.Printf("Successfully updated %d shots to display lyrics anchor in the Narrator field.", result.RowsAffected)
}
