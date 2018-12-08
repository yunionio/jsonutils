package jsonutils

import (
	"fmt"
	"testing"
)

func TestGetQueryStringArray(t *testing.T) {
	for _, queryString := range []string{"search=zone", "search=abc&search=123&search=b", "search.0=123&search.1=abc"} {
		jsonQuery, err := ParseQueryString(queryString)
		if err != nil {
			t.Errorf("Fail to parse query string")
		} else {
			t.Logf("query string: %s", jsonQuery.String())
		}
		search := GetQueryStringArray(jsonQuery, "search")
		t.Logf("%s", search)
	}
}

func TestGetArrayOfPrefix(t *testing.T) {
	json := NewDict()
	for i := 0; i < 10; i += 1 {
		json.Add(NewString(fmt.Sprintf("value.%d", i)), fmt.Sprintf("key.%d", i))
	}

	retArray := GetArrayOfPrefix(json, "key")
	if len(retArray) != 10 {
		t.Errorf("fail to getarrayofprefix")
	}
}
