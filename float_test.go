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
	"math"
	"testing"
)

func TestParseNonFiniteFloat(t *testing.T) {
	cases := []string{
		`{"a":NaN}`,
		`{"a":nan}`,
		`{"a":Inf}`,
		`{"a":Infinity}`,
		`{"a":+Inf}`,
		`{"a":-Inf}`,
		`{"a":-infinity}`,
		`[NaN]`,
	}
	for _, c := range cases {
		jo, err := Parse([]byte(c))
		if err != nil {
			t.Errorf("Parse(%q): %v", c, err)
			continue
		}
		out := jo.String()
		var v interface{}
		if err := json.Unmarshal([]byte(out), &v); err != nil {
			t.Errorf("Parse(%q) => %s, not a valid json document: %v", c, out, err)
			continue
		}
		if c == `[NaN]` {
			continue
		}
		// 非有限值没有 json 表示，按字符串保留原文
		if _, ok := v.(map[string]interface{})["a"].(string); !ok {
			t.Errorf("Parse(%q) => %s, want a string value", c, out)
		}
	}
}

func TestJSONFloatStringNonFinite(t *testing.T) {
	for _, f := range []float64{math.NaN(), math.Inf(1), math.Inf(-1)} {
		if got := NewFloat64(f).String(); got != "null" {
			t.Errorf("NewFloat64(%v).String() = %q, want %q", f, got, "null")
		}
	}
	// 有限值不受影响
	if got := NewFloat64(1.5).String(); got != "1.5" {
		t.Errorf("NewFloat64(1.5).String() = %q, want %q", got, "1.5")
	}
	if got := NewFloat32(1.5).String(); got != "1.5" {
		t.Errorf("NewFloat32(1.5).String() = %q, want %q", got, "1.5")
	}
}

func TestMarshalNonFiniteFloat(t *testing.T) {
	type S struct {
		A float64 `json:"a"`
		B float32 `json:"b"`
		C float64 `json:"c"`
	}
	jo := Marshal(S{A: math.NaN(), B: float32(math.Inf(1)), C: -1.5})
	out := jo.String()

	var v map[string]interface{}
	if err := json.Unmarshal([]byte(out), &v); err != nil {
		t.Fatalf("Marshal => %s, not a valid json document: %v", out, err)
	}
	if v["a"] != nil {
		t.Errorf("a = %v, want null", v["a"])
	}
	if v["b"] != nil {
		t.Errorf("b = %v, want null", v["b"])
	}
	if v["c"] != -1.5 {
		t.Errorf("c = %v, want -1.5", v["c"])
	}
}

func TestUnmarshalNonFiniteIntoFloat(t *testing.T) {
	type S struct {
		F float64 `json:"f"`
		G float32 `json:"g"`
	}
	cases := []string{
		`{"f":NaN}`,
		`{"f":"NaN"}`,
		`{"f":Infinity}`,
		`{"f":"Inf"}`,
		`{"f":"-Infinity"}`,
		`{"g":"NaN"}`,
		`{"g":Infinity}`,
	}
	for _, c := range cases {
		jo, err := ParseString(c)
		if err != nil {
			continue
		}
		var s S
		if err := jo.Unmarshal(&s); err == nil {
			if math.IsNaN(s.F) || math.IsInf(s.F, 0) ||
				math.IsNaN(float64(s.G)) || math.IsInf(float64(s.G), 0) {
				t.Errorf("Parse(%q): a non finite value reached a float field: %+v", c, s)
			}
		}
	}

	// a finite value is unaffected
	var s S
	jo, err := ParseString(`{"f":1.5,"g":-2.25}`)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if err := jo.Unmarshal(&s); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if s.F != 1.5 || s.G != -2.25 {
		t.Errorf("got %+v, want f=1.5 g=-2.25", s)
	}
}
