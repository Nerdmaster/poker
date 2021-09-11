package poker

import (
	"testing"
)

func makeHand(s string) (*Hand, error) {
	var cards, err = ParseCards(s)
	return NewHand(cards), err
}

func TestHandEvaluate(t *testing.T) {
	var tests = map[string]struct {
		handValue uint16
		hand      string
		handRank  HandRank
	}{
		// Five-card hands
		"Five-card worst hand ever": {7462, "2s 3d 4c 5h 7h", HighCard},
		"Five-card ace high":        {6252, "As Ks Jc 7h 5d", HighCard},
		"Five-card pair":            {3448, "As Ac Jc 7h 5d", OnePair},
		"Five-card two pair":        {2497, "As Ac Jc Jd 5d", TwoPair},
		"Five-card three of a kind": {1636, "As Ac Ad Jd 5d", ThreeOfAKind},
		"Five-card straight":        {1600, "As Ks Qd Jh Td", Straight},
		"Five-card flush":           {1542, "Ts 7s 4s 3s 2s", Flush},
		"Five-card full house":      {298, "4s 4c 4d 2s 2h", FullHouse},
		"Five-card four of a kind":  {19, "As Ac Ad Ah 5h", FourOfAKind},
		"Five-card straight flush":  {1, "As Ks Qs Js Ts", StraightFlush},

		// Six-card hands
		"Six-card ace high":        {6252, "3d As Ks Jc 7h 5d", HighCard},
		"Six-card pair":            {3448, "3d As Ac Jc 7h 5d", OnePair},
		"Six-card two pair":        {2497, "3d As Ac Jc Jd 5d", TwoPair},
		"Six-card three of a kind": {1636, "3d As Ac Ad Jd 5d", ThreeOfAKind},
		"Six-card straight":        {1600, "3d As Ks Qd Jh Td", Straight},
		"Six-card flush":           {1542, "3d Ts 7s 4s 3s 2s", Flush},
		"Six-card full house":      {298, "3d 4s 4c 4d 2s 2h", FullHouse},
		"Six-card four of a kind":  {19, "3d As Ac Ad Ah 5h", FourOfAKind},
		"Six-card straight flush":  {1, "3d As Ks Qs Js Ts", StraightFlush},

		// Seven-card hands
		"Seven-card ace high":        {6252, "2d 3d As Ks Jc 7h 5d", HighCard},
		"Seven-card pair":            {3448, "2d 3d As Ac Jc 7h 5d", OnePair},
		"Seven-card two pair":        {2497, "2d 3d As Ac Jc Jd 5d", TwoPair},
		"Seven-card three of a kind": {1636, "2c 3d As Ac Ad Jd 5d", ThreeOfAKind},
		"Seven-card straight":        {1600, "2d 3d As Ks Qd Jh Td", Straight},
		"Seven-card flush":           {1542, "2d 3d Ts 7s 4s 3s 2s", Flush},
		"Seven-card full house":      {298, "2d 3d 4s 4c 4d 2s 2h", FullHouse},
		"Seven-card four of a kind":  {19, "2d 3d As Ac Ad Ah 5h", FourOfAKind},
		"Seven-card straight flush":  {1, "2d 3d As Ks Qs Js Ts", StraightFlush},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			var hand, err = makeHand(tc.hand)
			if err != nil {
				t.Errorf("Expected valid Hand using %q, got error: %s", tc.hand, err)
				return
			}

			var result *HandResult
			result, err = hand.Evaluate()
			if err != nil {
				t.Errorf("hand.Evaluate() for %q: error: %s", tc.hand, err)
				return
			}

			var handVal = result.Score
			var handRank = GetHandRank(handVal)
			if handVal != tc.handValue {
				t.Errorf("%s gave a hand value of %d; expected %d", hand, handVal, tc.handValue)
			}
			if handRank != tc.handRank {
				t.Errorf("%s gave a hand rank of %q; expected %q", hand, handRank, tc.handRank)
			}
		})
	}
}

func TestHandDescribe(t *testing.T) {
	var tests = map[string]struct {
		hand string
		desc string
	}{
		"worst hand ever":     {"2s 3d 4c 5h 7h", "Seven High"},
		"ace high":            {"Ks Jc Ac 7h 5d", "Ace High"},
		"pair":                {"As Ac Jc 7h 5d", "One Pair, Aces"},
		"two pair":            {"As Ac Jc Jd 5d", "Two Pair, Aces And Jacks"},
		"three of a kind":     {"Jd Ac Ad 5d As", "Three Of A Kind, Aces"},
		"straight":            {"9d Qd Ks Jh As Td", "Ace-High Straight"},
		"flush":               {"Ts 7s 4s 3s 2s", "Ten-High Flush"},
		"full house":          {"2s 4c 4d 4s 2h", "Full House, Fours Over Twos"},
		"four of a kind":      {"As Ac Ah 5h Ad", "Four Of A Kind, Aces"},
		"low straight flush":  {"As 3s 5s 4s 2s", "Five-High Straight Flush"},
		"high straight flush": {"Ks Qs As Js Ts", "Royal Flush"},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			var hand, err = makeHand(tc.hand)
			if err != nil {
				t.Errorf("Expected valid hand using %q, got error: %s", tc.hand, err)
			}

			var result *HandResult
			result, err = hand.Evaluate()
			if err != nil {
				t.Errorf("Error with hand.Evaluate() on %q: %s", tc.hand, err)
			}

			var got = result.Describe()
			if tc.desc != got {
				t.Errorf("Expected %q to be described as %q, got %q", tc.hand, tc.desc, got)
			}
		})
	}
}
