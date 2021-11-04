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
	"time"

	"yunion.io/x/pkg/tristate"
)

func TestMarshal(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		type I interface{}
		t.Run("nil-self", func(t *testing.T) {
			var i I
			v := Marshal(i)
			if v != JSONNull {
				t.Errorf("expect json null, got %s", v)
			}
		})
		t.Run("nil with type", func(t *testing.T) {
			var i = func() I {
				var i *int
				return i
			}()
			v := Marshal(i)
			if v != JSONNull {
				t.Errorf("expect json null, got %s", v)
			}
		})
	})
	t.Run("anonymous interface", func(t *testing.T) {
		type I interface{}
		type M struct {
			N string
		}
		type S struct {
			I
		}
		s := S{I: M{N: "hello"}}
		want, _ := ParseString(`{"n": "hello"}`)
		got := Marshal(s)
		if !got.Equals(want) {
			t.Errorf("got %s, want %s", got, want)
		}
	})
	t.Run("anonymous struct", func(t *testing.T) {
		type M struct {
			N string
		}
		type S struct {
			M
		}
		s := S{M{N: "hello"}}
		want, _ := ParseString(`{"n": "hello"}`)
		got := Marshal(s)
		if !got.Equals(want) {
			t.Errorf("got %s, want %s", got, want)
		}
	})
	t.Run("anonymous struct ptr", func(t *testing.T) {
		type M struct {
			N string
		}
		type S struct {
			*M
		}
		s := S{&M{N: "hello"}}
		want, _ := ParseString(`{"n": "hello"}`)
		got := Marshal(s)
		if !got.Equals(want) {
			t.Errorf("got %s, want %s", got, want)
		}
	})
	t.Run("all", func(t *testing.T) {
		type EmbedBasicString string
		type EmbedStruct struct {
			EmbedStructM string
		}
		type EmbedStructPtr struct {
			EmbedStructPtrM string
		}
		type EmbedInterfaceStruct interface{}
		type EmbedInterfaceBasic interface{}
		type S struct {
			BasicString    string
			BasicStringPtr *string

			string // ignored for field name being not exported
			EmbedBasicString
			EmbedStruct
			*EmbedStructPtr
			EmbedInterfaceStruct
			EmbedInterfaceBasic
		}
		str := "BasicStringPtr"
		s := S{
			BasicString:    "BasicString",
			BasicStringPtr: &str,

			string:           "string",
			EmbedBasicString: "EmbedBasicString",
			EmbedStruct: EmbedStruct{
				EmbedStructM: "EmbedStructM",
			},
			EmbedStructPtr: &EmbedStructPtr{
				EmbedStructPtrM: "EmbedStructPtrM",
			},
			EmbedInterfaceStruct: EmbedInterfaceStruct(struct {
				EmbedInterfaceStructM string
			}{"EmbedInterfaceStructM"}),
			EmbedInterfaceBasic: EmbedInterfaceBasic("EmbedInterfaceBasic"),
		}
		want, _ := ParseString(`{
			"basic_string":"BasicString",
			"basic_string_ptr":"BasicStringPtr",

			"embed_basic_string":"EmbedBasicString",
			"embed_struct_m":"EmbedStructM",
			"embed_struct_ptr_m":"EmbedStructPtrM",
			"embed_interface_struct_m":"EmbedInterfaceStructM",
			"embed_interface_basic":"EmbedInterfaceBasic"
		}`)
		got := Marshal(s)
		if !got.Equals(want) {
			t.Errorf("got %s, want %s", got, want)
		}
	})
}

func TestJSONMarshal(t *testing.T) {
	t.Logf("%s", Marshal("string"))
	t.Logf("%s", Marshal(true))
	mapStringTest := map[string]string{
		"Name":     "Testtest",
		"Gender":   "Male",
		"Birthday": "2011-01-01",
		"Talent":   "",
	}
	t.Logf("%s", Marshal(mapStringTest))

	mapTest := map[string]interface{}{
		"Name":     "Test User",
		"Age":      23,
		"Gender":   "Male",
		"Birthday": time.Now(),
		"emptu":    "",
	}
	t.Logf("%s", Marshal(mapTest))
	arrTest := []interface{}{
		"Test user",
		23,
		"Male",
		time.Now(),
	}
	t.Logf("%s", Marshal(arrTest))

	type GenderType string

	type Occupation struct {
		Name     string
		Position string
	}
	type testStruct struct {
		Name       string
		Age        int
		Gender     GenderType `json:"__gender__"`
		Birthday   time.Time
		Position   Occupation
		password   string
		Json       JSONObject
		JsonDict   *JSONDict
		EmptyJson  *JSONArray
		EmptyJson2 JSONObject
		Mail       string `json:"mail,omitempty"`
		Ignore     int    `json:"-"`
	}
	jsonDict := NewDict()
	jsonDict.Add(NewString("jsonValue"), "jsonKey")
	testUser := testStruct{Name: "Test user", Age: 23, Gender: "Male", Birthday: time.Now(),
		Position: Occupation{Name: "Director", Position: "Direction"},
		password: "private info", Json: jsonDict, JsonDict: jsonDict,
		Mail:   "test@yunion.io",
		Ignore: 1,
	}
	t.Logf("%s", Marshal(testUser))
}

