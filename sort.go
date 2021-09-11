package poker

import "sort"

// SortAceHigh puts the most valuable cards on the "left", and considers an Ace
// to be the most valuable card, which is true in all non-grouped cases except
// certain straights
func (c CardList) SortAceHigh() {
	sort.Slice(c, func(i, j int) bool {
		var r1 = c[i].Rank()
		var r2 = c[j].Rank()

		if r1 == r2 {
			return c[i] > c[j]
		}
		return r1 > r2
	})
}

// SortAceLow considers an Ace to be the least valuable card, which is
// necessary strictly for low straights
func (c CardList) SortAceLow() {
	sort.Slice(c, func(i, j int) bool {
		var rI = c[i].Rank()
		var rJ = c[j].Rank()

		if rI == rJ {
			return c[i] > c[j]
		}

		var xI, xJ = int(rI), int(rJ)
		if rI == Ace {
			xI = 0
		}
		if rJ == Ace {
			xJ = 0
		}

		return xI > xJ
	})
}

type group struct {
	r    CardRank
	list CardList
}

// SortGroups breaks the cards into groups by rank, sorts each group, then
// reassembles the groupings in order of number in a group, then group rank.
// This kind of sort makes sense when displaying any of the "grouped by rank"
// hands, such as a pair, a full house, etc.
func (c CardList) SortGroups() {
	var groups []*group
	var lookup = make(map[CardRank]*group)
	for _, card := range c {
		var g = lookup[card.Rank()]
		if g == nil {
			g = &group{r: card.Rank()}
			groups = append(groups, g)
			lookup[card.Rank()] = g
		}
		g.list = append(g.list, card)
	}

	// Sort groups so those with the most cards are on the "left", then by card
	// value when groups are of equal size
	sort.Slice(groups, func(i, j int) bool {
		var li, lj = len(groups[i].list), len(groups[j].list)
		if li == lj {
			return groups[i].r > groups[j].r
		}
		return li > lj
	})

	// Sort cards within each group for consistency, while also aggregating
	// "groups" of one
	for _, g := range groups {
		g.list.SortAceHigh()
	}

	// Rebuild the entire list of cards
	var idx int
	for _, g := range groups {
		for _, card := range g.list {
			c[idx] = card
			idx++
		}
	}
}
