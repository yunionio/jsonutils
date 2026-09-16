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
	"errors"
	"testing"
)

func TestParseTrailingContent(t *testing.T) {
	invalid := []string{
		`{"a":1}{"b":2}`,
		`{"a":1} trailing`,
		`[1,2] x`,
		`1 2`,
		`{} []`,
		`null null`,
		`"a" "b"`,
	}
	for _, c := range invalid {
		if _, err := Parse([]byte(c)); err == nil {
			t.Errorf("Parse(%q) expect an error", c)
		}
	}

	valid := []string{
		`{"a":1}`,
		`{"a":1} `,
		"{\"a\":1}\n",
		"{\"a\":1}\t\r\n  ",
		`[1,2]`,
		`null`,
	}
	for _, c := range valid {
		if _, err := Parse([]byte(c)); err != nil {
			t.Errorf("Parse(%q): %v", c, err)
		}
	}
}

func TestParseStreamKeepsOffset(t *testing.T) {
	// ParseStream 是流式解析的入口，仍需返回消费到的位置
	jo, offset, err := ParseStream([]byte(`{"a":1}{"b":2}`), 0)
	if err != nil {
		t.Fatalf("ParseStream: %v", err)
	}
	if offset != 7 {
		t.Fatalf("offset = %d, want 7", offset)
	}
	if v, err := jo.Int("a"); err != nil || v != 1 {
		t.Errorf("a = %d, %v; want 1", v, err)
	}
}

type errorMarshalJSON struct{}

func (errorMarshalJSON) MarshalJSON() ([]byte, error) {
	return nil, errors.New("boom")
}

type invalidMarshalJSON struct{}

func (invalidMarshalJSON) MarshalJSON() ([]byte, error) {
	return []byte(`{"a":1} not json`), nil
}

func TestMarshalMarshalJSONFailure(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Marshal panicked: %v", r)
		}
	}()

	for _, v := range []interface{}{errorMarshalJSON{}, invalidMarshalJSON{}} {
		if got := Marshal(v); got.String() != "null" {
			t.Errorf("Marshal(%T) = %s, want null", v, got)
		}
	}
}
