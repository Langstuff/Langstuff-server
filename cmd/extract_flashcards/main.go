package main

import (
	"fmt"
	"io"
	"os"
	"raiden_fumo/lang_notebook_core/database"
	"raiden_fumo/lang_notebook_core/markdown_parsing"

	"github.com/akamensky/argparse"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func initializeDatabase() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		fmt.Println(3)
		panic("failed to connect database")
	}
	db.AutoMigrate(
		&database.Deck{},
		&database.Flashcard{},
		&database.SM2Record{},
		&database.Tag{},
		&database.FlashcardTagPair{},
	)
	return db
}

func main() {
	parser := argparse.NewParser("print", "Prints provided string to stdout")
	parser.NewCommand("parse", "Parse markdown file and extract flashcards")
	getCmd := parser.NewCommand("get", "Query")
	getTagsCmd := getCmd.NewCommand("tags", "Get tags")
	getFlashcardsCmd := getCmd.NewCommand("flashcards", "Get flashcards")
	tags := getCmd.StringList("t", "tag", &argparse.Options{Help: "tag for querying"})

	err := parser.Parse(os.Args)

	if err != nil {
		fmt.Print(parser.Usage(err))
		os.Exit(0)
	}
	db := initializeDatabase()

	if getCmd.Happened() {
		if getTagsCmd.Happened() {
			QueryTags(db)
		} else if getFlashcardsCmd.Happened() {
			QueryFlashcards(db, tags)
		}
	}


	bytes, err := io.ReadAll(os.Stdin)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}


	var deck database.Deck
	db.Where(database.Deck{Name: "Hello"}).FirstOrCreate(&deck, database.Deck{Name: "Hello", Description: "Hello, world"})

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
