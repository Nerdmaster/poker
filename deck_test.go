package poker

import (
	"math/rand"
	"testing"
)

func TestShuffle(t *testing.T) {
	var d1 = NewDeck(rand.NewSource(0))
	var d2 = NewDeck(rand.NewSource(0))

	var diff bool
	d1.Shuffle()
	for i, card := range d1.cards {
		if card != d2.cards[i] {
			diff = true
			break
		}
	}

	if !diff {
		t.Fatalf("Expected shuffled deck to differ from new deck, but it didn't!")
	}
}

func TestDraw(t *testing.T) {
	var deck = NewDeck(rand.NewSource(0))

	var cards = deck.Draw(5)
	if len(cards) != 5 {
		t.Fatalf("Expected deck.Draw(5) to return five cards, but got %d", len(cards))
	}
	if len(deck.cards) != 47 {
		t.Fatalf("Expected deck.Draw(5) to result in 47 cards remaining, but we have %d", len(deck.cards))
	}
}

func TestEmpty(t *testing.T) {
	var deck = NewDeck(rand.NewSource(0))
	if deck.Empty() {
		t.Fatalf("Newly initialized deck was empty")
	}

	deck.Draw(51)
	if deck.Empty() {
		t.Fatalf("Deck with 51 cards drawn was empty")
	}

	deck.Draw(1)
	if !deck.Empty() {
		t.Fatalf("Deck with 52 cards drawn wasn't reporting being empty")
	}
}
