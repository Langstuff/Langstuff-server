package main

import (
	"fmt"
	"os"
	"raiden_fumo/lang_notebook_core/database"

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
	parseCmd := parser.NewCommand("parse", "Parse markdown file and extract flashcards")
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
	} else if parseCmd.Happened() {
		extractFlashcards(db)
	}
}
