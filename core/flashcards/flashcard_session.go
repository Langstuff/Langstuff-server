package flashcards

type Action struct {
	CardId int
	Grade  int
}

type Session struct {
	ID           string
	CurrentCard  int
	FlashcardIDs []int
	Queue        []int
	ActionList   []Action
}

func (s *Session) ProcessAnswer(grade int) {
	// s.
}
