package poker

import (
	"encoding/json"
	"math"
	"math/rand"
	"strings"
	"testing"
	"time"
)

func TestRankString(t *testing.T) {
	var tests = map[string]struct {
		handValue uint16
		handRank  string
	}{
		"398: Flush":        {398, "Flush"},
		"2665: Two pair":    {2665, "Two Pair"},
		"6230: High Card":   {6230, "High Card"},
		"6529: High Card":   {6529, "High Card"},
		"6823: High Card":   {6823, "High Card"},
		"2669: Two Pair":    {2669, "Two Pair"},
		"4076: One Pair":    {4076, "One Pair"},
		"0: Straight Flush": {0, "Straight Flush"},
		"1607: Straight":    {1607, "Straight"},
		"7196: High Card":   {7196, "High Card"},
		"7221: High Card":   {7221, "High Card"},
		"6228: High Card":   {6228, "High Card"},
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
		"Five-card worst hand ever": {7462, `["2s", "3d", "4c", "5h", "7h"]`, HighCard},
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

		// Invalid hands
		"Four cards":  {math.MaxUint16, `["Ks", "Qs", "Js", "Ts"]`, HighCard},
		"Eight cards": {math.MaxUint16, `["Tc", "2d", "3d", "As", "Ks", "Qs", "Js", "Ts"]`, HighCard},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			var cards CardList
			var err = json.Unmarshal([]byte(tc.hand), &cards)
			if err != nil {
				t.Fatalf("Unmarshaling %q got an error: %s", tc.hand, err)
			}

			var handVal = cards.Evaluate()
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

func TestBestHand(t *testing.T) {
	var tests = map[string]struct {
		input string
		want  string
	}{
		// Five-card hands
		"Five card returns the same hand": {"As Ks Jc 7h 5d", "As Ks Jc 7h 5d"},

		// Seven-card hands
		"Seven-card ace high":        {"2d 3d As Ks Jc 7h 5d", "As Ks Jc 7h 5d"},
		"Seven-card pair":            {"2d 3d As Ac Jc 7h 5d", "As Ac Jc 7h 5d"},
		"Seven-card two pair":        {"2d 3d As Ac Jc Jd 5d", "As Ac Jc Jd 5d"},
		"Seven-card three of a kind": {"2c 3d 3h 3c Ad Jd 5d", "3d 3h 3c Ad Jd"},
		"Seven-card straight":        {"Jh 2d 3d Ks As Qd Td", "Jh Ks As Qd Td"},
		"Seven-card flush":           {"2s 3d Ts 7s 4s 3s 9s", "Ts 7s 4s 3s 9s"},
		"Seven-card full house":      {"2d 3d 4s 4c 4d 3s 3h", "3d 4s 4c 4d 3s"},
		"Seven-card four of a kind":  {"2d 3d 2s 2c Ad 2h 5h", "2d 2s 2c Ad 2h"},
		"Seven-card straight flush":  {"2s 3d As Ks 4s 3s 5s", "2s As 4s 3s 5s"},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			var hand CardList
			for _, c := range strings.Fields(tc.input) {
				hand = append(hand, newCardString(c))
			}
			var _, best = hand.BestHand()

			var got = CardList(best[:]).String()
			if got != tc.want {
				t.Fatalf("%s gave a best hand of %s; expected %s", tc.input, got, tc.want)
			}
		})
	}
}

var omahaTests = map[string]struct {
	handValue uint16
	hole      string
	community string
	handRank  HandRank
	bestHole  string
	bestComm  string
}{
	// Various ways of getting just a high card
	"high 1": {6251, "8c 3d As Ks", "Jc 7h 5d 2d 6s", HighCard, "As Ks", "Jc 7h 6s"},
	"high 2": {6251, "Qc 3d As Ks", "Jc 7h 5d 2d 6s", HighCard, "As Ks", "Jc 7h 6s"},
	"high 3": {6699, "Qc 3d Ks Ts", "Jc 7h 5d 2d 6s", HighCard, "Qc Ks", "Jc 7h 6s"},
	"Flush":  {722, "Js 8h Th 7d", "8c 3h Ah 5h 8s", Flush, "8h Th", "3h Ah 5h"},
	// No flush: you have to use two cards from your hand
	"No flush": {6699, "Qc 3d Ks Ts", "Jc 7c 5c 2c 6c", HighCard, "Qc Ks", "Jc 7c 6c"},
	// Lowest straight because you can't use three hole cards
	"Low str8": {1609, "Ac Qd Js 4h", "5c Kd Th 2s 3h", Straight, "Ac 4h", "5c 2s 3h"},
	// This gets the highest straight even though the combined cards are the same
	"High str8": {1600, "Ac Qd 3h 4h", "5c Kd Th 2s Js", Straight, "Ac Qd", "Kd Th Js"},

	// And now a few random hands just for extra assurance
	"Two pair":   {2853, "2c Jd 8h 4h", "Jc 8d 3h Qs Ad", TwoPair, "Jd 8h", "Jc 8d Ad"},
	"One pair":   {5093, "7h As Jc 6d", "6s Kd 2s 4s 3d", OnePair, "As 6d", "6s Kd 4s"},
	"Full House": {288, "Kc 2s 4c Qd", "3h 7c 4d 4s Kd", FullHouse, "Kc 4c", "4d 4s Kd"},
}

