package ltsv

import (
	"fmt"
	"testing"
)

var encodeTests = []struct {
	Name   string
	Input  interface{}
	Output string
	Check  func(s string) error
}{
	{
		Name: "map",
		Input: map[string]string{
			"hoge": "fuga",
			"piyo": "piyo",
		},
		Check: func(s string) error {
			expect1 := "hoge:fuga\tpiyo:piyo"
			expect2 := "piyo:piyo\thoge:fuga"
			if s != expect1 && s != expect2 {
				return fmt.Errorf("result is not expected: %s", s)
			}
			return nil
		},
	},
	{
		Name: "Simple",
		Input: &ss{
			User:   "songmu",
			Age:    36,
			Height: pfloat64(169.1),
			Weight: 66.6,
			Memo:   "songmu.jp",
		},
		Output: "user:songmu\tage:36\theight:169.1\tweight:66.6",
	},
	{
		Name: "Omit memo",
		Input: &ss{
			User: "songmu",
			Age:  36,
			Memo: "songmu.jp",
		},
		Output: "user:songmu\tage:36\tweight:0",
	},
}

func TestMarshal(t *testing.T) {
	for _, tt := range encodeTests {
		t.Logf("testing: %s\n", tt.Name)
		buf, err := Marshal(tt.Input)
		if err != nil {
			t.Errorf("%s(err): error should be nil but: %+v", tt.Name, err)
		}
		s := string(buf)
		if tt.Check != nil {
			err := tt.Check(s)
			if err != nil {
				t.Errorf("%s: %s", tt.Name, err)
			}
		} else {
			if s != tt.Output {
				t.Errorf("%s: out=%s, want=%s", tt.Name, s, tt.Output)
			}
		}
	}
}

func BenchmarkMarshalStruct(b *testing.B) {
	input := &ss{
		User:   "songmu",
		Age:    36,
		Height: pfloat64(169.1),
		Weight: 66.6,
		Memo:   "songmu.jp",
	}
	for i := 0; i < b.N; i++ {
		_, err := Marshal(input)
		if err != nil {
			b.Error(err)
		}
	}
}

func BenchmarkMarshalMap(b *testing.B) {
	input := map[string]string{
		"hoge": "fuga",
		"piyo": "piyo",
	}
	for i := 0; i < b.N; i++ {
		_, err := Marshal(input)
		if err != nil {
			b.Error(err)
		}
	}
}
