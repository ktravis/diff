package diff_test

import (
	"fmt"
	"strings"

	"github.com/ktravis/diff"
)

func ExamplePatience() {
	a := `
	func Foo(a, b int) int {
}
type Bar struct {
	A string
	B bool
}`

	b := `
	func Foo(a, b int) int {
	return a + b
}
type Bar struct {
	A string
	C []byte
}`

	d := diff.Patience(strings.Split(a, "\n"), strings.Split(b, "\n"))
	fmt.Printf("%s\n", diff.Items(d))
	// Output:
	//  func Foo(a, b int) int {
	// +	return a + b
	//  }
	//  type Bar struct {
	//  	A string
	// +	C []byte
	// -	B bool
	//  }
}

func ExampleMerge() {
	// Note: the periods below are added to maintain indentation.
	original := `
func Foo(a, b int) int {
}
type Bar struct {
... A string
... B bool
}`

	a := `
func Foo(a, b int) int {
... return a + b
}
type Bar struct {
... A string
... C []byte
}`

	b := `
func Foo(a, b int) int {
}
type Bar struct {
... A string
... B bool
... D int
}`

	lines := func(s string) []string {
		return strings.Split(s, "\n")
	}

	merged := diff.Merge(lines(original), lines(a), lines(b))
	fmt.Printf("%s\n", strings.Join(merged, "\n"))
	// Output:
	// func Foo(a, b int) int {
	// ... return a + b
	// }
	// type Bar struct {
	// ... A string
	// ... C []byte
	// ... D int
	// }
}
