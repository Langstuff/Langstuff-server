package main

import (
	"fmt"
	"os"
	"raiden_fumo/lang_notebook_core/cmd/flashcards/extract"
	"raiden_fumo/lang_notebook_core/cmd/flashcards/list"
	"raiden_fumo/lang_notebook_core/cmd/flashcards/play"
	"raiden_fumo/lang_notebook_core/database"

	"github.com/akamensky/argparse"
)

func main() {
	parser := argparse.NewParser("langstuff", "langstuff")
	flashcardCmd := parser.NewCommand("flashcards", "flashcards")

	extractFlashcardsCmd := flashcardCmd.NewCommand("extract", "extract")
	extractCmdDeck := extractFlashcardsCmd.String("d", "deck", &argparse.Options{Help: "deck"})
	extractFlashcardsTarget := extractFlashcardsCmd.Selector("t", "target", []string{
		"save",
		"qa",
		"csv",
	}, &argparse.Options{})

	listFlashcardsCmd := flashcardCmd.NewCommand("list", "list")

	playFlashcardsCmd := flashcardCmd.NewCommand("play", "play")
	flashcardCmdDeck := playFlashcardsCmd.String("d", "deck", &argparse.Options{Help: "deck"})
	flashcardCmdTags := playFlashcardsCmd.StringList("t", "tags", &argparse.Options{Help: "tags"})

	err := parser.Parse(os.Args)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	db := database.InitializeDatabase()

	if extractFlashcardsCmd.Happened() {
		var exportTarget extract.ExportTargetType = 0
		switch *extractFlashcardsTarget {
		case "save":
			exportTarget = extract.SAVE_IN_DATABASE
		case "qa":
			exportTarget = extract.EXPORT_QA
		case "csv":
			exportTarget = extract.EXPORT_CSV
		default:
			panic("unknown export target " + *extractFlashcardsTarget)
		}
		extract.Run(db, exportTarget, extractCmdDeck)
	} else if playFlashcardsCmd.Happened() {
		play.Run(db, flashcardCmdDeck, flashcardCmdTags)
	} else if listFlashcardsCmd.Happened() {
		list.Run(db)
	}
}
