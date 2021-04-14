package poker

import (
	"encoding/json"
	"testing"
)

func TestNewCard(t *testing.T) {
	var tests = map[string]struct {
		rank CardRank
		suit CardSuit
		want Card
	}{
		"Ace of hearts":  {Ace, Hearts, 268446761},
		"King of spades": {King, Spades, 134224677},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			var c = NewCard(tc.rank, tc.suit)
			if c != tc.want {
				t.Fatalf("Expected NewCard(%s, %s) to give card %d, but got %d (%q)", tc.rank, tc.suit, tc.want, c, c)
			}
		})
	}
}

func TestNewCardString(t *testing.T) {
	var tests = map[string]struct {
		input string
		want  Card
	}{
		"Ace of hearts":  {"Ah", 268446761},
		"King of spades": {"Ks", 134224677},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			var c, err = NewCardString(tc.input)
			if err != nil {
				t.Fatalf("Unexpected error with %q: %s", tc.input, err)
			}
			if c != tc.want {
				t.Fatalf("Expected %q to give card %d, but got %d (%q)", tc.input, tc.want, c, c)
			}
		})
	}
}

// This hack lets us ignore errors since invalid card text in *tests* better
// not be happening without a pretty obvious typo we fix immediately
func newCardString(s string) Card {
	var c, _ = NewCardString(s)
	return c
}

func TestMarshalJSON(t *testing.T) {
	var cards = []Card{
		newCardString("Ah"),
		newCardString("Kh"),
		newCardString("Qh"),
		newCardString("Jh"),
		newCardString("Th"),
	}

	var b, err = json.Marshal(cards)
	if err != nil {
		t.Fatalf("Unexpected error marshaling cards: %s", err)
	}
	var want = `["Ah","Kh","Qh","Jh","Th"]`
	if string(b) != want {
		t.Fatalf("Expected %q, but got %q", want, string(b))
	}
}

func TestUnmarshalJSON(t *testing.T) {
	var cards []Card
	var data = `["Ah","Kh","Qh","Jh","Th"]`

	var err = json.Unmarshal([]byte(data), &cards)
	if err != nil {
		t.Fatalf("Unexpected error unmarshaling %q: %s", data, err)
	}

	if len(cards) != 5 {
		t.Fatalf("Got %d cards, expected 5", len(cards))
	}

	if cards[0] != newCardString("Ah") {
		t.Fatalf("Card 0 was %q; expected Ah", cards[0])
	}
	if cards[1] != newCardString("Kh") {
		t.Fatalf("Card 1 was %q; expected Kh", cards[0])
	}
	if cards[2] != newCardString("Qh") {
		t.Fatalf("Card 2 was %q; expected Qh", cards[0])
	}
	if cards[3] != newCardString("Jh") {
		t.Fatalf("Card 3 was %q; expected Jh", cards[0])
	}
	if cards[4] != newCardString("Th") {
		t.Fatalf("Card 4 was %q; expected Th", cards[0])
	}
}

func TestString(t *testing.T) {
	var c = newCardString("3s")
	if c.String() != "3s" {
		t.Fatalf("3s gave us %q for its string", c)
	}
}

func TestBitRank(t *testing.T) {
	var c = newCardString("Ks")
	var br = c.BitRank()
	if br != 2048 {
		t.Fatalf("%q bitrank: got %d, expected 2048", c, br)
	}
}
