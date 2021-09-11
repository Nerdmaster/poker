package poker

import (
	"fmt"
	"strings"
)

// CardList is a simple alias for a list of ... wait for it... CARDS!
type CardList []Card

// ParseCards converts a string of whitespace-separated cards into a list
func ParseCards(s string) (CardList, error) {
	var cardStr = strings.Fields(s)
	var cards = make(CardList, len(cardStr))
	var err error
	for i, card := range cardStr {
		cards[i], err = NewCardString(card)
		if err != nil {
			return nil, fmt.Errorf("invalid hand %q: %w", s, err)
		}
	}

	return cards, nil
}

func (cards CardList) String() string {
	var slist []string
	for _, c := range cards {
		slist = append(slist, c.String())
	}

	return strings.Join(slist, " ")
}

func (cards *CardList) AddCard(c Card) {
	*cards = append(*cards, c)
}
