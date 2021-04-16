package poker

import (
	"math/rand"
)

// A Deck is a magical list of cards that happens to be able to initialize
// itself to a standard 52-card setup as well as be shuffled and have cards
// drawn, removing them from the deck.
type Deck struct {
	rnd *rand.Rand
	cards   []Card
}

// NewDeck returns a deck of 52 cards.  These are not shuffled in any way.
// It's like you just went to the store and bought them, but you don't have to
// spend time getting them out of that $#@*ING plastic wrap.
//
// rndSource can be any implementation of math/rand.Source; it is never seeded
// here, so cryptographically secure implementations are fine.
//
// A Deck is *not* safe for concurrent use.
func NewDeck(rndSource rand.Source) *Deck {
	var deck = &Deck{rnd: rand.New(rndSource)}
	deck.Reset()
	return deck
}

// Shuffle does what you think - randomizes the cards in the deck.  To
// re-initialize the deck with a full set of cards, use Reset().
func (d *Deck) Shuffle() {
	d.rnd.Shuffle(len(d.cards), func(i, j int) {
		d.cards[i], d.cards[j] = d.cards[j], d.cards[i]
	})
}

// Reset puts all cards back into the deck in their original order
func (d *Deck) Reset() {
	d.cards = make([]Card, 52)
	var idx = 0
	for rank := Deuce; rank <= Ace; rank++ {
		for _, suit := range []CardSuit{Spades, Hearts, Diamonds, Clubs} {
			d.cards[idx] = NewCard(rank, suit)
			idx++
		}
	}
}

// Draw returns up to n cards.  If n is larger than the number of cards left in
// the deck, only that many cards are returned.  A zero-length slice can be
// returned if the deck is empty.
func (d *Deck) Draw(n int) (cards []Card) {
	if len(d.cards) < n {
		n = len(d.cards)
	}

	cards, d.cards = d.cards[:n], d.cards[n:]
	return cards
}

// Count returns the number of cards left in the deck
func (d *Deck) Count() int {
	return len(d.cards)
}

// Empty returns true if the deck has no more cards
func (d *Deck) Empty() bool {
	return len(d.cards) == 0
}
