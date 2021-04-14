package poker

import "fmt"

type Card uint32

type CardRank uint32

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

var CharToCardRank = map[byte]CardRank{
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

type CardSuit uint32

const (
	Spades CardSuit = 1 << iota
	Hearts
	Diamonds
	Clubs
)

var CharToCardSuit = map[byte]CardSuit{
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

func NewCard(r CardRank, s CardSuit) Card {
	var rankPrime = primes[r]
	var bitRank uint32 = 1 << r << 16
	var suit = uint32(s) << 12
	var rank = uint32(r) << 8

	return Card(bitRank | suit | rank | rankPrime)
}

func NewCardString(s string) (Card, error) {
	if len(s) != 2 {
		return 0, fmt.Errorf("NewCardString(%q): need a two-rune string", s)
	}

	var rank = CharToCardRank[s[0]]
	if rank < Deuce || rank > Ace {
		return 0, fmt.Errorf("NewCardString(%q): invalid rank", s)
	}

	var suit = CharToCardSuit[s[1]]
	if suit != Spades && suit != Hearts && suit != Diamonds && suit != Clubs {
		return 0, fmt.Errorf("NewCardString(%q): invalid suit", s)
	}

	return NewCard(rank, suit), nil
}

func (c *Card) MarshalJSON() ([]byte, error) {
	return []byte("\"" + c.String() + "\""), nil
}

func (c *Card) UnmarshalJSON(b []byte) (err error) {
	*c, err = NewCardString(string(b[1:3]))
	return err
}

func (c Card) String() string {
	return c.Rank().String() + c.Suit().String()
}

func (c Card) Rank() CardRank {
	return CardRank((uint32(c) >> 8) & 0xF)
}

func (c Card) Suit() CardSuit {
	return CardSuit((uint32(c) >> 12) & 0xF)
}

func (c Card) BitRank() uint32 {
	return (uint32(c) >> 16) & 0x1FFF
}
