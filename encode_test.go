package ltsv

import (
	"encoding"
	"errors"
	"fmt"
	"testing"
)

type ms struct {
	err error
}

func (s *ms) MarshalText() ([]byte, error) {
	return []byte("ok"), s.err
}

var _ encoding.TextMarshaler = &ms{}

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
		Name: "Simple with nil pointer",
		Input: &ss{
			User:          "songmu",
			Age:           36,
			Height:        nil,
			Weight:        66.6,
			EmailVerified: true,
			Memo:          "songmu.jp",
		},
		Output: "user:songmu\tage:36\tweight:66.6\temail_verified:true",
	},
	{
		Name: "Simple without nil pointer",
		Input: &ss{
			User:          "songmu",
			Age:           36,
			Height:        pfloat64(169.1),
			Weight:        66.6,
			EmailVerified: false,
			Memo:          "songmu.jp",
		},
		Output: "user:songmu\tage:36\theight:169.1\tweight:66.6\temail_verified:false",
	},
	{
		Name: "Omit memo",
		Input: &ss{
			User: "songmu",
			Age:  36,
			Memo: "songmu.jp",
		},
		Output: "user:songmu\tage:36\tweight:0\temail_verified:false",
	},
	{
		Name: "Anthoer struct",
		Input: &struct {
			Name  string
			Value int `ltsv:"answer"`
		}{
			Name:  "the Answer",
			Value: 42,
		},
		Output: "name:the Answer\tanswer:42",
	},
	{
		Name: "TextMarshaler",
		Input: &struct {
			Struct interface{}
		}{
			Struct: &ms{},
		},
		Output: "struct:ok",
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
				t.Errorf("%s:\n  out =%s\n  want=%s", tt.Name, s, tt.Output)
			}
		}
	}
}

func TestMarshalError(t *testing.T) {
	errOK := errors.New("ok")
	s := struct {
		A *ms
		B *ms
	}{
		A: &ms{errOK},
		B: &ms{errOK},
	}

	_, err := Marshal(s)
	if err == nil {
		t.Errorf("got no error")
	}
	if got := err.(MarshalError).OfField("a"); got != errOK {
		t.Errorf("got error: %v", got)
	}
	if got := err.(MarshalError).OfField("b"); got != errOK {
		t.Errorf("got error: %v", got)
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
