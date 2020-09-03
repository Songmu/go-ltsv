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

type ss struct {
	User   string   `ltsv:"user"`
	Age    uint8    `ltsv:"age"`
	Height *float64 `ltsv:"height"`
	Weight float32
	Memo   string `ltsv:"-"`
}

func pfloat64(f float64) *float64 {
	return &f
}

var decodeTests = []struct {
	Name   string
	Input  string
	Output *ss
}{
	{
		Name:  "Simple",
		Input: "user:songmu\tage:36\theight:169.1\tweight:66.6",
		Output: &ss{
			User:   "songmu",
			Age:    36,
			Height: pfloat64(169.1),
			Weight: 66.6,
		},
	},
	{
		Name:  "Default values",
		Input: "user:songmu\tage:36",
		Output: &ss{
			User:   "songmu",
			Age:    36,
			Height: nil,
			Weight: 0.0,
		},
	},
	{
		Name:  "Hyphen and empty string as null number",
		Input: "user:songmu\tage:\theight:-",
		Output: &ss{
			User:   "songmu",
			Age:    0,
			Height: nil,
			Weight: 0.0,
		},
	},
}

func TestUnmarshal(t *testing.T) {
	m := make(map[string]string)
	Unmarshal([]byte("hoge: fuga\tpiyo: piyo"), &m)
	expect := map[string]string{
		"hoge": "fuga",
		"piyo": "piyo",
	}
	if !reflect.DeepEqual(m, expect) {
		t.Errorf("unmarsharl error:\n out:  %+v\n want: %+v", m, expect)
	}

	for _, tt := range decodeTests {
		t.Logf("testing: %s\n", tt.Name)
		s := &ss{}

		err := Unmarshal([]byte(tt.Input), s)
		if err != nil {
			t.Errorf("%s(err): error should be nil but: %+v", tt.Name, err)
		}

		if !reflect.DeepEqual(s, tt.Output) {
			t.Errorf("%s:\n out:  %+v\n want: %+v", tt.Name, s, tt.Output)
		}
	}
}

func BenchmarkUnmarshalStruct(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for _, tt := range decodeTests {
			s := &ss{}
			err := Unmarshal([]byte(tt.Input), s)
			if err != nil {
				b.Error(err)
			}
		}
	}
}
