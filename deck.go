package poker

import (
	"math/rand"
)

type Deck struct {
	cards []Card
}

func NewDeck() *Deck {
	var deck = &Deck{}
	deck.initialize()
	return deck
}

func (d *Deck) Shuffle() {
	rand.Shuffle(len(d.cards), func(i, j int) {
		d.cards[i], d.cards[j] = d.cards[j], d.cards[i]
	})
}

func (d *Deck) Draw(n int) (cards []Card) {
	if len(d.cards) < n {
		n = len(d.cards)
	}

	cards, d.cards = d.cards[:n], d.cards[n:]
	return cards
}

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
