package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/Nerdmaster/poker"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	var deck = poker.NewDeck()
	deck.Shuffle()
	var hand = deck.Draw(7)
	fmt.Println(hand)

	var val = poker.Evaluate(hand)
	fmt.Println(val)
	fmt.Println(poker.GetHandRank(val))
}
