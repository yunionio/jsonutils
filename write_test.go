package jsonutils

import "testing"

func TestJSONArrayStringNilElement(t *testing.T) {
	arr := &JSONArray{data: []JSONObject{NewString("a"), nil, NewInt(1)}}
	got := arr.String()
	want := `["a",null,1]`
	if got != want {
		t.Fatalf("got %s want %s", got, want)
	}

	desc := NewDict()
	desc.Set("nics", arr)
	cfg := NewDict()
	cfg.Set("desc", desc)
	if cfg.String() != `{"desc":{"nics":["a",null,1]}}` {
		t.Fatalf("nested dict got %s", cfg.String())
	}
}

func TestNewArrayNormalizesNil(t *testing.T) {
	arr := NewArray(NewString("a"), nil, NewInt(1))
	if arr.String() != `["a",null,1]` {
		t.Fatalf("got %s", arr.String())
	}
}
