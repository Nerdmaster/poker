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

func Evaluate(c []Card) uint16 {
	if len(c) == 5 {
		return evalFiveFast(c[0], c[1], c[2], c[3], c[4])
	}
	if len(c) < 5 || len(c) > 8 {
		return math.MaxUint16
	}

	return evalMore(c)
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
