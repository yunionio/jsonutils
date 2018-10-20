package jsonutils

import (
	"reflect"
	"testing"
	"time"

	"yunion.io/x/pkg/tristate"
)

type TestStruct struct {
	Name   string
	Age    int
	Grade  uint8
	Gender string
	Status string
	Json   JSONObject
	Json2  JSONObject
	Array  []string
	Tri    tristate.TriState
}

func TestJSONDictUnmarshal(t *testing.T) {
	var err error
	ts := TestStruct{Name: "test", Age: 23, Grade: 2, Gender: "Male", Status: "Enabled", Tri: tristate.True}
	t.Logf("%s", Marshal(ts))
	json := NewDict()
	json.Add(NewString("name1"), "name")
	json.Add(NewInt(19), "age")
	json.Add(NewInt(3), "grade")
	json.Add(NewStringArray([]string{"1", "2", "3"}), "array")
	json.Add(JSONFalse, "tri")
	subDict := NewDict()
	subDict.Add(NewString("value"), "key")
	subDict.Add(NewString("value2"), "key2")
	json.Add(subDict, "json")
	subArray := NewArray()
	subArray.Add(NewString("arr1"))
	subArray.Add(NewString("arr2"))
	subArray.Add(NewString("arr3"))
	subArray.Add(NewString("arr4"))
	subArray.Add(NewString("arr5"))
	json.Add(subArray, "json2")
	t.Logf("%s", json.String())
	err = json.Unmarshal(&ts)
	if err != nil {
		t.Errorf("unmarshal struct fail: %s", err)
	} else {
		t.Logf("%s", Marshal(ts))
	}

	val := make(map[string]string)
	err = json.Unmarshal(val)
	if err != nil {
		t.Errorf("unmarshal map fail: %s", err)
	} else {
		t.Logf("%s", Marshal(val))
	}
}

func TestJSONDict_Unmarshal(t *testing.T) {
	type TestStruct struct {
		Id   string
		Name string
		Dict JSONObject // *JSONDict
	}
	jsonDict := NewDict()
	jsonDict.Add(NewString("nameVal"), "name")
	jsonDict.Add(NewString("idVal"), "id")
	subDict, err := ParseString(`{"parent_task_id": "30247a37-0328-4c47-bf5e-796672118923", "__stages": [{"complete_at": "2018-05-24T03:00:43Z", "name": "on_init"}], "__request_context": {"request_id": "5c2bd"}}`)
	if err != nil {
		t.Errorf("Parse json error")
	}
	// subDict := NewDict()
	// subDict.Add(NewString("yes"), "answer")
	// subDict.Add(NewInt(24), "age")
	jsonDict.Add(subDict, "dict")
	t.Logf("%s", jsonDict.String())

	dest := TestStruct{}

	jsonDict.Unmarshal(&dest)

	t.Logf("%s", dest)
	t.Logf("%s", Marshal(dest).String())

}

func TestTime(t *testing.T) {
	type TimeStruct struct {
		EndTime time.Time
	}
	jsonDict := NewDict()
	jsonDict.Add(NewString(""), "end_time")
	t.Logf("json: %s", jsonDict.String())
	ts := TimeStruct{}
	err := jsonDict.Unmarshal(&ts)
	if err != nil {
		t.Errorf("unmarshal timestruct error %s", err)
	} else {
		t.Logf("unmarshal result %s", ts)
	}
}

