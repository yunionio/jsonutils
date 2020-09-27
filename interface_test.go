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

import "testing"

func TestStringInterface(t *testing.T) {
	str := "TestString"
	strJson := NewString(str)
	val := strJson.Interface()
	if val.(string) != str {
		t.Errorf("should be equal")
	}
}

func TestBoolInterface(t *testing.T) {
	json := JSONTrue
	val := json.Interface()
	if val.(bool) != true {
		t.Errorf("should be equal")
	}
}

func TestIntInterface(t *testing.T) {
	oval := 123
	json := NewInt(int64(oval))
	val := json.Interface()
	if val.(int64) != int64(oval) {
		t.Errorf("should be equal")
	}
}

func TestFloatInterface(t *testing.T) {
	oval := 123.223
	json := NewFloat64(float64(oval))
	val := json.Interface()
	if val.(float64) != float64(oval) {
		t.Errorf("should be equal")
	}
}
