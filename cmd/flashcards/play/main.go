package play

import (
	"bufio"
	"fmt"
	"os"
	"raiden_fumo/lang_notebook_core/core/database"
	"strconv"

	"gorm.io/gorm"
)

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

	new_values := flashcard.Calc(grade)

	flashcardSM2Info.Repetition = new_values.Repetitions
	flashcardSM2Info.Interval = new_values.Interval
	flashcardSM2Info.EasinessFactor = new_values.EasinessFactor
	db.Save(&flashcardSM2Info)
}

func Run(db *gorm.DB, deck *string, tags *[]string) {
	var flashcards = database.GetFlashcardsByTagList(db, *deck, *tags)
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