func TestMarshalPtr(t *testing.T) {
	type SPtrs struct {
		Bool   *bool
		Int    *int
		Float  *float64
		String *string
		Struct *struct{ Hmm int }
		Array  *[9]int
		Slice  *[]int
		Map    *map[string]int
	}
	// marshal nils
	ptrsNil := &SPtrs{}
	jsonNil := Marshal(ptrsNil)
	jsonStrNil := jsonNil.String()
	if jsonStrNil != "{}" {
		t.Errorf("Should omit nil values, got %s", jsonStrNil)
	}

	// parse null JSON values
	jsonStrNil2 := `
		{
			bool:    null,
			int:     null,
			float:   null,
			string:  null,
			struct:  null,
			array:   null,
			slice:   null,
			map:     null
		}
	`
	jsonObjNil, err := ParseString(jsonStrNil2)
	if err != nil {
		t.Errorf("parse json string error: %v", err)
	}
	jsonDictNil := jsonObjNil.(*JSONDict)
	if numFields := reflect.TypeOf(SPtrs{}).NumField(); jsonDictNil.Length() != numFields {
		t.Errorf("num fields want %d, got %d", numFields, jsonDictNil.Length())
	}

	// make nonNil
	vBool := true
	vInt := 99
	vFloat := 99.9
	vString := "9999"
	vStruct := struct{ Hmm int }{99999}
	vArray := [9]int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	vSlice := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	vMap := map[string]int{"999999": 1234567}
	ptrsNonNil := &SPtrs{
		Bool:   &vBool,
		Int:    &vInt,
		Float:  &vFloat,
		String: &vString,
		Struct: &vStruct,
		Array:  &vArray,
		Slice:  &vSlice,
		Map:    &vMap,
	}
	jsonStrNonNil := Marshal(ptrsNonNil).String()

	// unmarshal nils to non nils should perform override, partial if the source is not FULL
	jsonObjNil.Unmarshal(ptrsNonNil)
	jsonObj2 := Marshal(ptrsNonNil)
	jsonObj2Str := jsonObj2.String()
	if jsonObj2Str != "{}" {
		t.Errorf("unmarshal result should be {}, got %s", jsonObj2Str)
	}

	// unmarshal non nil str will restore correctly
	{
		jsonObj, err := ParseString(jsonStrNonNil)
		if err != nil {
			t.Errorf("parse error: %s", err)
		}
		ptrs := &SPtrs{}
		err = jsonObj.Unmarshal(ptrs)
		if err != nil {
			t.Errorf("unmarshal error: %s", err)
		}
		jsonStrAgain := Marshal(ptrs).String()
		if jsonStrAgain != jsonStrNonNil {
			t.Errorf("reverse failed: want %s, got %s", jsonStrNonNil, jsonStrAgain)
		}
	}
}

func TestUnmarshalNonNilPtr(t *testing.T) {
	t.Run("non-nil short-cap slice", func(t *testing.T) {
		s := `[43]`
		j, _ := ParseString(s)
		v := []int{}
		vp := &v
		err := j.Unmarshal(&vp)
		if err != nil {
			t.Errorf("expect no error, got %s", err)
			return
		}
		if len(v) != 1 {
			t.Errorf("expect length 0, got %d", len(v))
			return
		}
		if v[0] != 43 {
			t.Errorf("expect [43], got %#v", v)
			return
		}
	})
	t.Run("non-nil over-cap slice", func(t *testing.T) {
		s := `[43]`
		j, _ := ParseString(s)
		v := []int{1, 2}
		vp := &v
		err := j.Unmarshal(&vp)
		if err != nil {
			t.Errorf("expect no error, got %s", err)
			return
		}
		if len(v) != 1 {
			t.Errorf("expect length 0, got %d", len(v))
			return
		}
		if v[0] != 43 {
			t.Errorf("expect [43], got %#v", v)
			return
		}
	})
	t.Run("non-nil map", func(t *testing.T) {
		s := `{"a": "b"}`
		j, _ := ParseString(s)
		v := struct {
			A string
		}{}
		vp := &v
		err := j.Unmarshal(&vp)
		if err != nil {
			t.Errorf("expect no error, got %s", err)
			return
		}
		if v.A != "b" {
			t.Errorf("expect v.A == \"b\", got %#v", v)
			return
		}
	})
}

func TestJSONArrayUnmarshal(t *testing.T) {
	s := `[{"conf":{"cachedbadbbu":false,"conf":"none","count":0,"direct":false,"ra":false,"range":[],"size":[],"strip":0,"type":"hybrid","wt":false},"disks":[{"adapter":0,"driver":"Linux","enclousure":0,"index":0,"max_strip_size":0,"min_strip_size":0,"rotate":true,"size":100000,"slot":0}],"size":100000}]`
	jsonArr, err := ParseString(s)
	if err != nil {
		t.Errorf("parse json error")
	}

	dest := JSONArray{}
	jsonArr.Unmarshal(&dest)
	t.Logf("%s", dest)
	if Marshal(dest).String() != s {
		t.Errorf("TestJSONArrayUnmarshal errors")
	}
}

func TestUnmarshalCurrency(t *testing.T) {
	type SAccountBalance struct {
		USBalance     float64
		GermanBalance float32
	}
	jsonStr := `{"us_balance":"3,118.54", "german_balance":"3.490.000,89"}`
	json, err := ParseString(jsonStr)
	if err != nil {
		t.Errorf("parse %s error %s", jsonStr, err)
		return
	}
	balance := &SAccountBalance{}
	err = json.Unmarshal(&balance)
	if err != nil {
		t.Errorf("unmarshal %s fail %s", jsonStr, err)
		return
	}
}
