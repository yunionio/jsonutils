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
	"encoding/json"
	"testing"
	"unicode/utf8"
)

func TestQuoteStringRoundTrip(t *testing.T) {
	cases := []string{
		"a\xffb",
		"a\xc3\xa9b", // a valid multi byte sequence
		"中文",
		"a\x00b",
		"a\tb",
		`a"b`,
		`a\b`,
		"a\xed\xa0\x80b", // a surrogate encoded as utf-8
		"\xe4\xb8",       // a truncated sequence
		"\x80\x81\xfe\xff",
	}
	for _, c := range cases {
		quoted := quoteString(c)

		jo, err := ParseString(quoted)
		if err != nil {
			t.Errorf("ParseString(quoteString(%q) = %s): %v", c, quoted, err)
			continue
		}
		got, err := jo.GetString()
		if err != nil {
			t.Errorf("GetString(%s): %v", quoted, err)
			continue
		}
		if got != c {
			t.Errorf("quoteString(%q) = %s, parsed back as %q", c, quoted, got)
		}
		if !utf8.ValidString(quoted) {
			t.Errorf("quoteString(%q) = %s, the result is not valid utf-8", c, quoted)
		}
	}
}

func TestQuoteStringKeepsStandardUTF8(t *testing.T) {
	// a value that is valid utf-8 must not pull in the \x extension,
	// the output stays a standard json string
	cases := []string{"", "abc", "a中b", "a\tb", `a"b`, `a\b`, "a\x00b"}
	for _, c := range cases {
		quoted := quoteString(c)
		var got string
		if err := json.Unmarshal([]byte(quoted), &got); err != nil {
			t.Errorf("quoteString(%q) = %s, not a standard json string: %v", c, quoted, err)
			continue
		}
		if got != c {
			t.Errorf("quoteString(%q) = %s, encoding/json read %q", c, quoted, got)
		}
	}
}

func TestMarshalInvalidUTF8RoundTrip(t *testing.T) {
	v := struct {
		CommonName string
	}{
		CommonName: string([]byte{'y', 'u', 'n', 'i', 'o', 'n', 129, 10}),
	}
	jd := Marshal(v)
	back, err := ParseString(jd.String())
	if err != nil {
		t.Fatalf("ParseString(%s): %v", jd, err)
	}
	if !back.Equals(jd) {
		t.Errorf("round trip of % x gave %s", v.CommonName, back)
	}
}
