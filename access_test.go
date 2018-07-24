package jsonutils

import (
	"testing"
)

func TestNewDict(t *testing.T) {
	dict := NewDict()
	dict.Add(NewString("1"), "a", "b")
	dict2, _ := ParseString("{\"a\": {\"b\": \"1\"}}")
	if dict.String() != dict2.String() {
		t.Errorf("Fail %s != %s", dict, dict2)
	}
	dict = NewDict()
	dict2, _ = ParseString("{}")
	if dict.String() != dict2.String() {
		t.Errorf("Fail %s != %s", dict, dict2)
	}
}

func TestNewArray(t *testing.T) {
	arr := NewArray()
	arr.Add(NewString("1"), NewInt(1), NewFloat(1.0))
	arr2, _ := ParseString("[\"1\", 1, 1.0]")
	if arr.String() != arr2.String() {
		t.Errorf("Fail %s != %s", arr, arr2)
	}
	arr = NewArray()
	arr2, _ = ParseString("[]")
	if arr.String() != arr2.String() {
		t.Errorf("Fail %s != %s", arr, arr2)
	}
}
