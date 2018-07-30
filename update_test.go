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
