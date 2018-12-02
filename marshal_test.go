package jsonutils

import (
	"testing"
	"time"
	"yunion.io/x/onecloud/pkg/mcclient"
)

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
		Number int    `json:",omitempty"`
		Ignore string `json:"-"`
		Name   string `json:",allowempty"`
	}
	test := testStruct{}
	test.Ignore = "TETS"
	t.Logf("%s", Marshal(test))

	test.Number = 2
	t.Logf("%s", Marshal(test))

	type testStruct2 struct {
		Number int `json:",allowempty,string"`
		Ignore string
		Name   string
		Gender *JSONString
	}

	test2 := testStruct2{}
	test2.Gender = NewString("male")
	t.Logf("%s", Marshal(test2))
}


func TestTokenCredential(t *testing.T) {
	type SEmbededCredential struct {
		mcclient.TokenCredential
	}
	realToken := mcclient.SSimpleToken{
		User: "jackey",
		UserId: "jackey123456",
		Project: "system",
		ProjectId: "system_id",
	}

	valToken := SEmbededCredential{&realToken}

	jsonVal := Marshal(&valToken)

	t.Logf("%s", jsonVal)

	valToken1 := SEmbededCredential{}

	err := jsonVal.Unmarshal(&valToken1)

	if err != nil {
		t.Errorf("unmarshal fail %s", err)
	}

}