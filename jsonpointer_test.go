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

type pointerInnerInt struct {
	B int `json:"b"`
}

type pointerInnerStr struct {
	B string `json:"b"`
}

// 引用出现在节点之后
type pointerNodeFirst struct {
	A *pointerInnerInt `json:"a"`
	C *pointerInnerStr `json:"c"`
}

// 引用出现在节点之前
type pointerRefFirst struct {
	A *pointerInnerStr `json:"a"`
	B *pointerInnerInt `json:"b"`
}

// 引用与节点类型一致
type pointerSameType struct {
	A *pointerInnerInt `json:"a"`
	C *pointerInnerInt `json:"c"`
}

func TestJSONPointerTypeMismatch(t *testing.T) {
	var out pointerNodeFirst
	jo, err := ParseString(`{"a":{"___jnid_":1,"b":5},"c":<1>}`)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if err := jo.Unmarshal(&out); err == nil {
		t.Errorf("expect an error when the reference target has another type")
	}

	var out2 pointerRefFirst
	jo, err = ParseString(`{"a":<1>,"b":{"___jnid_":1,"b":5}}`)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if err := jo.Unmarshal(&out2); err == nil {
		t.Errorf("expect an error when the node has another type than the reference target")
	}
}

func TestJSONPointerSameType(t *testing.T) {
	var out pointerSameType
	jo, err := ParseString(`{"a":{"___jnid_":1,"b":5},"c":<1>}`)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if err := jo.Unmarshal(&out); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if out.A == nil || out.C == nil || out.A.B != 5 || out.C.B != 5 {
		t.Errorf("unexpected result: %+v %+v", out.A, out.C)
	}
}
