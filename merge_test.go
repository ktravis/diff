package diff

import (
	"fmt"
	"strings"
	"testing"
)

func TestMerge(t *testing.T) {
	cases := []struct {
		original string
		a, b     string
		expected string
	}{
		{
			original: `void func1() {
    x += 1
}

void func2() {
    x += 2
}`,
			a: `void func1() {
    x += 1
}

void functhreeover2() {
    x += 3/2;
}

void func2() {
    x += 1
}`,
			b: `void func1() {
    x += 1
}

void functhreehalves() {
    x += 1.5
}

void func2() {
    x += 1
}`,
			expected: `void func1() {
    x += 1
}

void functhreeover2() {
    x += 3/2;
}

void functhreehalves() {
    x += 1.5
}

void func2() {
    x += 1
}`,
		},
		{
			original: "a\nb\nc",
			a:        "a\nx1\nx\nc\na",
			b:        "a\nx1\nc",
			expected: "a\nx1\nx\nc\na",
		},
		{
			original: "a\nb\nc",
			a:        "1\n2\n3\n4",
			b:        "d",
			expected: "1\n2\n3\n4\nd",
		},
		{
			original: "a\na\na",
			a:        "a",
			b:        "a\nb",
			expected: "a\nb",
		},
	}

	for i, c := range cases {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			original := strings.Split(c.original, "\n")
			a := strings.Split(c.a, "\n")
			b := strings.Split(c.b, "\n")
			actual := strings.Join(Merge(original, a, b), "\n")
			if actual != c.expected {
				t.Fatalf("expected and actual merge did not match:\nexpected:\n%s\nactual:\n%s", c.expected, actual)
			}
		})
	}
}

func TestMergeMulti(t *testing.T) {
	a := strings.Split(`void func1() {
    x += 1
}

void func2() {
    x += 2
}`, "\n")
	b := strings.Split(`void func1() {
    x += 1
}

void functhreeover2() {
    x += 3/2;
}

void func2() {
    x += 1
}`, "\n")

	c := strings.Split(`void func1() {
    x += 1
}

void functhreehalves() {
    x += 1.5
}

void func2() {
    x += 1
}`, "\n")

	d := strings.Split(`void func1() {
    x += 1
    print("hi there")
}

void func2() {
    x += 3
}`, "\n")

	expected := `void func1() {
    x += 1
    print("hi there")
}

void functhreeover2() {
    x += 3/2;
}

void functhreehalves() {
    x += 1.5
}

void func2() {
    x += 1
    x += 3
}`
	actual := strings.Join(Merge(a, b, c, d), "\n")
	if actual != expected {
		t.Fatalf("expected and actual merge did not match:\nexpected:\n%s\nactual:\n%s", expected, actual)
	}
}
