package poker

import (
	"encoding/json"
	"testing"
)

func TestRankString(t *testing.T) {
	var tests = map[string]struct {
		handValue uint16
		handRank  string
	}{
		"398: Flush":      {398, "Flush"},
		"2665: Two pair":  {2665, "Two Pair"},
		"6230: High Card": {6230, "High Card"},
		"6529: High Card": {6529, "High Card"},
		"6823: High Card": {6823, "High Card"},
		"2669: Two Pair":  {2669, "Two Pair"},
		"4076: One Pair":  {4076, "One Pair"},
		"7196: High Card": {7196, "High Card"},
		"7221: High Card": {7221, "High Card"},
		"6228: High Card": {6228, "High Card"},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			var hr = GetHandRank(tc.handValue)
			if tc.handRank != hr.String() {
				t.Fatalf("Expected %q, but got %q", tc.handRank, hr)
			}
		})
	}
}

func TestEvaluate(t *testing.T) {
	var tests = map[string]struct {
		handValue uint16
		hand      string
		handRank  HandRank
	}{
		// Five-card hands
		"Five-card ace high":        {6252, `["As", "Ks", "Jc", "7h", "5d"]`, HighCard},
		"Five-card pair":            {3448, `["As", "Ac", "Jc", "7h", "5d"]`, OnePair},
		"Five-card two pair":        {2497, `["As", "Ac", "Jc", "Jd", "5d"]`, TwoPair},
		"Five-card three of a kind": {1636, `["As", "Ac", "Ad", "Jd", "5d"]`, ThreeOfAKind},
		"Five-card straight":        {1600, `["As", "Ks", "Qd", "Jh", "Td"]`, Straight},
		"Five-card flush":           {1542, `["Ts", "7s", "4s", "3s", "2s"]`, Flush},
		"Five-card full house":      {298, `["4s", "4c", "4d", "2s", "2h"]`, FullHouse},
		"Five-card four of a kind":  {19, `["As", "Ac", "Ad", "Ah", "5h"]`, FourOfAKind},
		"Five-card straight flush":  {1, `["As", "Ks", "Qs", "Js", "Ts"]`, StraightFlush},

		// Six-card hands
		"Six-card ace high":        {6252, `["3d", "As", "Ks", "Jc", "7h", "5d"]`, HighCard},
		"Six-card pair":            {3448, `["3d", "As", "Ac", "Jc", "7h", "5d"]`, OnePair},
		"Six-card two pair":        {2497, `["3d", "As", "Ac", "Jc", "Jd", "5d"]`, TwoPair},
		"Six-card three of a kind": {1636, `["3d", "As", "Ac", "Ad", "Jd", "5d"]`, ThreeOfAKind},
		"Six-card straight":        {1600, `["3d", "As", "Ks", "Qd", "Jh", "Td"]`, Straight},
		"Six-card flush":           {1542, `["3d", "Ts", "7s", "4s", "3s", "2s"]`, Flush},
		"Six-card full house":      {298, `["3d", "4s", "4c", "4d", "2s", "2h"]`, FullHouse},
		"Six-card four of a kind":  {19, `["3d", "As", "Ac", "Ad", "Ah", "5h"]`, FourOfAKind},
		"Six-card straight flush":  {1, `["3d", "As", "Ks", "Qs", "Js", "Ts"]`, StraightFlush},

		// Seven-card hands
		"Seven-card ace high":        {6252, `["2d", "3d", "As", "Ks", "Jc", "7h", "5d"]`, HighCard},
		"Seven-card pair":            {3448, `["2d", "3d", "As", "Ac", "Jc", "7h", "5d"]`, OnePair},
		"Seven-card two pair":        {2497, `["2d", "3d", "As", "Ac", "Jc", "Jd", "5d"]`, TwoPair},
		"Seven-card three of a kind": {1636, `["2c", "3d", "As", "Ac", "Ad", "Jd", "5d"]`, ThreeOfAKind},
		"Seven-card straight":        {1600, `["2d", "3d", "As", "Ks", "Qd", "Jh", "Td"]`, Straight},
		"Seven-card flush":           {1542, `["2d", "3d", "Ts", "7s", "4s", "3s", "2s"]`, Flush},
		"Seven-card full house":      {298, `["2d", "3d", "4s", "4c", "4d", "2s", "2h"]`, FullHouse},
		"Seven-card four of a kind":  {19, `["2d", "3d", "As", "Ac", "Ad", "Ah", "5h"]`, FourOfAKind},
		"Seven-card straight flush":  {1, `["2d", "3d", "As", "Ks", "Qs", "Js", "Ts"]`, StraightFlush},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			var cards []Card
			var err = json.Unmarshal([]byte(tc.hand), &cards)
			if err != nil {
				t.Fatalf("Unmarshaling %q got an error: %s", tc.hand, err)
			}

			var handVal = Evaluate(cards)
			var handRank = GetHandRank(handVal)
			if handVal != tc.handValue {
				t.Fatalf("%s gave a hand value of %d; expected %d", cards, handVal, tc.handValue)
			}
			if handRank != tc.handRank {
				t.Fatalf("%s gave a hand rank of %q; expected %q", cards, handRank, tc.handRank)
			}
		})
	}
}

func BenchmarkEvalFiveFast(b *testing.B) {
	var deck *Deck
	var hands = make([][]Card, 100)
	for i := 0; i < 100; i++ {
		deck = NewDeck()
		deck.Shuffle()
		hands[i] = deck.Draw(5)
	}

	var hl = len(hands)
	for i := 0; i < b.N; i++ {
		var hand = hands[i%hl]
		evalFiveFast(hand[0], hand[1], hand[2], hand[3], hand[4])
	}
}

func BenchmarkEvaluateFive(b *testing.B) {
	var deck *Deck
	var hands = make([][]Card, 100)
	for i := 0; i < 100; i++ {
		deck = NewDeck()
		deck.Shuffle()
		hands[i] = deck.Draw(5)
	}

	var hl = len(hands)
	for i := 0; i < b.N; i++ {
		var hand = hands[i%hl]
		Evaluate(hand)
	}
}

func BenchmarkEvaluateSeven(b *testing.B) {
	var deck *Deck
	var hands = make([][]Card, 100)
	for i := 0; i < 100; i++ {
		deck = NewDeck()
		deck.Shuffle()
		hands[i] = deck.Draw(7)
	}

	var hl = len(hands)
	for i := 0; i < b.N; i++ {
		var hand = hands[i%hl]
		Evaluate(hand)
	}
}
