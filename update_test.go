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

type testStruct struct {
	Name     string
	Gender   bool
	Age      int
	Position string
	nickname string
}

func TestUpdate(t *testing.T) {
	t1 := testStruct{Name: "alice", Gender: false, Age: 24, Position: "staff", nickname: "ali"}
	t2 := testStruct{Name: "bob", Gender: true, Age: 40, Position: "engineer", nickname: "bb"}
	t.Logf("t1: %s t1.nick: %s", Marshal(&t1).String(), t1.nickname)
	t.Logf("t2: %s t2.nick: %s", Marshal(&t2).String(), t2.nickname)
	err := Update(&t2, &t1)
	if err != nil {
		t.Errorf("update error %s", err)
	} else {
		t.Logf("t1: %s t1.nick: %s", Marshal(&t1).String(), t1.nickname)
		t.Logf("t2: %s t2.nick: %s", Marshal(&t2).String(), t2.nickname)
	}
}
