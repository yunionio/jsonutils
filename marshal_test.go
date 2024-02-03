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
				Name   string
				Age    int
				Weight float32
				Map    map[string]string
				Slice  []string
			}{
				Name:   "",
				Age:    0,
				Weight: 0.0,
			},
			want:    `{"age":0,"weight":0}`,
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

type SGuest struct {
	Id   string
	Name string

	Networks []*SGuestnetwork
}

type SNetwork struct {
	Id      string
	Name    string
	IpStart string
	IpEnd   string
	MaskLen int

	Guests []*SGuestnetwork
}

type SGuestnetwork struct {
	Ip      string
	Guest   *SGuest
	Network *SNetwork
}

type STopo struct {
	Guests   map[string]*SGuest
	Networks map[string]*SNetwork
}

func TestMarshalLoop(t *testing.T) {
	type SLoopStruct struct {
		Id     int
		Name   string
		Gender *bool
		Self   *SLoopStruct
	}

	male := true
	test := SLoopStruct{
		Id:     123,
		Name:   "John",
		Gender: &male,
	}
	test.Self = &test

	test1 := SLoopStruct{
		Id:     123,
		Name:   "John",
		Gender: &male,
	}

	test2 := SLoopStruct{
		Id:   2222,
		Name: "John2",
	}
	test3 := SLoopStruct{
		Id:   3333,
		Name: "John3",
	}
	test4 := SLoopStruct{
		Id:   4444,
		Name: "John4",
	}
	test2.Self = &test3
	test3.Self = &test4
	test4.Self = &test2

	loops := []*SLoopStruct{
		&test2,
		&test3,
		&test4,
	}

	guest1 := &SGuest{
		Id:   "guest1",
		Name: "vm1",
	}
	guest2 := &SGuest{
		Id:   "guest2",
		Name: "vm2",
	}
	net1 := &SNetwork{
		Id:      "net1",
		Name:    "vnet1",
		IpStart: "192.168.222.1",
		IpEnd:   "192.168.222.255",
		MaskLen: 24,
	}

	guest1net := &SGuestnetwork{
		Ip:      "192.168.222.254",
		Guest:   guest1,
		Network: net1,
	}
	guest2net := &SGuestnetwork{
		Ip:      "192.168.222.253",
		Guest:   guest2,
		Network: net1,
	}

	guest1.Networks = []*SGuestnetwork{
		guest1net,
	}
	guest2.Networks = []*SGuestnetwork{
		guest2net,
	}

	net1.Guests = []*SGuestnetwork{
		guest1net,
		guest2net,
	}

	topo := &STopo{
		Guests: map[string]*SGuest{
			guest1.Id: guest1,
			guest2.Id: guest2,
		},
		Networks: map[string]*SNetwork{
			net1.Id: net1,
		},
	}

	type SCdrom struct {
		Path    string
		ImageId string
	}

	type SDesc struct {
		Cdrom  *SCdrom
		Cdroms []*SCdrom
	}

	cdrom := &SCdrom{
		Path:    "rbd.cloudpods-test/image_caches_ba19c4fd-cb44-4ff0-8724-98a1e8bfb7e7",
		ImageId: "ba19c4fd-cb44-4ff0-8724-98a1e8bfb7e7",
	}

	desc := &SDesc{
		Cdrom:  cdrom,
		Cdroms: []*SCdrom{cdrom},
	}

	cases := []struct {
		in  interface{}
		out interface{}
	}{
		{
			in:  &test,
			out: &SLoopStruct{},
		},
		{
			in:  &test1,
			out: &SLoopStruct{},
		},
		{
			in:  &loops,
			out: &[]*SLoopStruct{},
		},
		{
			in:  topo,
			out: &STopo{},
		},
		{
			in:  desc,
			out: &SDesc{},
		},
	}
	for _, c := range cases {
		got := Marshal(c.in).String()
		t.Logf("marshal got: %s", got)

		json, err := ParseString(got)
		if err != nil {
			t.Errorf("parse json fail %s", err)
		} else {
			t.Logf("json: %s", json.String())
			err := json.Unmarshal(c.out)
			if err != nil {
				t.Errorf("unmarshal %s", err)
			} else {
				got2 := Marshal(c.out).String()
				if got2 != got {
					t.Errorf("got: %s != got2: %s", got, got2)
				}
			}
		}
	}

}

