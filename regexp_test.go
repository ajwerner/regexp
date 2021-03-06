package regexp

import (
	"fmt"
	"testing"
)

func TestCases(t *testing.T) {
	for _, c := range []testCase{
		{"a+", "a", true, false},
		{"a+", "", false, false},
		{"a", "aa", false, false},
		{"abc", "abc", true, false},
		{"abd", "abc", false, false},
		{"a.c", "a.c", true, false},
		{"a.c", "abc", true, false},
		{"a?bc", "abc", true, false},
		{"a?bc", "bc", true, false},
		{"a?ab", "ab", true, false},
		{"a+a+a+aa", "aaaaaa", true, false},
		{"a+a+a+aa", "aaaaaaaaaaa", true, false},
		{"a+a+a+aa", "aaaaaaaaaaaaaaaaaaa", true, false},
		{"a+a*a*aa", "aaaaaaaaaaaaaaaaaaa", true, false},
		{".+.", "asdfasdfasdf", true, false},
		{"..", "aaa", false, false},
		{"...", "aaa", true, false},
		{"....", "aa", false, false},
		{"ERROR: .*", "ERROR: file not found", true, false},
		{"ERROR: .*", "WARNING: file not found", false, false},
		{"**", "", false, true},
		{"(**)", "", false, true},
		{"*.*", "", false, true},
		{"+.?.", "aa", false, true},
		{"+.*", "a", false, true},
		{
			"a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?a?aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
			"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
			true, false,
		},
		{"", "", true, false},
		{"", "a", false, false},
		{"a", "", false, false},
		{"😃+", "😃😃😃😃😃😃", true, false},
		{"a|b", "a", true, false},
		{"a|b", "b", true, false},
		{"(ab)?a", "aba", true, false},
		{"(ab)?a", "a", true, false},
		{"(ab)*a", "ababababababababababababababababababababababababababababababababababababababababababababababababababababababababababababababababababababababababababababababababababababababababababababababababababababababababababababababababababababababababababababababababababababababababababababababababababababababababababababababababababababababababababababababababababababababababababababababababababababababababababababababababababababababababababababababababababababababababababababababababababababababababababababababababababababababababababababababababababababababababababababababa", true, false},
		{"(ab)+a", "a", false, false},
		{"(ab)+a", "ababa", true, false},
		{"(ab)+a", "abababa", true, false},
		{"\\*\\(ab\\)a", "abababa", false, false},
		{"\\*\\(ab\\)a", "*(ab)a", true, false},
		{")", "", false, true},
		{"as)", "", false, true},
		{"(as", "", false, true},
		{"(a(as)())(", "", false, true},
		{"(a(as)())()", "aas", true, false},
		{"]", "*(ab)a", false, false},
		{"\\y", "", false, true},
		{"()|a", "a", true, false},
		{"()||", "a", false, true},
		{"|", "", false, true},
		{"|a", "", false, true},
		{"()", "", true, false},
		{"()|()", "", true, false},
		{"\u0000", "\u0000", true, false},
		{"", "\u0000", false, false},
	} {
		t.Run(c.String(), c.run)
	}
}

type testCase struct {
	pattern, in  string
	result       bool
	compileError bool
}

func (c *testCase) String() string {
	if c.compileError {
		return fmt.Sprintf("error on #%q", c.pattern)
	}
	return fmt.Sprintf("#%q.MatchString(%q) == %v", c.pattern, c.in, c.result)
}

func (c *testCase) run(t *testing.T) {
	r, err := Compile(c.pattern)
	if err == nil && c.compileError {
		t.Fatalf("Expected error for #%q", c.pattern)
	} else if err != nil && !c.compileError {
		t.Fatalf("Unexpected error for %#q: %v", c.pattern, err)
	} else if c.compileError {
	} else if matches := r.MatchString(c.in); matches != c.result {
		t.Fatalf("match(%v, %v) != %v", c.pattern, c.in, c.result)
	}
	// TODO: compare against the stdlib implementation
	// re, reErr := regexp.Compile("^"+c.pattern+"$")
	// if reErr == nil && err != nil {
	// 	t.Fatalf("regexp returned nil but we returned %v", err)
	// } else if reErr != nil && err == nil {
	// 	t.Fatalf("regexp returned %v but we returned nil", reErr)
	// } else if reErr == nil && err == nil && re.MatchString(c.in) != c.result {
	// 	t.Fatalf("regexp results in %v for %q/%q unlike us", !c.result, c.pattern, c.in)
	// }
	defer func() {
		if r := recover(); r != nil {
			if !c.compileError {
				t.Fatalf("MustCompile panicked for #%q", c.pattern)
			}
		} else if c.compileError {
			t.Fatalf("MustCompile didn't panic for #%q", c.pattern)
		}
	}()
	MustCompile(c.pattern)
}
