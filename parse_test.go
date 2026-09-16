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
	"strings"
	"testing"
)

func TestParseNoPanic(t *testing.T) {
	cases := []string{
		``,
		`{`,
		`[`,
		`{"a":}`,
		`[}]`,
		`{"a"::}`,
		`{"a":,]`,
		`{"___jnid_":}`,
	}
	for _, c := range cases {
		func() {
			defer func() {
				if r := recover(); r != nil {
					t.Errorf("Parse(%q) panicked: %v", c, r)
				}
			}()
			if res, err := Parse([]byte(c)); err == nil {
				t.Logf("Parse(%q) => %s", c, res)
			}
		}()
	}
}

func TestParseInvalidNodeId(t *testing.T) {
	cases := []string{
		`{"___jnid_":"x"}`,
		`{"___jnid_":null}`,
		`{"___jnid_":1.5}`,
		`{"___jnid_":{}}`,
		`{"___jnid_":[]}`,
		`[{"___jnid_":1},{"___jnid_":1}]`,
		`{"a":{"___jnid_":7},"b":{"___jnid_":7}}`,
	}
	for _, c := range cases {
		if _, err := ParseTrusted([]byte(c)); err == nil {
			t.Errorf("ParseTrusted(%q) expect an error", c)
		}
	}
}

func TestParseNestingDepth(t *testing.T) {
	ok := strings.Repeat("[", maxParseDepth) + strings.Repeat("]", maxParseDepth)
	if _, err := Parse([]byte(ok)); err != nil {
		t.Errorf("Parse depth %d: %v", maxParseDepth, err)
	}

	tooDeep := strings.Repeat("[", maxParseDepth+1) + strings.Repeat("]", maxParseDepth+1)
	if _, err := Parse([]byte(tooDeep)); err == nil {
		t.Errorf("Parse depth %d expect an error", maxParseDepth+1)
	}

	okDict := strings.Repeat(`{"a":`, maxParseDepth) + "1" + strings.Repeat("}", maxParseDepth)
	if _, err := Parse([]byte(okDict)); err != nil {
		t.Errorf("Parse dict depth %d: %v", maxParseDepth, err)
	}

	tooDeepDict := strings.Repeat(`{"a":`, maxParseDepth+1) + "1" + strings.Repeat("}", maxParseDepth+1)
	if _, err := Parse([]byte(tooDeepDict)); err == nil {
		t.Errorf("Parse dict depth %d expect an error", maxParseDepth+1)
	}
}
