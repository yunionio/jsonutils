package jsonutils

import (
	"testing"
	"time"
)

type TestStruct struct {
	Name   string
	Age    int
	Gender string
	Status string
	Json   JSONObject
	Json2  JSONObject
	Array  []string
}

func TestJSONDictUnmarshal(t *testing.T) {
	var err error
	ts := TestStruct{Name: "test", Age: 23, Gender: "Male", Status: "Enabled"}
	t.Logf("%s", Marshal(ts))
	json := NewDict()
	json.Add(NewString("name1"), "name")
	json.Add(NewInt(19), "age")
	json.Add(NewStringArray([]string{"1", "2", "3"}), "array")
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
