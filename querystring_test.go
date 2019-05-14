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
)

func TestQueryString(t *testing.T) {
	target := "a=1&b=test&c=3&c=4"
	dict1, e := ParseQueryString(target)
	if e != nil {
		t.Errorf("Fail to parse %s: %s", target, e)
	}
	dict := NewDict()
	dict.Add(NewString("1"), "a")
	dict.Add(NewString("test"), "b")
	arr := NewArray()
	arr.Add(NewString("3"))
	arr.Add(NewString("4"))
	dict.Add(arr, "c")
	if dict1.QueryString() != target {
		t.Errorf("Fail 2 %s != %s", dict1.QueryString(), target)
	}
	if dict.QueryString() != target {
		t.Errorf("Fail 3 %s != %s", dict.QueryString(), target)
	}
}

func TestQueryBoolean(t *testing.T) {
	json := NewDict()
	json.Add(JSONTrue, "true_bool")
	json.Add(JSONFalse, "false_bool")
	json.Add(NewString("true"), "true_string")
	json.Add(NewString("false"), "false_string")
	json.Add(NewInt(1), "true_number")
	json.Add(NewInt(0), "false_number")

	t.Logf("true_bool %v", QueryBoolean(json, "true_bool", false))
	t.Logf("true_string %v", QueryBoolean(json, "true_string", false))
	t.Logf("true_number %v", QueryBoolean(json, "true_number", false))

	t.Logf("false_bool %v", QueryBoolean(json, "false_bool", false))
	t.Logf("false_string %v", QueryBoolean(json, "false_string", false))
	t.Logf("false_number %v", QueryBoolean(json, "false_number", false))
}
