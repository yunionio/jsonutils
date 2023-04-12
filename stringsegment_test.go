// Copyright 2019 Yunion
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
