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

// backslash is the escape character, spelled out to keep the test cases readable
var backslash = string(rune(92))

func TestQuoteString(t *testing.T) {
	cases := []string{
		"a\xc3\xa9b", // valid multi byte sequence
		"中文",
		"a\x00b",
		"a\tb",
		`a"b`,
		`a\b`,
	}
	for _, c := range cases {
		quoted := quoteString(c)
		var got string
		if err := json.Unmarshal([]byte(quoted), &got); err != nil {
			t.Errorf("quoteString(%q) = %s, not a valid json string: %v", c, quoted, err)
			continue
		}
		if c != got {
			t.Errorf("quoteString(%q) = %s, the value changed", c, quoted)
		}
		if !utf8.ValidString(got) {
			t.Errorf("quoteString(%q) = %s, the result is not valid utf-8", c, quoted)
		}
	}
}

func TestPrettyStringKeys(t *testing.T) {
	cases := []string{
		`{"a\":1,\"x":2}`,
		`{"a\nb":1}`,
		`{"中文":"a"}`,
		`{"a":1,"b":[1,2,{"c":"d"}]}`,
		`{"":{"":""}}`,
	}
	for _, c := range cases {
		jo, err := ParseString(c)
		if err != nil {
			t.Fatalf("parse %s: %v", c, err)
		}
		pretty := jo.PrettyString()

		back, err := ParseString(pretty)
		if err != nil {
			t.Errorf("Parse(PrettyString(%s)): %v\n%s", c, err, pretty)
			continue
		}
		if !jo.Equals(back) {
			t.Errorf("PrettyString round trip of %s gave %s\n%s", c, back, pretty)
		}

		// 美化输出必须是合法的 json 文档
		var v interface{}
		if err := json.Unmarshal([]byte(pretty), &v); err != nil {
			t.Errorf("PrettyString(%s) is not a valid json document: %v\n%s", c, err, pretty)
		}
	}
}

func TestParseSurrogatePair(t *testing.T) {
	cases := []string{
		`"` + backslash + `ud83d` + backslash + `ude00"`,   // a pair
		`"a` + backslash + `ud83d` + backslash + `ude00b"`, // a pair inside text
		`"` + backslash + `ud800"`,                         // unpaired high surrogate
		`"` + backslash + `udfff"`,                         // unpaired low surrogate
		`"` + backslash + `ud83dx"`,                        // high surrogate then text
		`"` + backslash + `ud83d` + backslash + `u0041"`,   // high surrogate then a bmp escape
		`"` + backslash + `u0041` + backslash + `ude00"`,   // bmp escape then a low surrogate
	}
	for _, c := range cases {
		jo, err := ParseString(c)
		if err != nil {
			t.Errorf("Parse(%s): %v", c, err)
			continue
		}
		got, _ := jo.GetString()

		var std string
		if err := json.Unmarshal([]byte(c), &std); err != nil {
			t.Errorf("json.Unmarshal(%s): %v", c, err)
			continue
		}
		if got != std {
			t.Errorf("Parse(%s) = % x, encoding/json = % x", c, got, std)
		}
		// 重新解析输出应当得到同一个值
		back, err := ParseString(jo.String())
		if err != nil {
			t.Errorf("ParseString(%s): %v", jo.String(), err)
			continue
		}
		if !jo.Equals(back) {
			t.Errorf("Parse(%s) is not stable: %s", c, jo.String())
		}
	}
}
