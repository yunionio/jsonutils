package jsonutils

import (
	"reflect"
	"testing"
)

func TestDeepCopyBasic(t *testing.T) {
	cases := []struct {
		name string
		in   JSONObject
	}{
		{
			name: "string",
			in:   NewString(`hello world`),
		},
		{
			name: "int",
			in:   NewInt(2048),
		},
		{
			name: "float",
			in:   NewFloat(2048.4096),
		},
		{
			name: "bool true",
			in:   JSONTrue,
		},
		{
			name: "bool false",
			in:   JSONFalse,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			vc := DeepCopy(c.in)
			if !reflect.DeepEqual(vc, c.in) {
				t.Errorf("want %s, got %s", c.in.String(), vc.String())
			}
		})
	}
}

func TestDeepCopyCompound(t *testing.T) {
	cases := []struct {
		name string
		in   string
	}{
		{
			name: "dict",
			in:   `{"a": 0, "b": true, "c": 3.14}`,
		},
		{
			name: "array",
			in:   `["hello world", 0, 3.14, true, false]`,
		},
		{
			name: "dict (2 level)",
			in:   `{"a": 0, "b": true, "c": 3.14, d:{}, e:[]}`,
		},
		{
			name: "array (2 level)",
			in:   `["hello world", 0, 3.14, true, false, {}, []]`,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			obj, err := ParseString(c.in)
			if err != nil {
				t.Fatalf("must parse, got %s", err)
			}
			objC := DeepCopy(obj)
			if !reflect.DeepEqual(obj, objC) {
				t.Fatalf("want %s, got %s", obj.String(), objC.String())
			}
			switch v := objC.(type) {
			case *JSONArray:
				elems, _ := v.GetArray()
				elems[0] = NewString("who knows me")
			case *JSONDict:
				v.Set("id", NewString("nobody"))
			}
			if reflect.DeepEqual(obj, objC) {
				t.Fatalf("not copied: %s", obj.String())
			}
		})
	}
}
