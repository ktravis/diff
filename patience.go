// Package diff implements the patience diff algorithm, as described by
// http://alfedenzo.livejournal.com/170301.html
// Also provided are functions to compute the Longest Common Subsequence
// of string slices, and a custom n-way Merge function for lines of text.
package diff

type ItemType int

const (
	Insertion ItemType = iota
	Unchanged
	Deletion
)

func (t ItemType) String() string {
	switch t {
	case Insertion:
		return "+"
	case Unchanged:
		return " "
	case Deletion:
		return "-"
	}
	return "<error>"
}

// Item represents a single node in a diff result - it contains the type, which
// indicates if the wrapped text was inserted, deleted, or unchanged from the
// original, as well as the text itself alongside its index in the source.
type Item struct {
	Type ItemType
	Line
}

func (i Item) String() string {
	return i.Type.String() + i.Text
}

// Items is a wrapper type for a slice of Item structs which facilitates
// "pretty-printing" of the lines.
type Items []Item

func (it Items) String() string {
	out := ""
	for _, x := range it {
		if len(out) > 0 {
			out += "\n"
		}
		out += x.String()
	}
	return out
}

func uniqueLines(a []Line) (out []Line) {
	type entry struct {
		index, count int
	}
	m := make(map[string]entry)
	for _, x := range a {
		e := m[x.Text]
		e.count += 1
		e.index = x.Index
		m[x.Text] = e
	}
	for _, x := range a {
		if m[x.Text].count == 1 {
			out = append(out, x)
		}
	}
	return
}

// Patience implements the patience diff algorithm, as described by
// http://alfedenzo.livejournal.com/170301.html, between to slices
// of strings.
func Patience(a, b []string) []Item {
	if len(a) == 0 && len(b) == 0 {
		return nil
	}

	aLines := ToLines(a)
	bLines := ToLines(b)

	items := make([]Item, 0)
	if len(a) == 0 {
		for _, line := range bLines {
			items = append(items, Item{Insertion, line})
		}
		return items
	}
	if len(b) == 0 {
		for _, line := range aLines {
			items = append(items, Item{Deletion, line})
		}
		return items
	}

	for i := 0; i < len(a) && i < len(b); i++ {
		if a[i] != b[i] {
			break
		}
		items = append(items, Item{Unchanged, aLines[i]})
	}
	if n := len(items); n != 0 {
		for _, item := range Patience(a[n:], b[n:]) {
			item.Index += n
			items = append(items, item)
		}
		return items
	}

	suffix := make([]Item, 0)
	for i := 0; i < len(a) && i < len(b); i++ {
		if a[len(a)-i-1] != b[len(b)-i-1] {
			break
		}
		item := Item{
			Type: Unchanged,
			Line: aLines[len(a)-i-1],
		}
		suffix = append([]Item{item}, suffix...)
	}
	if len(suffix) != 0 {
		for _, item := range Patience(a[:len(a)-len(suffix)], b[:len(b)-len(suffix)]) {
			items = append(items, item)
		}
		items = append(items, suffix...)
		return items
	}

	table := NewLCSTable(uniqueLines(aLines), uniqueLines(bLines))
	lcs := table.LongestCommonSubsequence()

	if len(lcs) == 0 {
		table := NewLCSTable(aLines, bLines)
		for _, d := range table.Diff() {
			items = append(items, d)
		}
		return items
	}

	lastA := 0
	lastB := 0

	for _, x := range lcs {
		for _, item := range Patience(a[lastA:x[0]], b[lastB:x[1]]) {
			item.Index += lastA
			items = append(items, item)
		}
		items = append(items, Item{
			Type: Unchanged,
			Line: aLines[x[0]],
		})
		lastA = x[0] + 1
		lastB = x[1] + 1
	}
	for _, item := range Patience(a[lastA:], b[lastB:]) {
		item.Index += lastA
		items = append(items, item)
	}

	return items
}
