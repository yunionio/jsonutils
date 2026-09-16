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

type nodeRefInner struct {
	B int `json:"b"`
}

type nodeRefTwoFields struct {
	A *nodeRefInner `json:"a"`
	C *nodeRefInner `json:"c"`
}

type nodeRefCycle struct {
	Name string        `json:"name"`
	Next *nodeRefCycle `json:"next,omitempty"`
}

func TestNodeReferenceNotResolved(t *testing.T) {
	jo, err := ParseString(`{"a":{"___jnid_":1,"b":5},"c":<1>}`)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}

	// <1> is kept as a plain value, it is not resolved to the node
	c, err := jo.Get("c")
	if err != nil {
		t.Fatalf("get c: %v", err)
	}
	if str, err := c.GetString(); err != nil || str != "<1>" {
		t.Errorf("c = %v, %v; want the string \"<1>\"", c, err)
	}

	// the reserved key stays an ordinary key of the nested object
	a, err := jo.Get("a")
	if err != nil {
		t.Fatalf("get a: %v", err)
	}
	dict, ok := a.(*JSONDict)
	if !ok {
		t.Fatalf("a is %T, want a dict", a)
	}
	if !dict.Contains("___jnid_") {
		t.Errorf("___jnid_ should be kept as an ordinary key, got %v", dict.SortedKeys())
	}

	// and the two fields must not end up pointing at the same object
	var s nodeRefTwoFields
	if err := jo.Unmarshal(&s); err == nil {
		if s.A != nil && s.C != nil && s.A == s.C {
			t.Errorf("the two fields must not point at the same object")
		}
	}
}

func TestNodeReferenceDoesNotAlias(t *testing.T) {
	// a document can not make two fields share one object
	jo, err := ParseString(`{"a":{"___jnid_":1,"b":5},"c":<1>}`)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	var s nodeRefTwoFields
	if err := jo.Unmarshal(&s); err == nil && s.A != nil && s.A == s.C {
		t.Errorf("unexpected alias: %+v", s)
	}
}

func TestMarshalUnmarshalKeepsCycle(t *testing.T) {
	a := &nodeRefCycle{Name: "a"}
	b := &nodeRefCycle{Name: "b"}
	a.Next = b
	b.Next = a

	// Marshal keeps writing the reference syntax, it is what makes a cyclic
	// object terminate
	jo := Marshal(a)

	var c nodeRefCycle
	if err := jo.Unmarshal(&c); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if c.Name != "a" || c.Next == nil || c.Next.Name != "b" {
		t.Fatalf("unexpected result: %+v", c)
	}
	if c.Next.Next != &c {
		t.Errorf("the cycle does not point back at the root")
	}
}

func TestParseTrustedResolvesReference(t *testing.T) {
	jo, err := ParseTrustedString(`{"a":{"___jnid_":1,"b":5},"c":<1>}`)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	var s nodeRefTwoFields
	if err := jo.Unmarshal(&s); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if s.A == nil || s.C == nil || s.A.B != 5 || s.C.B != 5 {
		t.Fatalf("unexpected result: %+v %+v", s.A, s.C)
	}
	if s.A != s.C {
		t.Errorf("a trusted document resolves the reference to the same object")
	}

	// the reserved key is consumed, as it was before
	a, err := jo.Get("a")
	if err != nil {
		t.Fatalf("get a: %v", err)
	}
	if dict, ok := a.(*JSONDict); ok && dict.Contains("___jnid_") {
		t.Errorf("___jnid_ should be consumed for a trusted document")
	}
}

func TestParseTrustedRoundTrip(t *testing.T) {
	// the whole point of the trusted entry point: a cyclic object survives
	// a trip through its text form
	a := &nodeRefCycle{Name: "a"}
	b := &nodeRefCycle{Name: "b"}
	a.Next = b
	b.Next = a

	text := Marshal(a).String()
	jo, err := ParseTrustedString(text)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	var c nodeRefCycle
	if err := jo.Unmarshal(&c); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if c.Next == nil || c.Next.Next != &c {
		t.Errorf("the cycle was not restored: %+v", c)
	}
}
