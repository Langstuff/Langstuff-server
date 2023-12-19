package flashcards

type FlashcardGrade interface{}

type Flashcard interface {
	CalculateStep(grade FlashcardGrade) *Flashcard
	Save()
}
