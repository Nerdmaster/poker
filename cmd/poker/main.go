package main

import (
	"log"
	"math/rand"
	"sort"
	"strconv"
	"time"

	"github.com/Nerdmaster/poker"
)

type Player struct {
	name   string
	hand   *poker.Hand
	result *poker.HandResult
}

func main() {
	var deck = poker.NewDeck(rand.NewSource(time.Now().UnixNano()))
	deck.Shuffle()

	var players []*Player
	for i := 0; i < 5; i++ {
		players = append(players, &Player{hand: poker.NewHand(nil)})
	}
	for i, p := range players {
		p.name = "Player " + strconv.Itoa(i+1)
		log.Printf("Dealing to %s", p.name)
		deck.Deal(p.hand)
		deck.Deal(p.hand)
	}

	log.Printf("Dealing community cards")
	var comm = deck.Draw(5)
	log.Printf("Community cards: %s", comm)

	for _, p := range players {
		var res, err = p.hand.Evaluate(comm...)
		if err != nil {
			log.Fatalf("Error evaluating hand %s: %s", p.hand, err)
		}
		p.result = res
	}

	sort.Slice(players, func(i, j int) bool {
		return players[i].result.Score < players[j].result.Score
	})

	for i, p := range players {
		log.Printf("In position %d, we have %s (%q) with %s (%s)", i+1, p.name, p.hand, p.result.Best5, p.result.Describe())
	}
}
