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
	"strings"
	"testing"
	"time"
)

// buildKeysDict builds a json object whose keys are in the requested order
func buildKeysDict(n int, descending bool) string {
	sb := &strings.Builder{}
	sb.WriteByte('{')
	for i := 0; i < n; i++ {
		key := i + 1
		if descending {
			key = n - i
		}
		if i > 0 {
			sb.WriteByte(',')
		}
		fmt.Fprintf(sb, "\"k%07d\":%d", key, key)
	}
	sb.WriteByte('}')
	return sb.String()
}

func TestParseDictKeyOrder(t *testing.T) {
	for _, descending := range []bool{false, true} {
		jo, err := Parse([]byte(buildKeysDict(1000, descending)))
		if err != nil {
			t.Fatalf("parse (descending=%v): %v", descending, err)
		}
		dict := jo.(*JSONDict)
		if dict.Length() != 1000 {
			t.Fatalf("got %d keys, want 1000", dict.Length())
		}
		keys := dict.SortedKeys()
		for i := 1; i < len(keys); i++ {
			if keys[i-1] >= keys[i] {
				t.Fatalf("keys not sorted at %d: %s >= %s", i, keys[i-1], keys[i])
			}
		}
		if v, err := dict.Int("k0000001"); err != nil || v != 1 {
			t.Errorf("k0000001 = %d, %v; want 1", v, err)
		}
		if v, err := dict.Int("k0001000"); err != nil || v != 1000 {
			t.Errorf("k0001000 = %d, %v; want 1000", v, err)
		}
	}
}

func TestParseDictDuplicateKey(t *testing.T) {
	jo, err := ParseString(`{"a":1,"b":2,"a":3}`)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	dict := jo.(*JSONDict)
	if dict.Length() != 2 {
		t.Fatalf("got %d keys, want 2", dict.Length())
	}
	if v, err := dict.Int("a"); err != nil || v != 3 {
		t.Errorf("a = %d, %v; want the last value 3", v, err)
	}
	if v, err := dict.Int("b"); err != nil || v != 2 {
		t.Errorf("b = %d, %v; want 2", v, err)
	}
}

func TestParseDictManyKeys(t *testing.T) {
	const n = 200000

	// keys are in descending order, which is the worst case for a sorted
	// map that shifts its tail on every out of order insert
	data := buildKeysDict(n, true)

	start := time.Now()
	jo, err := Parse([]byte(data))
	elapsed := time.Since(start)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if elapsed > 30*time.Second {
		t.Errorf("parsing %d keys took %v", n, elapsed)
	}

	dict := jo.(*JSONDict)
	if dict.Length() != n {
		t.Fatalf("got %d keys, want %d", dict.Length(), n)
	}
	keys := dict.SortedKeys()
	for i := 1; i < len(keys); i++ {
		if keys[i-1] >= keys[i] {
			t.Fatalf("keys not sorted at %d: %s >= %s", i, keys[i-1], keys[i])
		}
	}
	if v, err := dict.Int(fmt.Sprintf("k%07d", n)); err != nil || v != int64(n) {
		t.Errorf("last key = %d, %v; want %d", v, err, n)
	}
}

func BenchmarkParseDictDescendingKeys(b *testing.B) {
	data := []byte(buildKeysDict(10000, true))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := Parse(data); err != nil {
			b.Fatalf("parse: %v", err)
		}
	}
}
