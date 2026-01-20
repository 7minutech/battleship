package gamelogic

import "testing"

func TestGetWords(t *testing.T) {

	cases := []struct {
		input    string
		expected []string
	}{
		{"move", []string{"move"}},
		{"move a1", []string{"move", "a1"}},
		{"place b2", []string{"place", "b2"}},
		{"a b c d", []string{"a", "b", "c", "d"}},
		{"", []string{""}},
	}

	for _, c := range cases {
		actual := GetWords(c.input)
		if len(actual) != len(c.expected) {
			t.Errorf("len(%s) = %d; not %d)", c.input, len(actual), len(c.expected))
		}
		for i := range actual {
			if actual[i] != c.expected[i] {
				t.Errorf("%s[%d] = %s; not %s", c.input, i, actual[i], c.expected[i])
			}
		}
	}

}
