package extract

import (
	"fmt"
	"io"
	"os"
	"raiden_fumo/lang_notebook_core/database"
	"raiden_fumo/lang_notebook_core/markdown_parsing"

	"gorm.io/gorm"
)

func extractFlashcards(db *gorm.DB, deckName *string) {
	bytes, err := io.ReadAll(os.Stdin)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	var deck database.Deck
	db.Where(database.Deck{Name: "Hello"}).
		FirstOrCreate(&deck, database.Deck{Name: "Hello", Description: "Hello, world"})

	pairs := markdown_parsing.ExtractLearnPairs(bytes)
	for _, pair := range pairs {
		fmt.Printf("Q: %s\nA: %s\n\n", pair.Second, pair.First)
		var flashcard database.Flashcard
		db.Where(database.Flashcard{DeckID: deck.ID, SideA: pair.First, SideB: pair.Second}).
			Attrs(database.Flashcard{DeckID: deck.ID, SideA: pair.First, SideB: pair.Second}).
			FirstOrCreate(&flashcard)
		for _, tag := range *pair.Tags {
			dbTagEntry := database.Tag{}
			db.Where(database.Tag{Name: tag}).FirstOrCreate(&dbTagEntry)
			dbFlashcardTagPair := database.FlashcardTagPair{}
			db.Where(database.FlashcardTagPair{FlashcardID: flashcard.ID, TagID: dbTagEntry.ID}).
				FirstOrCreate(&dbFlashcardTagPair)
		}
	}
}
