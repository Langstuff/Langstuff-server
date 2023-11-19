package main

import (
	"fmt"
	"raiden_fumo/lang_notebook_core/database"

	"gorm.io/gorm"
)

func QueryFlashcards(db *gorm.DB, tags *[]string) {
	var tagEntries []database.Tag
	db.Where("name IN ?", *tags).Find(&tagEntries)
	fmt.Println(tagEntries)
	for _, entry := range tagEntries {
		flashcards := entry.GetFlashcards(db)
		for _, fc := range flashcards {
			fmt.Println(fc.SideA, fc.SideB)
		}
	}
}
