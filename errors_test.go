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

package jsonutils // import "yunion.io/x/jsonutils"

import "testing"

func TestNewJSONError(t *testing.T) {
	cases := []struct {
		str []byte
		pos int
	}{
		{
			str: []byte(`{"a":1,"b":2}`),
			pos: 0,
		},
		{
			str: []byte(`{"a":1,"b":2}`),
			pos: 1,
		},
		{
			str: []byte(`{"a":1,"b":2}`),
			pos: 10,
		},
		{
			str: []byte(`{"a":1,"b":2}`),
			pos: 11,
		},
		{
			str: []byte(`{"a":1,"b":2}`),
			pos: 12,
		},
		{
			str: []byte(`{"a":1,"b":2}`),
			pos: 5,
		},
		{
			str: []byte(`{"a":1,"b":2}`),
			pos: 30,
		},
		{
			str: []byte(``),
			pos: 0,
		},
		{
			str: []byte(``),
			pos: 1,
		},

		{
			str: []byte(` `),
			pos: 1,
		},
		{
			str: []byte(` `),
			pos: 2,
		},
	}
	for _, c := range cases {
		err := NewJSONError(c.str, c.pos, "expected a number")
		t.Logf("error: %v", err)
	}
}
