package list

import (
	"fmt"
	"raiden_fumo/lang_notebook_core/core/database"

	"gorm.io/gorm"
)

func Run(db *gorm.DB) {
	var tags []database.Tag
	db.Find(&tags)
	for _, tag := range tags {
		fmt.Println(tag.Name)
	}
}
