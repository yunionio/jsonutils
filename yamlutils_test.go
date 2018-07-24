package jsonutils

import (
	"testing"
)

func TestIndentLines(t *testing.T) {
	cases := []struct {
		in       []string
		in_array bool
		out      []string
	}{
		{[]string{"abc", "def"}, true, []string{"- abc", "  def"}},
		{[]string{"abc", "def"}, false, []string{"  abc", "  def"}},
	}
	for _, c := range cases {
		got := indentLines(c.in, c.in_array)
		if got[0] != c.out[0] {
			t.Errorf("IndentLines(%s, %v) = %s != %s", c.in, c.in_array, got, c.out)
		}
	}
}
