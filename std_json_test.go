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
	"bytes"
	"fmt"
	"reflect"
	"testing"
)

type MyState int

func (u MyState) MarshalJSON() ([]byte, error) {
	var buf bytes.Buffer
	switch u {
	case 1:
		_, err := buf.WriteString(`"Success"`)
		return buf.Bytes(), err
	case 2:
		_, err := buf.WriteString(`"Fail"`)
		return buf.Bytes(), err
	case 3:
		_, err := buf.WriteString(`"Pending"`)
		return buf.Bytes(), err
	}

	return nil, fmt.Errorf("unsupported value: %v!", u)
}

func (u *MyState) UnmarshalJSON(data []byte) error {
	switch string(data) {
	case `"Success"`:
		*u = 1
		return nil
	case `"Fail"`:
		*u = 2
		return nil
	case `"Pending"`:
		*u = 3
		return nil
	}
	return fmt.Errorf("unsupported value: %v!", *u)
}

func TestIsImplementStdUnmarshaler(t *testing.T) {
	tests := []struct {
		name  string
		args  reflect.Value
		isNil bool
	}{
		{
			"Implement Unmarshaler",
			reflect.Indirect(reflect.ValueOf(new(MyState))),
			false,
		},
		{
			"Implement Unmarshaler",
			reflect.ValueOf(new(MyState)),
			false,
		},
		{
			"Not implement Unmarshaler",
			reflect.ValueOf(new(map[string]string)),
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsImplementStdUnmarshaler(tt.args); !reflect.DeepEqual(got == nil, tt.isNil) {
				t.Errorf("IsImplementStdUnmarshaler() = %v, want %v", got, tt.isNil)
			}
		})
	}
}

func TestStdJSONUnmarhaler(t *testing.T) {
	jStr, _ := ParseString(`"Success"`)

	tests := []struct {
		input    string
		expected int
	}{
		{
			`"Success"`,
			1,
		},
		{
			`"Fail"`,
			2,
		},
		{
			`"Pending"`,
			3,
		},
	}

	for _, tt := range tests {
		jo, err := ParseString(tt.input)
		if err != nil {
			t.Errorf("Parse string %q", tt.input)
			return
		}
		state := new(MyState)
		if err := jo.Unmarshal(state); err != nil {
			t.Errorf("Unmarshal to state: %v", err)
			return
		}
		if int(*state) != tt.expected {
			t.Errorf("Unmarshal state = %d, want %d", *state, tt.expected)
		}
	}

	rs := new(MyState)
	if err := jStr.Unmarshal(rs); err != nil {
		t.Errorf("unmarshal error: %v", err)
	}
}

func TestStdJSONMarhaler(t *testing.T) {
	state := MyState(1)
	expected := `"Success"`
	got := Marshal(state).String()
	if expected != got {
		t.Errorf("got = %q, expected = %q", got, expected)
	}
}
