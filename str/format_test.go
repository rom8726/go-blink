package str

import (
	"testing"
)

func TestFormat(t *testing.T) {
	cases := []struct {
		Format   string
		Expected string
		Params   map[string]interface{}
	}{
		{
			"${one}",
			"1",
			map[string]interface{}{"one": 1},
		},
		{
			"${not-a-param",
			"${not-a-param",
			nil,
		},
		{
			"${one} + ${two}",
			"1 + 2",
			map[string]interface{}{"one": 1, "two": 2},
		}, {
			"add ${one} to ${two} equals 3",
			"add 1 to 2 equals 3",
			map[string]interface{}{"one": 1, "two": 2},
		}, {
			"Hello, ${absent}!",
			"Hello, !",
			map[string]interface{}{},
		}, {
			"Hello, ${name}! It's ${one} + ${two}",
			"Hello, John Doe! It's 1 + 2",
			map[string]interface{}{"name": "John Doe", "one": 1, "two": 2},
		},
	}

	for _, c := range cases {
		s := Format(c.Format, c.Params)
		if s != c.Expected {
			t.Fatal("Expected: ", c.Expected, "actual: ", s)
		}
	}
}

func TestFormatStruct(t *testing.T) {
	cases := []struct {
		Format   string
		Expected string
		Params   interface{}
	}{
		{
			"${One}",
			"1",
			struct {
				One int64
			}{1},
		},
		{
			"${not-a-param",
			"${not-a-param",
			struct{}{},
		},
		{
			"${One} + ${Two}",
			"1 + 2",
			struct {
				One int
				Two int
			}{1, 2},
		}, {
			"add ${One} to ${Two} equals 3",
			"add 1 to 2 equals 3",
			struct {
				One int
				Two int
			}{1, 2},
		}, {
			"Hello, ${Absent}!",
			"Hello, !",
			struct{}{},
		}, {
			"Hello, ${Name}! It's ${One} + ${Two}",
			"Hello, John Doe! It's 1 + 2",
			struct {
				Name string
				One  int
				Two  int
			}{"John Doe", 1, 2},
		},
	}

	for _, c := range cases {
		s := FormatStruct(c.Format, c.Params)
		if s != c.Expected {
			t.Fatal("Expected: ", c.Expected, "actual: ", s)
		}
	}
}

func BenchmarkFormatter_Format(b *testing.B) {
	f := NewFormatter("Hello, ${name}! It's ${one} + ${two}")

	for i := 0; i < b.N; i++ {
		p := map[string]interface{}{"name": "John Doe", "one": 1, "two": 2}
		f.Format(p)
	}
}

func BenchmarkFormatter_FormatStruct(b *testing.B) {
	f := NewFormatter("Hello, ${Name}! It's ${One} + ${Two}")

	for i := 0; i < b.N; i++ {
		p := struct {
			Name string
			One  int
			Two  int
		}{"John Doe", 1, 2}
		f.FormatStruct(p)
	}
}
