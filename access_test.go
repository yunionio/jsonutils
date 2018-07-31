package jsonutils

import (
	"reflect"
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

func TestNewBool(t *testing.T) {
	type args struct {
		val bool
	}
	tests := []struct {
		name string
		args args
		want *JSONBool
	}{
		{
			name: "New-bool-true",
			args: args{true},
			want: &JSONBool{data: true},
		},
		{
			name: "New-bool-false",
			args: args{false},
			want: &JSONBool{data: false},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewBool(tt.args.val); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewBool() = %v, want %v", got, tt.want)
			}
		})
	}
}
