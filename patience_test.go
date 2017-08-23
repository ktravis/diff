package diff

import (
	"fmt"
	"strings"
	"testing"
)

func diffToString(items []Item) string {
	s := make([]string, len(items))
	for i, c := range items {
		pre := " "
		switch c.Type {
		case Insertion:
			pre = "+"
		case Deletion:
			pre = "-"
		}
		s[i] = fmt.Sprintf("%s %s", pre, c.Text)
	}
	return strings.Join(s, "\n")
}

func TestPatience(t *testing.T) {
	a := strings.Split(`void func1() {
    x += 1
}

void functhreeover2() {
    x += 3/2;
}

void func2() {
    x += 2
}`, "\n")

	b := strings.Split(`void func1() {
    x += 1
}

void functhreehalves() {
    x += 1.5
}

void functhreeover2() {
    x += 3/2;
}

void func2() {
    x += 1
}`, "\n")

	expected := `  void func1() {
      x += 1
  }
  
+ void functhreehalves() {
+     x += 1.5
+ }
+ 
  void functhreeover2() {
      x += 3/2;
  }
  
  void func2() {
+     x += 1
-     x += 2
  }`

	diff := Patience(a, b)
	diffStr := diffToString(diff)
	if diffStr != expected {
		t.Fatalf("Diff was not correct:\nexpected:\n%v\ngot:\n%v", expected, diffStr)
	}
}
