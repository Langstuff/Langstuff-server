package extract

import "gorm.io/gorm"

type ExportTargetType int

const (
	SAVE_IN_DATABASE ExportTargetType = iota
	EXPORT_CSV
	EXPORT_QA
)

func Run(db *gorm.DB, exportTarget ExportTargetType, deck *string) {
	extractFlashcards(db, exportTarget, deck)
}
