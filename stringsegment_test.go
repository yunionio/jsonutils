package jsonutils

import (
	"reflect"
	"testing"
)

func TestStringSegments(t *testing.T) {
	cases := []struct {
		input []string
		want  [][]sTextNumber
	}{
		{
			input: []string{"provider"},
			want: [][]sTextNumber{
				{
					{
						text:     "provider",
						isNumber: false,
					},
				},
			},
		},
		{
			input: []string{"provider.0", "provider.1"},
			want: [][]sTextNumber{
				{
					{
						text:     "provider",
						isNumber: false,
					},
					{
						number:   0,
						isNumber: true,
					},
				},
				{
					{
						text:     "provider",
						isNumber: false,
					},
					{
						number:   1,
						isNumber: true,
					},
				},
			},
		},
		{
			input: []string{"provider.0.name"},
			want: [][]sTextNumber{
				{
					{
						text:     "provider",
						isNumber: false,
					},
					{
						number:   0,
						isNumber: true,
					},
					{
						text:     "name",
						isNumber: false,
					},
				},
			},
		},
		{
			input: []string{"provider.0.1"},
			want: [][]sTextNumber{
				{
					{
						text:     "provider",
						isNumber: false,
					},
					{
						number:   0,
						isNumber: true,
					},
					{
						number:   1,
						isNumber: true,
					},
				},
			},
		},
	}
	for _, c := range cases {
		got := strings2stringSegments(c.input)
		if !reflect.DeepEqual(got, sStringSegments(c.want)) {
			t.Errorf("input %s got %#v != want %#v", c.input, got, c.want)
		}
		got2 := stringSegments2Strings(got)
		if !reflect.DeepEqual(c.input, got2) {
			t.Errorf("input %s != got %#v", c.input, got2)
		}
	}
}
