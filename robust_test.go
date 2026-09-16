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
	"testing"
)

func TestUnmarshalTargetMustBeWritable(t *testing.T) {
	type S struct {
		A int `json:"a"`
	}
	jo, err := ParseString(`{"a":1}`)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}

	// a struct value is not addressable, there is nowhere to write to
	if err := jo.Unmarshal(S{}); err == nil {
		t.Error("expect an error for a struct value")
	}

	var nilPtr *S
	if err := jo.Unmarshal(nilPtr); err == nil {
		t.Error("expect an error for a nil pointer")
	}

	var nilMap map[string]string
	if err := jo.Unmarshal(nilMap); err == nil {
		t.Error("expect an error for a nil map")
	}

	// a pointer still works
	var s S
	if err := jo.Unmarshal(&s); err != nil {
		t.Fatalf("unmarshal into a pointer: %v", err)
	}
	if s.A != 1 {
		t.Errorf("A = %d, want 1", s.A)
	}

	// a non nil map is a reference type, it still works
	m := make(map[string]string)
	if err := jo.Unmarshal(m); err != nil {
		t.Fatalf("unmarshal into a map: %v", err)
	}
	if m["a"] != "1" {
		t.Errorf("m[a] = %q, want %q", m["a"], "1")
	}
}

func TestUnmarshalDanglingReference(t *testing.T) {
	type Inner struct {
		B int `json:"b"`
	}
	type S struct {
		A *Inner `json:"a"`
		C *Inner `json:"c"`
	}

	jo, err := ParseString(`{"c":<999>}`)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	var s S
	if err := jo.Unmarshal(&s); err == nil {
		t.Errorf("expect an error for a reference with no node, got %+v", s)
	}

	// a resolvable reference is unaffected
	jo, err = ParseString(`{"a":{"___jnid_":1,"b":5},"c":<1>}`)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	var s2 S
	if err := jo.Unmarshal(&s2); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if s2.A == nil || s2.C == nil || s2.A.B != 5 {
		t.Errorf("unexpected result: %+v %+v", s2.A, s2.C)
	}
}

func TestParseQueryStringErrorPropagation(t *testing.T) {
	// a is already a scalar, a.b has nowhere to go
	if _, err := ParseQueryString("a=1&a.b=2"); err == nil {
		t.Error("expect an error when a scalar conflicts with a nested key")
	}

	jo, err := ParseQueryString("a.b=2&a.c=3")
	if err != nil {
		t.Fatalf("ParseQueryString: %v", err)
	}
	if v, err := jo.Int("a", "b"); err != nil || v != 2 {
		t.Errorf("a.b = %d, %v; want 2", v, err)
	}
	if v, err := jo.Int("a", "c"); err != nil || v != 3 {
		t.Errorf("a.c = %d, %v; want 3", v, err)
	}
}
