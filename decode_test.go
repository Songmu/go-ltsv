package ltsv

import (
	"reflect"
	"testing"
)

func TestData2Map(t *testing.T) {
	l, _ := data2map([]byte("hoge: fuga\tpiyo: piyo"))

	expect := ltsvMap{
		"hoge": "fuga",
		"piyo": "piyo",
	}
	if !reflect.DeepEqual(l, expect) {
		t.Errorf("result of data2map not expected: %#v", l)
	}

	_, err := data2map([]byte("hoge"))
	if err.Error() != "not a ltsv: hoge" {
		t.Errorf("something went wrong")
	}
}

func TestUnmarshal(t *testing.T) {
	m := make(map[string]string)
	Unmarshal([]byte("hoge: fuga\tpiyo: piyo"), &m)

	expect := map[string]string{
		"hoge": "fuga",
		"piyo": "piyo",
	}
	if !reflect.DeepEqual(m, expect) {
		t.Errorf("result of data2map not expected: %#v", m)
	}

	type ss struct {
		User   string   `ltsv:"user"`
		Age    uint8    `ltsv:"age"`
		Height *float64 `ltsv:"height"`
		Weight float32
	}
	s := &ss{}
	Unmarshal([]byte("user:songmu\tage:36\theight:169.1\tweight:66.6"), s)
	hei := 169.1
	expect2 := &ss{
		User:   "songmu",
		Age:    36,
		Height: &hei,
		Weight: 66.6,
	}
	if !reflect.DeepEqual(s, expect2) {
		t.Errorf("result of data2map not expected: %#v", s)
	}

	s2 := &ss{}
	Unmarshal([]byte("user:songmu\tage:36"), s2)
	expect3 := &ss{
		User: "songmu",
		Age:  36,
		// Height: nil,
		Weight: 0.0,
	}
	if !reflect.DeepEqual(s2, expect3) {
		t.Errorf("result of data2map not expected: %#v", s2)
	}

	s3 := &ss{}
	Unmarshal([]byte("user:songmu\tage:-\theight:-"), s3)
	expect4 := &ss{
		User: "songmu",
		Age:  0,
		// Height: nil,
		Weight: 0.0,
	}
	if !reflect.DeepEqual(s3, expect4) {
		t.Errorf("result of data2map not expected: %#v", s3)
	}

}
