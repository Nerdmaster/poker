package poker

import (
	"math/rand"
)

// A Deck is a magical list of cards that happens to shuffle them and allow
// them to be drawn into a hand
type Deck struct {
	cards []Card
}

// NewDeck returns a deck of 52 cards.  These are not shuffled in any way.
// It's like you just went to the store and bought them, but you don't have to
// spend time getting them out of that $#@*ING plastic wrap.
func NewDeck() *Deck {
	var deck = &Deck{}
	deck.initialize()
	return deck
}

// Shuffle does what you think - randomizes the cards in the deck.  rand.Seed
// is *not* called here, so seed it yourself to ensure quality randomization.
// As this uses math/rand, this operation is *not* cryptographically secure.
// This deck is far more useful as a toy to play with than a high-volume poker
// server where real currency is exchanged.
func (d *Deck) Shuffle() {
	rand.Shuffle(len(d.cards), func(i, j int) {
		d.cards[i], d.cards[j] = d.cards[j], d.cards[i]
	})
}

// Draw returns up to n cards.  If n is larger than the number of cards left in
// the deck, only that many cards is returned.
func (d *Deck) Draw(n int) (cards []Card) {
	if len(d.cards) < n {
		n = len(d.cards)
	}

	cards, d.cards = d.cards[:n], d.cards[n:]
	return cards
}

// Empty returns true if the deck has no more cards
func (d *Deck) Empty() bool {
	return len(d.cards) == 0
}

func (d *Deck) initialize() {
	d.cards = make([]Card, 52)
	var idx = 0
	for rank := Deuce; rank <= Ace; rank++ {
		for _, suit := range []CardSuit{Spades, Hearts, Diamonds, Clubs} {
			d.cards[idx] = NewCard(rank, suit)
			idx++
		}
	}
}
