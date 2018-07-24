package jsonutils

import (
	"testing"
	"time"
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
	}
	jsonDict := NewDict()
	jsonDict.Add(NewString("jsonValue"), "jsonKey")
	testUser := testStruct{Name: "Test user", Age: 23, Gender: "Male", Birthday: time.Now(),
		Position: Occupation{Name: "Director", Position: "Direction"},
		password: "private info", Json: jsonDict, JsonDict: jsonDict}
	t.Logf("%s", Marshal(testUser))
}
