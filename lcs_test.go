package diff

import (
	"fmt"
	"strings"
	"testing"
)

func splitLines(s string) (out []Line) {
	for i, l := range strings.Split(s, "\n") {
		out = append(out, Line{i, strings.TrimSpace(l)})
	}
	return
}

func TestLCSTable(t *testing.T) {
	cases := []struct {
		a, b            []Line
		expectedLengths []int
		expectedItems   string
	}{
		{
			splitLines("a\nx\nb"),
			splitLines("a\nb\nc"),
			[]int{0, 0, 0, 0, 0, 1, 1, 1, 0, 1, 1, 1, 0, 1, 2, 2},
			`  a
- x
  b
+ c`,
		},
		{
			splitLines("g\na\nc"),
			splitLines("a\ng\nc\na\nt"),
			[]int{0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 1, 0, 1, 1, 1, 2, 2, 0, 1, 1, 2, 2, 2},
			`+ a
  g
+ c
  a
+ t
- c`,
		},
	}

	for i, c := range cases {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			table := NewLCSTable(c.a, c.b)
			if len(table.lengths) != len(c.expectedLengths) {
				t.Fatalf("LCSTable lengths were not correct:\nexpected: %v\ngot:      %v", c.expectedLengths, table.lengths)
			}
			for i, l := range table.lengths {
				if l != c.expectedLengths[i] {
					t.Errorf("LCSTable lengths were not correct:\nexpected: %v\ngot:      %v", c.expectedLengths, table.lengths)
					break
				}
			}
			diff := table.Diff()
			diffStr := diffToString(diff)
			if diffStr != c.expectedItems {
				t.Fatalf("Diff was not correct:\nexpected:\n%v\ngot:\n%v", c.expectedItems, diffStr)
			}
		})
	}
}

func TestLongestCommonSubsequence(t *testing.T) {
	cases := []struct {
		a, b        []Line
		expectedLcs [][2]int
	}{
		{
			splitLines("a\nx\nb"),
			splitLines("a\nb\nc"),
			[][2]int{{0, 0}, {2, 1}},
		},
		{
			splitLines("X\nX\nX\na\nX\nX\nX\nb\nX\nX\nX\nc"),
			splitLines("Y\nY\na\nY\nY\nb\nY\nY\nc"),
			[][2]int{{3, 2}, {7, 5}, {11, 8}},
		},
	}
	for i, c := range cases {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			table := NewLCSTable(c.a, c.b)
			lcs := table.LongestCommonSubsequence()
			if len(lcs) != len(c.expectedLcs) {
				t.Fatalf("LCS was not correct:\nexpected:\n%v\ngot:\n%v", c.expectedLcs, lcs)
			}
			for i := range c.expectedLcs {
				if lcs[i] != c.expectedLcs[i] {
					t.Fatalf("LCS was not correct:\nexpected:\n%v\ngot:\n%v", c.expectedLcs, lcs)
				}
			}
		})
	}
}
