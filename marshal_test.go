package jsonutils

import (
	"testing"
	"time"
)

func TestMarshal(t *testing.T) {
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
