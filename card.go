package poker

import "fmt"

// Card represents a unique card from a standard 52-card deck
type Card uint32

// CardRank is basically an enum 0-12 with Deuce being 0, Ace being 12
type CardRank uint32

// All possible card ranks in a standard deck
const (
	Deuce CardRank = iota
	Three
	Four
	Five
	Six
	Seven
	Eight
	Nine
	Ten
	Jack
	Queen
	King
	Ace
)

var charToCardRank = map[byte]CardRank{
	'2': Deuce,
	'3': Three,
	'4': Four,
	'5': Five,
	'6': Six,
	'7': Seven,
	'8': Eight,
	'9': Nine,
	'T': Ten,
	'J': Jack,
	'Q': Queen,
	'K': King,
	'A': Ace,
}

func (r CardRank) String() string {
	switch r {
	case Deuce:
		return "2"
	case Three:
		return "3"
	case Four:
		return "4"
	case Five:
		return "5"
	case Six:
		return "6"
	case Seven:
		return "7"
	case Eight:
		return "8"
	case Nine:
		return "9"
	case Ten:
		return "T"
	case Jack:
		return "J"
	case Queen:
		return "Q"
	case King:
		return "K"
	case Ace:
		return "A"
	}

	return ""
}

// CardSuit represents spades, hearts, diamonds, or clubs
type CardSuit uint32

// The four possible card suits
const (
	Spades CardSuit = 1 << iota
	Hearts
	Diamonds
	Clubs
)

var charToCardSuit = map[byte]CardSuit{
	's': Spades,
	'h': Hearts,
	'd': Diamonds,
	'c': Clubs,
}

func (s CardSuit) String() string {
	switch s {
	case Spades:
		return "s"
	case Hearts:
		return "h"
	case Diamonds:
		return "d"
	case Clubs:
		return "c"
	}

	return ""
}

// NewCard takes a rank and suit and returns a card.  Invalid ranks or suits
// will give result in an undefined Card value, so always use the CardRank and
// CardSuit constants.
func NewCard(r CardRank, s CardSuit) Card {
	var rankPrime = primes[r]
	var bitRank uint32 = 1 << r << 16
	var suit = uint32(s) << 12
	var rank = uint32(r) << 8

	return Card(bitRank | suit | rank | rankPrime)
}

// NewCardString takes a two-character string and returns a card.  Unlike
// NewCard, there is also the possibility of an error being returned because of
// how many ways a string could end up *not* representing anything meaningful.
func NewCardString(s string) (Card, error) {
	if len(s) != 2 {
		return 0, fmt.Errorf("NewCardString(%q): need a two-rune string", s)
	}

	var rank = charToCardRank[s[0]]
	if rank < Deuce || rank > Ace {
		return 0, fmt.Errorf("NewCardString(%q): invalid rank", s)
	}

	var suit = charToCardSuit[s[1]]
	if suit != Spades && suit != Hearts && suit != Diamonds && suit != Clubs {
		return 0, fmt.Errorf("NewCardString(%q): invalid suit", s)
	}

	return NewCard(rank, suit), nil
}

// MarshalJSON implements json.Marshaler to convert a single card into its JSON
// string representation
func (c *Card) MarshalJSON() ([]byte, error) {
	return []byte("\"" + c.String() + "\""), nil
}

// UnmarshalJSON implements json.Unmarshaler to take a JSON string and turn it
// into a card
func (c *Card) UnmarshalJSON(b []byte) (err error) {
	*c, err = NewCardString(string(b[1:3]))
	return err
}

func (c Card) String() string {
	return c.Rank().String() + c.Suit().String()
}

// Rank returns the CardRank value for this card
func (c Card) Rank() CardRank {
	return CardRank((uint32(c) >> 8) & 0xF)
}

// Suit returns the CardSuit value for this card
func (c Card) Suit() CardSuit {
	return CardSuit((uint32(c) >> 12) & 0xF)
}