func TestJSONMarshalTag(t *testing.T) {
	type testStruct struct {
		Number int    `json:",omitzero"`
		Ignore string `json:"-"`
		Name   string `json:"OS:Name:Test,allowempty"`
	}
	test := testStruct{}
	test.Ignore = "TETS"
	j1 := Marshal(test)
	t.Logf("%s", j1)

	if j1.Contains("number") {
		t.Fatalf("omitzero field presents")
	}
	if !j1.Contains("OS:Name:Test") {
		t.Fatalf("allowempty field should present")
	}

	test.Number = 2
	j2 := Marshal(test)
	t.Logf("%s", j2)
	if !j2.Contains("number") {
		t.Fatalf("Non-zero omitzero field shoudl present")
	}
	if !j2.Contains("OS:Name:Test") {
		t.Fatalf("allowempty field should present")
	}

	type testStruct2 struct {
		Number      int `json:"Number,allowzero,string"`
		EmptyNumber int `json:"EmptyNumber,omitzero"`
		Ignore      string
		Name        string
		Gender      *JSONString
	}

	test2 := testStruct2{}
	test2.Gender = NewString("male")
	j := Marshal(test2)
	t.Logf("%s", Marshal(test2))
	if !j.Contains("Number") {
		t.Fatalf("allowzero field should present")
	}
	numJson, _ := j.Get("Number")
	if _, ok := numJson.(*JSONString); !ok {
		t.Fatalf("forcestring field should be JSONString")
	}
	if j.Contains("EmptyNumber") {
		t.Fatalf("omitzero field should not present")
	}
}

func TestMarshalDeprecatedBy(t *testing.T) {
	type Struct0 struct {
		Name string `json:"name"`

		Cloudregion string `json:"cloudregion"`
		// Deprecated
		CloudregionId string `json:"cloudregion_id" yunion-deprecated-by:"cloudregion"`
		// Deprecated
		RegionId string `json:"region_id" yunion-deprecated-by:"cloudregion_id"`
		// Deprecated
		Region string `json:"region" yunion-deprecated-by:"region_id"`

		Loop0 string `json:"loop0" yunion-deprecated-by:"loop1"`
		Loop1 string `json:"loop1" yunion-deprecated-by:"loop2"`
		Loop2 string `json:"loop2" yunion-deprecated-by:"loop0"`
	}
	s := Struct0{}
	s.Cloudregion = "region0"
	jsonS := Marshal(s)
	s2 := Struct0{}
	err := jsonS.Unmarshal(&s2)
	if err != nil {
		t.Fatalf("unmarshal s3 should success %s", err)
	}
	if s2.CloudregionId != s.Cloudregion {
		t.Errorf("expect s2.CloudregionId(%s) == s.Cloudregion(%s)", s2, s)
	}
	if s2.RegionId != s.Cloudregion {
		t.Errorf("expect s2.RegionId(%s) == s.Cloudregion(%s)", s2, s)
	}
	if s2.Region != s.Cloudregion {
		t.Errorf("expect s2.Region(%s) == s.Cloudregion(%s)", s2, s)
	}
}

func TestMarshalAll(t *testing.T) {
	cases := []struct {
		in      interface{}
		want    string
		wantAll string
	}{
		{
			in:      true,
			want:    "true",
			wantAll: "true",
		},
		{
			in:      tristate.None,
			want:    "null",
			wantAll: "null",
		},
		{
			in: struct {
				Name   string            `json:",omitempty"`
				Age    int               `json:",omitzero"`
				Weight float32           `json:",omitzero"`
				Map    map[string]string `json:",omitempty"`
				Slice  []string          `json:",omitempty"`
			}{
				Name:   "",
				Age:    0,
				Weight: 0.0,
			},
			want:    "{}",
			wantAll: `{"age":0,"map":null,"name":"","slice":null,"weight":0}`,
		},
		{
			in: struct {
				Name   string            `json:",omitempty"`
				Age    int               `json:",omitzero"`
				Weight float32           `json:",omitzero"`
				Map    map[string]string `json:",allowempty"`
				Slice  []string          `json:",allowempty"`
			}{
				Name:   "",
				Age:    0,
				Weight: 0.0,
			},
			want:    "{}",
			wantAll: `{"age":0,"map":null,"name":"","slice":null,"weight":0}`,
		},
		{
			in: struct {
				Name   string            `json:",omitempty"`
				Age    int               `json:",omitzero"`
				Weight float32           `json:",omitzero"`
				Map    map[string]string `json:",omitempty"`
				Slice  []string          `json:",omitempty"`
			}{
				Name:   "",
				Age:    0,
				Weight: 0.0,
				Map:    make(map[string]string),
				Slice:  make([]string, 0),
			},
			want:    "{}",
			wantAll: `{"age":0,"map":{},"name":"","slice":[],"weight":0}`,
		},
		{
			in: struct {
				Name   string  `json:",allowempty"`
				Age    int     `json:",allowzero"`
				Weight float32 `json:",allowzero"`
			}{
				Name:   "",
				Age:    0,
				Weight: 0.0,
			},
			want:    `{"age":0,"name":"","weight":0}`,
			wantAll: `{"age":0,"name":"","weight":0}`,
		},
		{
			in: map[string]string{
				"name": "",
			},
			want:    `{"name":""}`,
			wantAll: `{"name":""}`,
		},
		{
			in: []string{
				"",
			},
			want:    `[""]`,
			wantAll: `[""]`,
		},
	}
	for _, c := range cases {
		got := Marshal(c.in)
		if got.String() != c.want {
			t.Errorf("want %s got %s", c.want, got)
		}
		gotAll := MarshalAll(c.in)
		if gotAll.String() != c.wantAll {
			t.Errorf("wantAll %s gotAll %s", c.wantAll, gotAll)
		}
	}
}
