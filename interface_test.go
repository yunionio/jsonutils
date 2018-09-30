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
	json := NewFloat(float64(oval))
	val := json.Interface()
	if val.(float64) != float64(oval) {
		t.Errorf("should be equal")
	}
}