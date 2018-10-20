package jsonutils

import "testing"

func TestNormalizeCurrencyString(t *testing.T) {
	cases := []struct {
		in   string
		want string
	}{
		{"3,000.89", "3000.89"},
		{"3,000.", "3000."},
		{"3.000,89", "3000.89"},
		{"3.000,", "3000."},
	}
	for _, c := range cases {
		got := normalizeCurrencyString(c.in)
		if got != c.want {
			t.Errorf("normalize %s got %s want %s", c.in, got, c.want)
		}
	}
}
