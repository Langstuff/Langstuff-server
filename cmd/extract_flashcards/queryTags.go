package main

import (
	"fmt"
	"raiden_fumo/lang_notebook_core/database"

	"gorm.io/gorm"
)

func QueryTags(db *gorm.DB) {
	var tags []database.Tag
	db.Find(&tags)
	for _, tag := range(tags) {
		fmt.Println(tag.Name)
	}
}