func TestEvaluateOmaha(t *testing.T) {
	for name, tc := range omahaTests {
		t.Run(name, func(t *testing.T) {
			var hole, community CardList
			for _, c := range strings.Fields(tc.hole) {
				hole = append(hole, newCardString(c))
			}
			for _, c := range strings.Fields(tc.community) {
				community = append(community, newCardString(c))
			}

			var handVal = hole.EvaluateOmaha(community)
			var handRank = GetHandRank(handVal)
			if handRank != tc.handRank {
				t.Fatalf("%s,%s gave a hand rank of %q; expected %q", hole, community, handRank, tc.handRank)
			}
			if handVal != tc.handValue {
				t.Fatalf("%s,%s gave a hand value of %d; expected %d", hole, community, handVal, tc.handValue)
			}
		})
	}
}

func TestBestOmahaHand(t *testing.T) {
	for name, tc := range omahaTests {
		t.Run(name, func(t *testing.T) {
			var hole, community CardList
			for _, c := range strings.Fields(tc.hole) {
				hole = append(hole, newCardString(c))
			}
			for _, c := range strings.Fields(tc.community) {
				community = append(community, newCardString(c))
			}

			var _, bestH, bestC = hole.BestOmahaHand(community)
			var gotHole = CardList(bestH[:]).String()
			var gotComm = CardList(bestC[:]).String()
			if gotHole != tc.bestHole || gotComm != tc.bestComm {
				t.Fatalf("%s,%s: expected %s,%s to be best, but got %s,%s", hole, community, tc.bestHole, tc.bestComm, gotHole, gotComm)
			}
		})
	}
}

func BenchmarkEvalFiveFast(b *testing.B) {
	var deck *Deck
	var hands = make([]CardList, 100)
	for i := 0; i < 100; i++ {
		deck = NewDeck(rand.NewSource(time.Now().UnixNano()))
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
	var hands = make([]CardList, 100)
	for i := 0; i < 100; i++ {
		deck = NewDeck(rand.NewSource(time.Now().UnixNano()))
		deck.Shuffle()
		hands[i] = deck.Draw(5)
	}

	var hl = len(hands)
	for i := 0; i < b.N; i++ {
		var hand = hands[i%hl]
		hand.Evaluate()
	}
}

func BenchmarkEvaluateSeven(b *testing.B) {
	var deck *Deck
	var hands = make([]CardList, 100)
	for i := 0; i < 100; i++ {
		deck = NewDeck(rand.NewSource(time.Now().UnixNano()))
		deck.Shuffle()
		hands[i] = deck.Draw(7)
	}

	var hl = len(hands)
	for i := 0; i < b.N; i++ {
		var hand = hands[i%hl]
		hand.Evaluate()
	}
}

func BenchmarkBestHandSeven(b *testing.B) {
	var deck *Deck
	var hands = make([]CardList, 100)
	for i := 0; i < 100; i++ {
		deck = NewDeck(rand.NewSource(time.Now().UnixNano()))
		deck.Shuffle()
		hands[i] = deck.Draw(7)
	}

	var hl = len(hands)
	for i := 0; i < b.N; i++ {
		var hand = hands[i%hl]
		hand.BestHand()
	}
}

func BenchmarkEvaluateOmaha(b *testing.B) {
	var deck *Deck
	var hands = make([]CardList, 100)
	for i := 0; i < 100; i++ {
		deck = NewDeck(rand.NewSource(time.Now().UnixNano()))
		deck.Shuffle()
		hands[i] = deck.Draw(9)
	}

	var hl = len(hands)
	for i := 0; i < b.N; i++ {
		var hand = hands[i%hl]
		CardList(hand[:4]).EvaluateOmaha(hand[4:])
	}
}

func BenchmarkBestOmahaHand(b *testing.B) {
	var deck *Deck
	var hands = make([]CardList, 100)
	for i := 0; i < 100; i++ {
		deck = NewDeck(rand.NewSource(time.Now().UnixNano()))
		deck.Shuffle()
		hands[i] = deck.Draw(9)
	}

	var hl = len(hands)
	for i := 0; i < b.N; i++ {
		var hand = hands[i%hl]
		CardList(hand[:4]).EvaluateOmaha(hand[4:])
	}
}
