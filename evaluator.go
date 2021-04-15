package poker

import (
	"math"
)

func evalFiveFast(c1, c2, c3, c4, c5 Card) uint16 {
	var q = (c1 | c2 | c3 | c4 | c5) >> 16
	if (c1 & c2 & c3 & c4 & c5 & 0xf000) != 0 {
		return flushes[q]
	}
	if unique5[q] != 0 {
		return unique5[q]
	}

	var product = (c1 & 0xff) * (c2 & 0xff) * (c3 & 0xff) * (c4 & 0xff) * (c5 & 0xff)
	return hashValues[findFast(uint32(product))]
}

func findFast(prod uint32) uint32 {
	var a, b, r uint32
	prod += 0xe91aaa35
	prod ^= prod >> 16
	prod += prod << 8
	prod ^= prod >> 4
	b = (prod >> 8) & 0x1ff
	a = (prod + (prod << 2)) >> 19
	r = a ^ uint32(hashAdjust[b])
	return r
}

// Evaluate returns the score for the best five-card poker hand found.  Hands
// can be 5, 6, or 7 cards, otherwise the return will be math.MaxUint16.
func Evaluate(c []Card) uint16 {
	if len(c) == 5 {
		return evalFiveFast(c[0], c[1], c[2], c[3], c[4])
	}
	if len(c) < 5 || len(c) > 8 {
		return math.MaxUint16
	}

	return evalMore(c)
}

// BestHand returns the score as well as the best five-card hand found
func BestHand(cards []Card) (score uint16, best [5]Card) {
	score = math.MaxUint16
	if len(cards) < 5 || len(cards) > 7 {
		return
	}
	if len(cards) == 5 {
		best[0] = cards[0]
		best[1] = cards[1]
		best[2] = cards[2]
		best[3] = cards[3]
		best[4] = cards[4]
		return evalFiveFast(cards[0], cards[1], cards[2], cards[3], cards[4]), best
	}

	var perms = perms7
	if len(cards) == 6 {
		perms = perms6
	}

	for _, perm := range perms {
		var val = evalFiveFast(
			cards[perm[0]],
			cards[perm[1]],
			cards[perm[2]],
			cards[perm[3]],
			cards[perm[4]],
		)
		if val < score {
			score = val
			best[0] = cards[perm[0]]
			best[1] = cards[perm[1]]
			best[2] = cards[perm[2]]
			best[3] = cards[perm[3]]
			best[4] = cards[perm[4]]
		}
	}

	return
}

var perms6 = [][5]int{
	{0, 1, 2, 3, 4},
	{0, 1, 2, 3, 5},
	{0, 1, 2, 4, 5},
	{0, 1, 3, 4, 5},
	{0, 2, 3, 4, 5},
	{1, 2, 3, 4, 5},
}

var perms7 = [][5]int{
	{0, 1, 2, 3, 4},
	{0, 1, 2, 3, 5},
	{0, 1, 2, 3, 6},
	{0, 1, 2, 4, 5},
	{0, 1, 2, 4, 6},
	{0, 1, 2, 5, 6},
	{0, 1, 3, 4, 5},
	{0, 1, 3, 4, 6},
	{0, 1, 3, 5, 6},
	{0, 1, 4, 5, 6},
	{0, 2, 3, 4, 5},
	{0, 2, 3, 4, 6},
	{0, 2, 3, 5, 6},
	{0, 2, 4, 5, 6},
	{0, 3, 4, 5, 6},
	{1, 2, 3, 4, 5},
	{1, 2, 3, 4, 6},
	{1, 2, 3, 5, 6},
	{1, 2, 4, 5, 6},
	{1, 3, 4, 5, 6},
	{2, 3, 4, 5, 6},
}

// all permutations of the four hole cards in Omaha - exactly two must be used
var omahaHolePerms = [][2]int{
	{0, 1},
	{0, 2},
	{0, 3},
	{1, 2},
	{1, 3},
	{2, 3},
}

// all permutations of the five community cards in Omaha - exactly three must be used
var omahaCommunityPerms = [][3]int{
	{0, 1, 2},
	{0, 1, 3},
	{0, 1, 4},
	{0, 2, 3},
	{0, 2, 4},
	{0, 3, 4},
	{1, 2, 3},
	{1, 2, 4},
	{1, 3, 4},
	{2, 3, 4},
}

func evalMore(cards []Card) uint16 {
	var minimum uint16 = math.MaxUint16

	var perms = perms7
	if len(cards) == 6 {
		perms = perms6
	}

	for _, perm := range perms {
		var score = evalFiveFast(
			cards[perm[0]],
			cards[perm[1]],
			cards[perm[2]],
			cards[perm[3]],
			cards[perm[4]],
		)
		if score < minimum {
			minimum = score
		}
	}

	return minimum
}

// EvaluateOmaha returns the best five-card hand value that can be generated
// give four hole cards and five community cards.
//
// Omaha has nine cards, but the rules require you to use exactly two of them
// to make a hand.  This might seem complicated, but it drastically reduces the
// permutations compared to a full nine-card evaluation.
func EvaluateOmaha(hole, community []Card) uint16 {
	var minimum uint16 = math.MaxUint16
	if len(hole) != 4 || len(community) != 5 {
		return minimum
	}

	for _, holeP := range omahaHolePerms {
		for _, commP := range omahaCommunityPerms {
			var score = evalFiveFast(
				hole[holeP[0]],
				hole[holeP[1]],
				community[commP[0]],
				community[commP[1]],
				community[commP[2]],
			)
			if score < minimum {
				minimum = score
			}
		}
	}
	return minimum
}

// BestOmahaHand returns the score as well as the actual hole and community
// cards used.  A lot of code is copied from EvaluateOmaha so that that
// function can be streamlined for cases where speed is of the essence, where
// this is more useful when you need to present what actually "won" out of the
// many permutations of Omaha hand possibilities.
func BestOmahaHand(hole, community []Card) (score uint16, bestH [2]Card, bestC [3]Card) {
	score = math.MaxUint16
	if len(hole) != 4 || len(community) != 5 {
		return
	}

	for _, holeP := range omahaHolePerms {
		for _, commP := range omahaCommunityPerms {
			var eval = evalFiveFast(
				hole[holeP[0]],
				hole[holeP[1]],
				community[commP[0]],
				community[commP[1]],
				community[commP[2]],
			)
			if eval < score {
				bestH[0] = hole[holeP[0]]
				bestH[1] = hole[holeP[1]]
				bestC[0] = community[commP[0]]
				bestC[1] = community[commP[1]]
				bestC[2] = community[commP[2]]
				score = eval
			}
		}
	}
	return
}
