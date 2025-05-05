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

func TestHexchar2num(t *testing.T) {
	cases := []struct {
		in, want byte
	}{
		{'F', 15},
		{'A', 10},
		{'0', 0},
		{'1', 1},
	}
	for _, c := range cases {
		got, _ := hexchar2num(c.in)
		if got != c.want {
			t.Errorf("Hexchar2num(%c) == %c, want %c", c.in, got, c.want)
		}
	}
	_, e := hexchar2num('G')
	if e == nil {
		t.Errorf("Hexchar2num(G) should raise error")
	}
}

func TestHexstr2byte(t *testing.T) {
	cases := []struct {
		in   []byte
		want byte
	}{
		{[]byte{'F', 'F'}, 255},
		{[]byte{'0', '0'}, 0},
		{[]byte{'1', '0'}, 16},
	}
	for _, c := range cases {
		got, _ := hexstr2byte(c.in)
		if got != c.want {
			t.Errorf("hexstr2byte(%s) == %d, want %d", c.in, got, c.want)
		}
	}
}

func TestHexstr2rune(t *testing.T) {
	cases := []struct {
		in   []byte
		want rune
	}{
		{[]byte("00FF"), 255},
		{[]byte("0000"), 0},
		{[]byte("0010"), 16},
	}
	for _, c := range cases {
		got, _ := hexstr2rune(c.in)
		if got != c.want {
			t.Errorf("hexstr2rune(%s) == %d, want %d", c.in, got, c.want)
		}
	}
}

func TestReadString(t *testing.T) {
	cases := []struct {
		in         []byte
		want       string
		want_quote bool
	}{
		{[]byte("\"00FF\""), "00FF", true},
		{[]byte("0"), "0", false},
		{[]byte("\"a\\nb\\n\""), "a\nb\n", true},
		{[]byte("123\n22"), "123", false},
		{[]byte("abc:"), "abc", false},
	}
	for _, c := range cases {
		got, quote, _, _ := parseString(c.in, 0)
		if got != c.want || quote != c.want_quote {
			t.Errorf("readString(%s) == %s %v, want %s %v", c.in, got, quote, c.want, c.want_quote)
		}
	}
}

func TestJSONParse(t *testing.T) {
	cases := []struct {
		in, out string
	}{
		{"{'name': '大家好'}", `{"name": "\xe5\xa4\xa7\xe5\xae\xb6\xe5\xa5\xbd"}`},
		{`["\xe5\xa5\xbd"]`, `["好"]`},
	}
	for _, c := range cases {
		t.Logf("in: %s out: %s", c.in, c.out)
		got, _ := ParseString(c.in)
		got2, _ := ParseString(c.out)
		t.Logf("%s %s", got, got2)
		if got.String() != got2.String() {
			t.Errorf("JSONParse: %s(%s) != %s(%s)", c.in, got, c.out, got2)
		}
	}
}

func BenchmarkParseString(b *testing.B) {
	cases := []struct {
		name string
		c    string
	}{
		{
			name: "all",
			c:    `{"abc": 12, "def": [1,2,"123",4.43], "ghi": "hahahah"}`,
		},
	}

	for _, c := range cases {
		b.Run(c.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				ParseString(c.c)
			}
		})
	}
}

func Benchmark_quoteString(b *testing.B) {
	cases := []struct {
		name string
		in   string
	}{
		{
			name: "no escape",
			in:   "hello world",
		},
		{
			name: "escape",
			in:   "hello\nworld\r\t\\",
		},
	}
	for _, c := range cases {
		b.Run(c.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				quoteString(c.in)
			}
		})
	}
}

func BenchmarkStringify(b *testing.B) {
	cases := []struct {
		name string
		c    string
		obj  JSONObject
	}{
		{
			name: "all",
			c:    `{"abc": 12, "def": [1,2,"123",4.43], "ghi": "hahahah"}`,
		},
	}
	for _, c := range cases {
		var err error
		c.obj, err = ParseString(c.c)
		if err != nil {
			b.Fatalf("%s: bad case: %v", c.name, err)
		}

		b.Run(c.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				c.obj.String()
			}
		})
	}
}

func TestParseParams(t *testing.T) {
	str := `{"__request_context":{"request_id":"cbde96","service_name":"region","trace":{"debug":true,"duration":0,"id":"0","kind":"SERVER","local_endpoint":{"port":0,"service_name":"region"},"name":"delete","remote_endpoint":{"addr":"10.168.222.188","port":57240,"service_name":"(unknown_service)"},"shared":false,"tags":{"resource":"cloudaccounts"},"timestamp":"2020-04-11T06:37:07.100833Z","trace_id":"1866608c"}},"__stages":[{"complete_at":"2020-04-11T14:44:57Z","name":"on_init"}],"parent_task_id":"65617a87-3ecd-40c0-8add-70e17fec8ab2"}`
	json, err := ParseString(str)
	if err != nil {
		t.Fatalf("%s", err)
	}
	if str != json.String() {
		t.Fatalf("%s!=%s", str, json.String())
	}
}

func TestParseJsonStreams(t *testing.T) {
	cases := []struct {
		jsonStr []byte
		jsonLen int
	}{
		{
			jsonStr: []byte(`{"abc": 12, "def": [1,2,"123",4.43], "ghi": "hahahah"} {"abc": 12}`),
			jsonLen: 2,
		},
		{
			jsonStr: []byte(`    {"abc": 12, "def": [1,2,"123",4.43], "ghi": "hahahah"} {"abc": 12}         `),
			jsonLen: 2,
		},
		{
			jsonStr: []byte(`    {"abc": 12, "def": [1,2,"123",4.43], "ghi": "hahahah"} [] {"abc": 12}         `),
			jsonLen: 3,
		},
		{
			jsonStr: []byte(`    {"abc": 12, "def": [1,2,"123",4.43], "ghi": "hahahah"} [] {"abc": 12}     {}   [] `),
			jsonLen: 5,
		},
		{
			jsonStr: []byte(`  abc  {"abc": 12, "def": [1,2,"123",4.43], "ghi": "hahahah"} [] {"abc": 12}     {}   [] `),
			jsonLen: 5,
		},
		{
			jsonStr: []byte(`  abc ["abc"]  {"abc": 12,
			 "def": [1,2,"123",4.43], "ghi": "hahahah"}
			 [
			 ] {"abc": 12}
			   {}   [] `),
			jsonLen: 6,
		},
		{
			jsonStr: []byte(`  abc  [{"abc": 12, "def": [1,2,"123",4.43], "ghi": "hahahah"} ["abc"] {"abc": 12}     {}   [] `),
			jsonLen: 5,
		},
		{
			jsonStr: []byte(`  abc   {"abc": 12, "def": [1,2,"123",4.43], "ghi": "hahahah"} ["abc"] {"abc": 12}     {}   [] `),
			jsonLen: 5,
		},
	}
	for _, c := range cases {
		results, err := ParseJsonStreams(c.jsonStr)
		if err != nil {
			t.Errorf("%s", err)
		} else if len(results) != c.jsonLen {
			t.Errorf("got %d != expect %d", len(results), c.jsonLen)
		}
		t.Logf("%s", Marshal(results))
	}
}
