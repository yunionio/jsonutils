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

func TestDeepCopyBasic(t *testing.T) {
	cases := []struct {
		name string
		in   JSONObject
	}{
		{
			name: "string",
			in:   NewString(`hello world`),
		},
		{
			name: "int",
			in:   NewInt(2048),
		},
		{
			name: "float",
			in:   NewFloat(2048.4096),
		},
		{
			name: "bool true",
			in:   JSONTrue,
		},
		{
			name: "bool false",
			in:   JSONFalse,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			vc := DeepCopy(c.in)
			if !reflect.DeepEqual(vc, c.in) {
				t.Errorf("want %s, got %s", c.in.String(), vc.String())
			}
		})
	}
}

func TestDeepCopyCompound(t *testing.T) {
	cases := []struct {
		name string
		in   string
	}{
		{
			name: "dict",
			in:   `{"a": 0, "b": true, "c": 3.14}`,
		},
		{
			name: "array",
			in:   `["hello world", 0, 3.14, true, false]`,
		},
		{
			name: "dict (2 level)",
			in:   `{"a": 0, "b": true, "c": 3.14, d:{}, e:[]}`,
		},
		{
			name: "array (2 level)",
			in:   `["hello world", 0, 3.14, true, false, {}, []]`,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			obj, err := ParseString(c.in)
			if err != nil {
				t.Fatalf("must parse, got %s", err)
			}
			objC := DeepCopy(obj)
			if !obj.Equals(objC) {
				t.Fatalf("want\n%s\n, got\n%s", obj.String(), objC.String())
			}
			switch v := objC.(type) {
			case *JSONArray:
				v.SetAt(0, NewString("who knows me"))
			case *JSONDict:
				v.Set("id", NewString("nobody"))
			}
			if obj.Equals(objC) {
				t.Fatalf("not copied: %s", obj.String())
			}
		})
	}
}
