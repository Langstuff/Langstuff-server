package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"raiden_fumo/lang_notebook_core/database"
	"strconv"

	"github.com/akamensky/argparse"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Card struct {
	interval       uint
	easinessFactor float64
	repetitions    uint
}

func sm2Calc(grade uint, card Card) Card {
	outCard := Card{}
	if grade >= 3 {
		if card.repetitions == 0 {
			outCard.interval = 1
		} else if card.repetitions == 1 {
			outCard.interval = 6
		} else {
			outCard.interval = uint(math.Floor(float64(card.interval) * card.easinessFactor))
		}
		outCard.repetitions = card.repetitions + 1
	} else {
		outCard.repetitions = 0
		outCard.interval = 1
	}
	outCard.easinessFactor = card.easinessFactor +
		(0.1 - float64(5-grade)*(0.08+float64(5-grade)*0.02))
	if card.easinessFactor < 1.3 {
		card.easinessFactor = 1.3
	}
	return outCard
}

func readNumber(scanner *bufio.Scanner) (uint, bool) {
	for {
		fmt.Print("Grade: ")
		if scanner.Scan() {
			num, err := strconv.Atoi(scanner.Text())
			if err != nil {
				fmt.Println("Invalid input:", scanner.Text())
				continue
			}
			if num < 0 || num > 5 {
				fmt.Println("0 >= num <= 5")
			}
			return uint(num), false
		} else {
			return 0, true
		}
	}
}

func getFlashcardLearningInfo(db *gorm.DB, flashcard *database.Flashcard) *database.SM2Record {
	flashcardSM2Info := &database.SM2Record{}
	db.Where(database.SM2Record{FlashcardID: flashcard.ID}).
		FirstOrCreate(flashcardSM2Info, database.SM2Record{EasinessFactor: 2.5})
	return flashcardSM2Info
}

func smStep(grade uint, db *gorm.DB, flashcard *database.Flashcard) {
	flashcardSM2Info := getFlashcardLearningInfo(db, flashcard)

	new_values := sm2Calc(
		grade,
		Card{
			flashcardSM2Info.Repetition,
			float64(flashcardSM2Info.EasinessFactor),
			flashcardSM2Info.Interval,
		},
	)

	flashcardSM2Info.Repetition = new_values.repetitions
	flashcardSM2Info.Interval = new_values.interval
	flashcardSM2Info.EasinessFactor = new_values.easinessFactor
	db.Save(&flashcardSM2Info)
}

func main() {
	parser := argparse.NewParser("play_flashcards", "Play flashcards")
	deck := parser.String("d", "deck", &argparse.Options{Help: "deck"})
	tags := parser.StringList("t", "tags", &argparse.Options{Help: "tags"})
	err := parser.Parse(os.Args)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		fmt.Println(3)
		panic("failed to connect to database")
	}

	var flashcards = database.GetFlashcardsByTagList(db, deck, tags)
	scanner := bufio.NewScanner(os.Stdin)
	for _, flashcard := range flashcards {
		fmt.Println("-------------------")
		fmt.Print(flashcard.SideA)
		scanner.Scan()
		fmt.Println(flashcard.SideB)
		grade, eof_met := readNumber(scanner)
		if eof_met {
			fmt.Println("EOF, exit")
			os.Exit(0)
		}

		smStep(grade, db, &flashcard)
	}
}
