package poker

import "fmt"

// HandResult is the complex data created after analyzing a hand. It contains
// the source cards (user and community), the five cards that made the best
// hand, sorted for readability, a raw score, and a human-friendly description.
type HandResult struct {
	Hand        CardList
	Community   CardList
	best        [5]Card
	Best5       CardList
	Rank        HandRank
	Description string
	Score       uint16
}

func (hr *HandResult) evaluateRaw() error {
	if len(hr.Hand) < 5 || len(hr.Hand) > 7 {
		return fmt.Errorf("evaluateRaw(): %w", ErrInvalidCardCount)
	}

	hr.Score, hr.best = BestHand(hr.Hand)
	return nil
}

func (hr *HandResult) evaluateTexas() error {
	if len(hr.Community) < 3 || len(hr.Community) > 5 {
		return fmt.Errorf("evaluateTexas(): %w", ErrInvalidCardCount)
	}

	var eval = append(hr.Hand, hr.Community...)
	hr.Score, hr.best = BestHand(eval)
	return nil
}

func (hr *HandResult) evaluateOmaha() error {
	if len(hr.Community) < 3 || len(hr.Community) > 5 {
		return fmt.Errorf("evaluateOmaha(): %w", ErrInvalidCardCount)
	}

	var bestH [2]Card
	var bestC [3]Card
	hr.Score, bestH, bestC = BestOmahaHand(hr.Hand, hr.Community)
	hr.best = [5]Card{bestH[0], bestH[1], bestC[0], bestC[1], bestC[2]}
	return nil
}

// sort takes the hand result and makes it human-friendly based on its rank.
// "sorts" the best hand's cards so they're easy to read (e.g., "Ah 2d Ac 3s
// As" becomes "Ah Ac As 3s 2d") and marks the hand as having been evaluated.
//
// If the score was invalid upon calling this method, no sorting takes place
// and "evaluated" is set to false.
func (hr *HandResult) sort() {
	switch GetHandRank(hr.Score) {
	case StraightFlush, Straight:
		// If we have two cards both greater than five, Ace must be high, otherwise
		// it's low
		if hr.Best5[0].Rank() > Five && hr.Best5[1].Rank() > Five {
			hr.Best5.SortAceHigh()
		} else {
			hr.Best5.SortAceLow()
		}
	case FourOfAKind, FullHouse, ThreeOfAKind, TwoPair, OnePair:
		hr.Best5.SortGroups()
	case Flush, HighCard:
		hr.Best5.SortAceHigh()
	}
}

// Describe gives an explanation about the hand: "Full House, Aces Over Kings",
// "Two pair, Kings And Threes", etc.
//
// If this is called without one of the Evaluate methods first having been
// called, the hand is described as "N/A".
func (hr *HandResult) Describe() string {
	var high = hr.Best5[0].Rank()
	var low = hr.Best5[4].Rank()
	var base = hr.Rank.String()

	switch hr.Rank {
	case StraightFlush, Straight:
		if hr.Rank == StraightFlush && high == Ace {
			return "Royal Flush"
		}
		return high.Name() + "-High " + base
	case FourOfAKind:
		return base + ", " + high.Plural()
	case FullHouse:
		return base + ", " + high.Plural() + " Over " + low.Plural()
	case ThreeOfAKind:
		return base + ", " + high.Plural()
	case TwoPair:
		return base + ", " + high.Plural() + " And " + hr.Best5[2].Rank().Plural()
	case OnePair:
		return base + ", " + high.Plural()
	case Flush:
		return high.Name() + "-High " + base
	case HighCard:
		return high.Name() + " High"
	}

	panic("ERROR: Unknown hand rank!")
}
