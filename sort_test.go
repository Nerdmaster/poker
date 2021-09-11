package poker

import "testing"

func TestSortGroups(t *testing.T) {
	var tests = map[string]struct {
		input    string
		expected string
	}{
		"Seven-card ace high":        {"2d 3d As Ks Jc 7h 5d", "As Ks Jc 7h 5d 3d 2d"},
		"Seven-card pair":            {"2d 3d Js Ac Jc 7h 5d", "Jc Js Ac 7h 5d 3d 2d"},
		"Seven-card two pair":        {"2d 3d As Ac Jc Jd 5d", "Ac As Jc Jd 5d 3d 2d"},
		"Seven-card three of a kind": {"2c 3d As Ac Ad Jd 5d", "Ac Ad As Jd 5d 3d 2c"},
		"Seven-card straight":        {"2d 3d As Ks Qd Jh Td", "As Ks Qd Jh Td 3d 2d"},
		"Seven-card flush":           {"2d 3d Ts 7s 4s 3s 2s", "3d 3s 2d 2s Ts 7s 4s"},
		"Seven-card full house":      {"7d 3d 4s 2c 4d 2s 2h", "2c 2h 2s 4d 4s 7d 3d"},
		"Crazy":                      {"4d 7c 4c 7d 4s 4h 8d 8c", "4c 4d 4h 4s 8c 8d 7c 7d"},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			var cards, err = ParseCards(tc.input)
			if err != nil {
				t.Errorf("Expected valid CardList using %q, got error: %s", tc.input, err)
				return
			}

			cards.SortGroups()
			var got = cards.String()
			if got != tc.expected {
				t.Errorf("Expected %q to sort to %q, but got %q", tc.input, tc.expected, got)
			}
		})
	}
}

func TestSortAceLow(t *testing.T) {
	var tests = map[string]struct {
		input    string
		expected string
	}{
		"Low straight":                 {"2d 5d Ad 3d 4d", "5d 4d 3d 2d Ad"},
		"Another low because it broke": {"Ac Kd 4h 5d 3d Kc 2c", "Kc Kd 5d 4h 3d 2c Ac"},
		"High straight":                {"Ks Jh Qd As Td", "Ks Qd Jh Td As"},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			var cards, err = ParseCards(tc.input)
			if err != nil {
				t.Errorf("Expected valid CardList using %q, got error: %s", tc.input, err)
				return
			}

			cards.SortAceLow()
			var got = cards.String()
			if got != tc.expected {
				t.Errorf("Expected %q to sort to %q, but got %q", tc.input, tc.expected, got)
			}
		})
	}
}
