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

// Evaluate returns the score for the best five-card poker hand found.  The
// lower the score, the better the hand, with a royal flush being 1, and the
// worst-possible high card (2, 3, 4, 5, 7) being 7462.
//
// Hands can be 5, 6, or 7 cards, otherwise the return will be math.MaxUint16.
func (cl CardList) Evaluate() uint16 {
	if len(cl) == 5 {
		return evalFiveFast(cl[0], cl[1], cl[2], cl[3], cl[4])
	}
	if len(cl) < 5 || len(cl) > 7 {
		return math.MaxUint16
	}

	return cl.evalMore()
}

// BestHand returns the same score as Evaluate, but also the full set of five
// cards that made up the best hand.  In cases where you only have five card,
// this is useless, but it can be more helpful when trying to see *why* a given
// seven-card hand won, especially for Poker novices.
func (cl CardList) BestHand() (score uint16, best [5]Card) {
	score = math.MaxUint16
	if len(cl) < 5 || len(cl) > 7 {
		return
	}
	if len(cl) == 5 {
		best[0] = cl[0]
		best[1] = cl[1]
		best[2] = cl[2]
		best[3] = cl[3]
		best[4] = cl[4]
		return evalFiveFast(cl[0], cl[1], cl[2], cl[3], cl[4]), best
	}

	var perms = perms7
	if len(cl) == 6 {
		perms = perms6
	}

	for _, perm := range perms {
		var val = evalFiveFast(
			cl[perm[0]],
			cl[perm[1]],
			cl[perm[2]],
			cl[perm[3]],
			cl[perm[4]],
		)
		if val < score {
			score = val
			best[0] = cl[perm[0]]
			best[1] = cl[perm[1]]
			best[2] = cl[perm[2]]
			best[3] = cl[perm[3]]
			best[4] = cl[perm[4]]
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

func (cl CardList) evalMore() uint16 {
	var minimum uint16 = math.MaxUint16

	var perms = perms7
	if len(cl) == 6 {
		perms = perms6
	}

	for _, perm := range perms {
		var score = evalFiveFast(
			cl[perm[0]],
			cl[perm[1]],
			cl[perm[2]],
			cl[perm[3]],
			cl[perm[4]],
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
func (cl CardList) EvaluateOmaha(community CardList) uint16 {
	var minimum uint16 = math.MaxUint16
	if len(cl) != 4 || len(community) != 5 {
		return minimum
	}

	for _, holeP := range omahaHolePerms {
		for _, commP := range omahaCommunityPerms {
			var score = evalFiveFast(
				cl[holeP[0]],
				cl[holeP[1]],
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
func (cl CardList) BestOmahaHand(community []Card) (score uint16, bestH [2]Card, bestC [3]Card) {
	score = math.MaxUint16
	if len(cl) != 4 || len(community) != 5 {
		return
	}

	for _, holeP := range omahaHolePerms {
		for _, commP := range omahaCommunityPerms {
			var eval = evalFiveFast(
				cl[holeP[0]],
				cl[holeP[1]],
				community[commP[0]],
				community[commP[1]],
				community[commP[2]],
			)
			if eval < score {
				bestH[0] = cl[holeP[0]]
				bestH[1] = cl[holeP[1]]
				bestC[0] = community[commP[0]]
				bestC[1] = community[commP[1]]
				bestC[2] = community[commP[2]]
				score = eval
			}
		}
	}
	return
}
