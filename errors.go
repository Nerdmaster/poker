package poker

// PokerError is a simple string wrapper to make constants into usable and
// meaningful returned error values
type PokerError string

// Hand errors
const (
	ErrEmptyDeck        PokerError = "cannot draw from empty deck"
	ErrInvalidCardCount PokerError = "invalid card count"
)

func (e PokerError) Error() string {
	return string(e)
}
