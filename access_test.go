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

func TestNewDict(t *testing.T) {
	dict := NewDict()
	dict.Add(NewString("1"), "a", "b")
	dict2, _ := ParseString("{\"a\": {\"b\": \"1\"}}")
	if dict.String() != dict2.String() {
		t.Errorf("Fail %s != %s", dict, dict2)
	}
	dict = NewDict()
	dict2, _ = ParseString("{}")
	if dict.String() != dict2.String() {
		t.Errorf("Fail %s != %s", dict, dict2)
	}
}

func TestNewArray(t *testing.T) {
	arr := NewArray()
	arr.Add(NewString("1"), NewInt(1), NewFloat64(1.0))
	arr2, _ := ParseString("[\"1\", 1, 1.0]")
	if arr.String() != arr2.String() {
		t.Errorf("Fail %s != %s", arr, arr2)
	}
	arr = NewArray()
	arr2, _ = ParseString("[]")
	if arr.String() != arr2.String() {
		t.Errorf("Fail %s != %s", arr, arr2)
	}
}

func TestNewBool(t *testing.T) {
	type args struct {
		val bool
	}
	tests := []struct {
		name string
		args args
		want *JSONBool
	}{
		{
			name: "New-bool-true",
			args: args{true},
			want: &JSONBool{data: true},
		},
		{
			name: "New-bool-false",
			args: args{false},
			want: &JSONBool{data: false},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewBool(tt.args.val); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewBool() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestJSONDictRemove(t *testing.T) {
	jd := NewDict()
	for k, v := range map[string]JSONObject{
		"Hello": NewString("world"),
		"hello": NewString("world"),
		"HELLO": NewString("world"),
	} {
		jd.Set(k, v)
	}
	var removed bool
	if removed = jd.Remove("HEllo"); removed {
		t.Fatalf("case sensitive remove, want false, got true")
	}
	if removed = jd.Remove("HELLO"); !removed {
		t.Fatalf("case sensitive remove, want true, got false")
	}
	if removed = jd.Remove("HELLO"); removed {
		t.Fatalf("case sensitive remove, want false, got true")
	}
	if removed = jd.RemoveIgnoreCase("hello"); !removed {
		t.Fatalf("case insensitive remove, want true, got false")
	}
	if removed = jd.RemoveIgnoreCase("hello"); removed {
		t.Fatalf("case insensitive false, want true, got true")
	}
}

func TestJSONSpecialChar(t *testing.T) {
	cases := []struct {
		in string
	}{
		{
			in: string([]byte{'y', 'u', 'n', 'i', 'o', 'n', '\n'}),
		},
		{
			in: "中文Engilish\bok\n\t",
		},
		{
			in: string([]byte{'y', 'u', 'n', 'i', 'o', 'n', 10, 10, 10, 10}),
		},
		{
			in: string([]byte{'y', 'u', 'n', 'i', 'o', 'n', 129, 10, 10, 10}),
		},
		{
			in: string([]byte{'y', 'u', 'n', 'i', 'o', 'n', 8, 8, 8, 8}),
		},
		{
			in: "中文 空格；符号。 中文：",
		},
	}
	for _, c := range cases {
		v := struct {
			CommonName string
		}{
			CommonName: c.in,
		}
		jd := Marshal(v)
		output := jd.String()
		t.Log("output: ", output)
		newjd, err := ParseString(output)
		if err != nil {
			t.Errorf("ParseString %s fail: %s", output, err)
		}
		if !newjd.Equals(jd) {
			t.Errorf("newJd %s !=  jd %s", newjd, jd)
		}
		newStr, _ := newjd.GetString("common_name")
		t.Logf("%s %x", newStr, newStr)
	}
}
