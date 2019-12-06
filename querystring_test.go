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
)

func TestQueryString(t *testing.T) {
	type Target struct {
		A string
		B string
		C []string
	}
	target1 := "a=1&b=test&c=3&c=4"
	dict1, e := ParseQueryString(target1)
	if e != nil {
		t.Errorf("Fail to parse %s: %s", target1, e)
	}
	target2 := "a=1&b=test&c.0=3&c.1=4"
	dict2, e := ParseQueryString(target2)
	if e != nil {
		t.Errorf("Fail to parse %s: %s", target2, e)
	}
	s1 := Target{}
	err := dict1.Unmarshal(&s1)
	if err != nil {
		t.Errorf("unmarshal fail %s", err)
	}
	s2 := Target{}
	err = dict2.Unmarshal(&s2)
	if err != nil {
		t.Errorf("unmarshal fail %s", err)
	}
	if !reflect.DeepEqual(s1, s2) {
		t.Errorf("s1 %#v != s2 %#v", s1, s2)
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

func TestQueryString2(t *testing.T) {
	type Person struct {
		Name     string
		Gender   string
		IsLeader bool
		Age      int
		Alias    []string
		Friends  []Person
	}
	val := Person{
		Name:     "Perter",
		Gender:   "Male",
		IsLeader: true,
		Age:      20,
		Alias:    []string{"John", "Smith"},
		Friends: []Person{
			{
				Name:     "Alice",
				Gender:   "Female",
				IsLeader: true,
				Age:      32,
				Alias:    []string{"Emily", "Emma"},
			},
			{
				Name:     "Tom",
				Gender:   "Male",
				IsLeader: false,
				Age:      33,
				Friends: []Person{
					{
						Name:   "Johnson",
						Gender: "Female",
						Alias:  []string{"Google"},
					},
				},
			},
			{
				Name:     "longAlias",
				Gender:   "Female",
				IsLeader: false,
				Age:      44,
				Alias: []string{
					"a10", "a11", "a12", "a13", "a14", "a15", "a16", "a17", "a18", "a19",
					"a0", "a1", "a2", "a3", "a4", "a5", "a6", "a7", "a8", "a9",
					"a20",
				},
			},
		},
	}
	valJson := Marshal(val)
	queryString := valJson.QueryString()
	t.Logf("%s", queryString)
	valJson2, err := ParseQueryString(queryString)
	if err != nil {
		t.Errorf("ParseQueryString fail %s", err)
	} else {
		val2 := Person{}
		err := valJson2.Unmarshal(&val2)
		if err != nil {
			t.Errorf("valJson2.Unmarshal %s", err)
		} else if !reflect.DeepEqual(val, val2) {
			t.Errorf("val %#v != val2 %#v", val, val2)
		}
	}

}
