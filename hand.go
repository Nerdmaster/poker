package poker

import (
	"fmt"
	"math"
	"strings"
)

// HandError is a simple string wrapper to make constants into usable and
// meaningful returned error values
type HandError string

// Hand errors
const (
	ErrEmptyDeck        HandError = "cannot draw from empty deck"
	ErrInvalidCardCount HandError = "invalid card count"
)

func (e HandError) Error() string {
	return string(e)
}

// A Hand is just a CardList with user-specific behaviors. It could be a
// five-card draw hand, hole cards, etc.
//
// Hand is primarily built for higher-level interactions than the low-level
// Card-based functions, but it's really nothing more than a list of cards with
// some simple behaviors.
//
// A Hand's zero value is usable, but cards will have to be added from a deck
// before it can be useful.
type Hand struct {
	cards CardList
}

// NewHand takes a CardList and turns it into a usable hand
func NewHand(cards CardList) *Hand {
	return &Hand{cards: cards}
}

// String returns a human-readable(ish) string representing the hand
func (h *Hand) String() string {
	var list = make([]string, len(h.cards))
	for i, c := range h.cards {
		list[i] = c.String()
	}
	return strings.Join(list, " ")
}

// Evaluate computes the best hand and its score based on cards in this hand
// and the community cards, if any. If there are no community cards, this is a
// "standard" poker hand (for five-card draw, seven-card stud, etc.).
//
// If the hand is invalid for evaluation (fewer than five cards total, 3 hole
// cards but community cards were offered up, etc.), the score will be the
// worst possible (MaxUint16), and there will be no description of the hand.
func (h *Hand) Evaluate(community ...Card) (hr *HandResult, err error) {
	hr = &HandResult{}

	// Copy cards, don't just reuse the slices
	hr.Hand = make(CardList, len(h.cards))
	copy(hr.Hand, h.cards)
	hr.Community = make(CardList, len(community))
	copy(hr.Community, community)

	if len(community) == 0 {
		err = hr.evaluateRaw()
	} else {
		switch len(h.cards) {
		case 2:
			err = hr.evaluateTexas()
		case 4:
			err = hr.evaluateOmaha()
		default:
			return nil, fmt.Errorf("%w: hole cards must be two or four when community cards are present", ErrInvalidCardCount)
		}
	}

	if err != nil {
		return nil, fmt.Errorf("error evaluating hand: %w", err)
	}
	if hr.Score == math.MaxUint16 {
		return nil, fmt.Errorf("unknown error evaluating hand")
	}

	hr.Best5 = CardList(hr.best[:])
	hr.Rank = GetHandRank(hr.Score)
	hr.sort()
	return hr, nil
}

// AddCard puts the card into this player's hand
func (h *Hand) AddCard(c Card) {
	h.cards = append(h.cards, c)
}