func TestMarshalPointer(t *testing.T) {
	type SCdrom struct {
		Path    string
		ImageId string
	}

	type SDesc struct {
		Cdrom  *SCdrom
		Cdroms []*SCdrom
	}

	type SDescHost struct {
		Cdroms []*SCdrom
	}

	cdrom := &SCdrom{
		Path:    "rbd.cloudpods-test/image_caches_ba19c4fd-cb44-4ff0-8724-98a1e8bfb7e7",
		ImageId: "ba19c4fd-cb44-4ff0-8724-98a1e8bfb7e7",
	}

	desc := &SDesc{
		Cdrom:  cdrom,
		Cdroms: []*SCdrom{cdrom},
	}

	got := Marshal(desc).String()
	t.Logf("got: %s", got)
	json, err := ParseString(got)
	if err != nil {
		t.Errorf("ParseString error %s", err)
	} else {
		newdesc := &SDescHost{}
		err := json.Unmarshal(newdesc)
		if err != nil {
			t.Errorf("Unmarshal error %s", err)
		} else if Marshal(desc.Cdroms[0]).String() != Marshal(newdesc.Cdroms[0]).String() {
			t.Errorf("%s != %s", Marshal(desc.Cdroms[0]).String(), Marshal(newdesc.Cdroms[0]).String())
		}
	}
}

func TestMarshalAllowEmpty(t *testing.T) {
	type arrayStruct struct {
		OrgNodeId []string          `json:"org_node_id,allowempty"`
		OrgNodes  map[string]string `json:"org_nodes,allowempty"`
	}
	input := arrayStruct{
		OrgNodeId: make([]string, 0),
		OrgNodes:  make(map[string]string),
	}
	jsonStr := Marshal(input).String()
	wantJsonStr := `{"org_node_id":[],"org_nodes":{}}`
	if jsonStr != wantJsonStr {
		t.Errorf("marshal want %s got %s", wantJsonStr, jsonStr)
	} else {
		json, err := ParseString(jsonStr)
		if err != nil {
			t.Errorf("parsestring fail %s", err)
		} else {
			t.Logf("parstring output JSON: %s", json.String())
			val := arrayStruct{}
			err := json.Unmarshal(&val)
			if err != nil {
				t.Errorf("Unmarshal fail %s", err)
			} else {
				if !reflect.DeepEqual(val, input) {
					t.Errorf("want %s got %s", Marshal(input), Marshal(val))
				}
			}
		}
	}
}

func TestMarshalOmitEmpty(t *testing.T) {
	type arrayStruct struct {
		OrgNodeId []string          `json:"org_node_id,omitempty"`
		OrgNodes  map[string]string `json:"org_nodes,omitempty"`
	}
	input := arrayStruct{
		OrgNodeId: make([]string, 0),
		OrgNodes:  make(map[string]string),
	}
	jsonStr := Marshal(input).String()
	wantJsonStr := `{}`
	if jsonStr != wantJsonStr {
		t.Errorf("marshal want %s got %s", wantJsonStr, jsonStr)
	} else {
		json, err := ParseString(jsonStr)
		if err != nil {
			t.Errorf("parsestring fail %s", err)
		} else {
			t.Logf("parstring output JSON: %s", json.String())
			val := arrayStruct{}
			err := json.Unmarshal(&val)
			if err != nil {
				t.Errorf("Unmarshal fail %s", err)
			} else {
				if Marshal(input).String() != Marshal(val).String() {
					t.Errorf("want %s got %s", Marshal(input), Marshal(val))
				}
			}
		}
	}
}
