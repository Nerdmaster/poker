package poker

type HandRank int

const (
	StraightFlush HandRank = iota + 1
	FourOfAKind
	FullHouse
	Flush
	Straight
	ThreeOfAKind
	TwoPair
	OnePair
	HighCard
)

func GetHandRank(v uint16) HandRank {
	if v > 6185 {
		return HighCard
	}
	if v > 3325 {
		return OnePair
	}
	if v > 2467 {
		return TwoPair
	}
	if v > 1609 {
		return ThreeOfAKind
	}
	if v > 1599 {
		return Straight
	}
	if v > 322 {
		return Flush
	}
	if v > 166 {
		return FullHouse
	}
	if v > 10 {
		return FourOfAKind
	}
	return StraightFlush
}

func (r HandRank) String() string {
	switch r {
	case StraightFlush:
		return "Straight Flush"
	case FourOfAKind:
		return "Four Of A Kind"
	case FullHouse:
		return "Full House"
	case Flush:
		return "Flush"
	case Straight:
		return "Straight"
	case ThreeOfAKind:
		return "Three Of A Kind"
	case TwoPair:
		return "Two Pair"
	case OnePair:
		return "One Pair"
	case HighCard:
		return "High Card"
	}

	return ""
}
