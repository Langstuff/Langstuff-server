package extract

import (
	"fmt"
	"io"
	"os"
	"raiden_fumo/lang_notebook_core/core/database"
	"raiden_fumo/lang_notebook_core/markdown_parsing"

	"gorm.io/gorm"
)

func saveInDatabase(db *gorm.DB, deckID uint, pair markdown_parsing.Pair) {
	var flashcard database.Flashcard
	db.Where(database.Flashcard{DeckID: deckID, SideA: pair.First, SideB: pair.Second}).
		Attrs(database.Flashcard{DeckID: deckID, SideA: pair.First, SideB: pair.Second}).
		FirstOrCreate(&flashcard)
	for _, tag := range *pair.Tags {
		dbTagEntry := database.Tag{}
		db.Where(database.Tag{Name: tag}).FirstOrCreate(&dbTagEntry)
		dbFlashcardTagPair := database.FlashcardTagPair{}
		db.Where(database.FlashcardTagPair{FlashcardID: flashcard.ID, TagID: dbTagEntry.ID}).
			FirstOrCreate(&dbFlashcardTagPair)
	}
}

func exportInCSV(db *gorm.DB, deckID uint, pair markdown_parsing.Pair) {
	fmt.Printf("\"%s\",\"%s\"\n", pair.Second, pair.First)
}

func exportInQA(db *gorm.DB, deckID uint, pair markdown_parsing.Pair) {
	fmt.Printf("Q: %s\nA: %s\n\n", pair.Second, pair.First)
}

func extractFlashcards(db *gorm.DB, exportTarget ExportTargetType, deckName *string) {
	bytes, err := io.ReadAll(os.Stdin)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	var deck database.Deck
	if exportTarget == SAVE_IN_DATABASE {
		db.Where(database.Deck{Name: *deckName}).
			FirstOrCreate(&deck, database.Deck{Name: *deckName})
	}

	pairs := markdown_parsing.ExtractLearnPairs(bytes)
	for _, pair := range pairs {
		switch exportTarget {
		case EXPORT_CSV:
			exportInCSV(db, deck.ID, pair)
		case SAVE_IN_DATABASE:
			saveInDatabase(db, deck.ID, pair)
		case EXPORT_QA:
			exportInQA(db, deck.ID, pair)
		default:
			panic("Unknown export target " + fmt.Sprint(exportTarget))
		}
	}
}
