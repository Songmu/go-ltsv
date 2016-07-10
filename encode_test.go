package ltsv

import "testing"

func TestMarshal(t *testing.T) {
	data := map[string]string{
		"hoge": "fuga",
		"piyo": "piyo",
	}
	expect1 := []byte("hoge:fuga\tpiyo:piyo")
	expect2 := []byte("piyo:piyo\thoge:fuga")
	r, _ := Marshal(data)
	if string(r) != string(expect1) && string(r) != string(expect2) {
		t.Errorf("result is not expected: %s", string(r))
	}
}
