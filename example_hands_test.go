package poker_test

import (
	"fmt"
	"math/rand"

	"github.com/Nerdmaster/poker"
)

// ExampleTextHoldEm shows how to set up a single hand with two hole cards,
// then evaluate them against the five community cards.
func Example_texasHoldEmHand() {
	var deck = poker.NewDeck(rand.NewSource(3))
	deck.Shuffle()
	var hand = poker.NewHand(nil)

	deck.Deal(hand)
	deck.Deal(hand)

	var community = deck.Draw(5)
	fmt.Printf("Hole: %s\n", hand)
	fmt.Printf("Community: %s\n", community)
	fmt.Println()
	var res, err = hand.Evaluate(community...)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Hand: %s (%s)", res.Best5, res.Describe())

	// Output:
	// Hole: 9d 3d
	// Community: 3h Qh 7h 4h 4s
	//
	// Hand: 4h 4s 3d 3h Qh (Two Pair, Fours And Threes)
}

func Example_omahaTable() {
	type player struct {
		Name   string
		Hand   *poker.Hand
		Result *poker.HandResult
	}

	var deck = poker.NewDeck(rand.NewSource(3))
	deck.Shuffle()

	var players = []*player{
		{Name: "Alice", Hand: poker.NewHand(nil)},
		{Name: "Bob", Hand: poker.NewHand(nil)},
		{Name: "Frogurt", Hand: poker.NewHand(nil)},
		{Name: "Gatsby", Hand: poker.NewHand(nil)},
		{Name: "Tedros", Hand: poker.NewHand(nil)},
	}

	// Deal four cards to each player
	for i := 0; i < 4; i++ {
		for _, p := range players {
			deck.Deal(p.Hand)
		}
	}

	// Hack up Tedros to show this is indeed using Omaha rules: his hold cards
	// will be all the aces
	var fouraces, err = poker.ParseCards("Ah As Ad Ac")
	if err != nil {
		panic(err)
	}
	players[4].Hand = poker.NewHand(fouraces)

	var community = deck.Draw(5)
	var winner *player
	fmt.Printf("Community: %s\n", community)
	fmt.Println()
	for i, p := range players {
		if i > 0 {
			fmt.Println()
		}
		p.Result, err = p.Hand.Evaluate(community...)
		if err != nil {
			panic(err)
		}
		fmt.Println(p.Name + ":")
		fmt.Printf("  - Hole cards: %s\n", p.Hand)
		fmt.Printf("  - Hand: %s (%s)\n", p.Result.Describe(), p.Result.Best5)

		if winner == nil || p.Result.Score < winner.Result.Score {
			winner = p
		}
	}

	fmt.Printf("\nWinner: %s\n", winner.Name)

	// Output:
	// Community: 5s 8s Jc 8c 3s
	//
	// Alice:
	//   - Hole cards: 9d 4h 3c 4c
	//   - Hand: Two Pair, Eights And Fours (8c 8s 4c 4h Jc)
	//
	// Bob:
	//   - Hole cards: 3d 4s 7c 5c
	//   - Hand: Two Pair, Eights And Fives (8c 8s 5c 5s 7c)
	//
	// Frogurt:
	//   - Hole cards: 3h 8d 6s 9s
	//   - Hand: Full House, Eights Over Threes (8c 8d 8s 3h 3s)
	//
	// Gatsby:
	//   - Hole cards: Qh Tc 2c Kd
	//   - Hand: One Pair, Eights (8c 8s Kd Qh Jc)
	//
	// Tedros:
	//   - Hole cards: Ah As Ad Ac
	//   - Hand: Two Pair, Aces And Eights (Ah As 8c 8s Jc)
	//
	// Winner: Frogurt
}
