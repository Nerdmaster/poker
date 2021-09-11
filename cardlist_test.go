package poker

import "testing"

func TestParseCards(t *testing.T) {
	var tests = map[string]struct {
		source        string
		expectedCards CardList
		expectError   bool
	}{
		"Simple": {"2h 5c", CardList{NewCard(Deuce, Hearts), NewCard(Five, Clubs)}, false},
		"Six-card straight": {"2d 3d As Ks Qd Jh", CardList{
			NewCard(Deuce, Diamonds),
			NewCard(Three, Diamonds),
			NewCard(Ace, Spades),
			NewCard(King, Spades),
			NewCard(Queen, Diamonds),
			NewCard(Jack, Hearts),
		}, false},
		"Big hand": {"2h 5c 7h Td Qc Kd As 2s", CardList{
			NewCard(Deuce, Hearts),
			NewCard(Five, Clubs),
			NewCard(Seven, Hearts),
			NewCard(Ten, Diamonds),
			NewCard(Queen, Clubs),
			NewCard(King, Diamonds),
			NewCard(Ace, Spades),
			NewCard(Deuce, Spades),
		}, false},
		"Empty string": {"", CardList{}, false},
		"Invalid":      {"2", nil, true},
		"Invalid 2":    {"!fdasfjk324231FDAS", nil, true},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			var cards, err = ParseCards(tc.source)
			if tc.expectError {
				if err == nil {
					t.Errorf("Expected errors, but got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("Expected no error, but got %s", err)
				return
			}

			if len(tc.expectedCards) != len(cards) {
				t.Errorf("Expected %d cards, but got %d", len(tc.expectedCards), len(cards))
				return
			}

			for i, c := range cards {
				if tc.expectedCards[i] != c {
					t.Errorf("Expected card %d to be %s, but got %s", i, tc.expectedCards[i], c)
				}
			}
		})
	}
}
