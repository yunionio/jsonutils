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

func TestDiff(t *testing.T) {
	a := struct {
		Name string `json:"name"`
		Desc string `json:"description"`
		Male bool   `json:"is_male"`
	}{
		Name: "john",
		Desc: "jonn's json",
		Male: true,
	}
	old := Marshal(a).(*JSONDict)
	new := old.CopyExcludes("name")
	new.Add(NewString("manager"), "title")
	new.Set("is_male", JSONFalse)

	aNoB, aDiffB, aAndB, bNoA := Diff(old, new)

	t.Logf("a: %s b: %s", old, new)
	t.Logf("aNoB: %s aDiffB: %s aAndB: %s bNoA: %s", aNoB, aDiffB, aAndB, bNoA)

	if !aNoB.Contains("name") {
		t.Errorf("a shoud remove name")
	}
	if !bNoA.Contains("title") {
		t.Errorf("b should add title")
	}
	if !aDiffB.Contains("is_male") {
		t.Errorf("a diff b should contains is_male")
	}
	if !aAndB.Contains("description") {
		t.Errorf("a and b should contains description")
	}
}

func TestDictKeyOrder(t *testing.T) {
	dict := NewDict()
	for i := 0; i < 26; i += 1 {
		key := fmt.Sprintf("%c%c", 'A'+i, 'a'+i)
		dict.Add(NewInt(int64(i)), key)
	}
	dictStr := dict.String()
	for i := 0; i < 26; i += 1 {
		if dict.String() != dictStr {
			t.Logf("dict string changed!!!")
		}
	}
}
