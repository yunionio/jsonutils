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
