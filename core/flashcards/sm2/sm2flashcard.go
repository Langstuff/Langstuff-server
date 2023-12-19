package sm2

import (
	"math"
	"raiden_fumo/lang_notebook_core/core/database"
)

type Card struct {
	Interval       uint
	EasinessFactor float64
	Repetitions    uint
}

func (card *Card) Calc(grade uint) Card {
	outCard := Card{}
	if grade >= 3 {
		if card.Repetitions == 0 {
			outCard.Interval = 1
		} else if card.Repetitions == 1 {
			outCard.Interval = 6
		} else {
			outCard.Interval = uint(math.Floor(float64(card.Interval) * card.EasinessFactor))
		}
		outCard.Repetitions = card.Repetitions + 1
	} else {
		outCard.Repetitions = 0
		outCard.Interval = 1
	}
	outCard.EasinessFactor = card.EasinessFactor +
		(0.1 - float64(5-grade)*(0.08+float64(5-grade)*0.02))
	if card.EasinessFactor < 1.3 {
		card.EasinessFactor = 1.3
	}
	return outCard
}

func (card *Card) Save() {

}

func (card *Card) Load(cardDbEntry database.Flashcard) {

}
