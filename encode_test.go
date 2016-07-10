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

	type ss struct {
		User   string  `ltsv:"user"`
		Age    uint8   `ltsv:"age"`
		Height float64 `ltsv:"height"`
		Weight float32
		Memo   string `ltsv:"-"`
	}
	s := &ss{
		User:   "songmu",
		Age:    36,
		Height: 169.1,
		Weight: 66.6,
		Memo:   "songmu.jp",
	}
	expect := []byte("user:songmu\tage:36\theight:169.1\tweight:66.6")
	buf, _ := Marshal(s)
	if string(buf) != string(expect) {
		t.Errorf("result is not expected: %s", string(buf))
	}
}
