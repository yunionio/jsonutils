package jsonutils

import "testing"

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
