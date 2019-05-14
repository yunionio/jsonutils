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
