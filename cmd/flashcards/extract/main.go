package extract

import "gorm.io/gorm"

func Run(db *gorm.DB, deck *string) {
	extractFlashcards(db, deck)
}
