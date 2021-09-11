package poker

import "strings"

// CardList is a simple alias for a list of ... wait for it... CARDS!
type CardList []Card

func (cards CardList) String() string {
	var slist []string
	for _, c := range cards {
		slist = append(slist, c.String())
	}

	return strings.Join(slist, " ")
}
