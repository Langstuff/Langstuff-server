package markdown_parsing

import "testing"

var cases = []struct{
	source string
	a string
	b string
} {
	{"a : b", "a", "b"},
	{" a:b", "a", "b"},
	{"a:b  ", "a", "b"},
	{"a:b ", "a", "b"},
	{"a     :b ", "a", "b"},
	{"a:  b ", "a", "b"},
	{"  test a :  test b ", "test a", "test b"},
	{"\ttest a\t:\ttest b\t", "test a", "test b"},
}

func TestParsePair(t *testing.T) {
	for _, testCase := range cases {
		a, b, err := parsePair(testCase.source)
		if err != nil {
			t.Error(testCase.source + ": " + err.Error())
		}
		if a != testCase.a {
			t.Errorf("Expected '%s', got '%s'", testCase.a, a)
		}
		if b != testCase.b {
			t.Errorf("Expected '%s', got '%s'", testCase.b, b)
		}
	}
}
