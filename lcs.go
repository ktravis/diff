package diff

// Line represents an atom of text from the source material, contextualized
// by its original index.
type Line struct {
	Index int
	Text  string
}

// ToLines is a convenience function for creating a slice of Line structs from a
// slice of strings as input for NewLCSTable.
func ToLines(a []string) []Line {
	out := make([]Line, len(a))
	for i, line := range a {
		out[i] = Line{i, line}
	}
	return out
}

// LCSTable is a data structure used to compute the LCS and traditional
// LCS-based diff.
type LCSTable struct {
	lengths []int
	a, b    []Line
}

// NewLCSTable constructs an LCSTable, pre-computing the necessary len(a)*len(b)
// array of lengths required for future operations.
func NewLCSTable(a, b []Line) *LCSTable {
	t := &LCSTable{
		lengths: make([]int, (len(a)+1)*(len(b)+1)),
		a:       a,
		b:       b,
	}

	for i, _ := range a {
		for j, _ := range b {
			k := (i+1)*(len(b)+1) + (j + 1)
			if a[i].Text == b[j].Text {
				t.lengths[k] = t.getLength(i, j) + 1
			} else {
				nextA := t.getLength(i+1, j)
				nextB := t.getLength(i, j+1)
				if nextA > nextB {
					t.lengths[k] = nextA
				} else {
					t.lengths[k] = nextB
				}
			}
		}
	}
	return t
}

func (t *LCSTable) getLength(ai, bi int) int {
	return t.lengths[ai*(len(t.b)+1)+bi]
}

func (t *LCSTable) LongestCommonSubsequence() [][2]int {
	return t.recursiveLcs(len(t.a), len(t.b))
}

func (t *LCSTable) recursiveLcs(i, j int) [][2]int {
	if i == 0 || j == 0 {
		return nil
	}
	if t.a[i-1].Text == t.b[j-1].Text {
		next := [2]int{t.a[i-1].Index, t.b[j-1].Index}
		return append(t.recursiveLcs(i-1, j-1), next)
	}
	if t.getLength(i, j-1) > t.getLength(i-1, j) {
		return t.recursiveLcs(i, j-1)
	}
	return t.recursiveLcs(i-1, j)
}

// Diff returns a diff of the two sets of lines the LCSTable was created with,
// as determined by the LCS.
func (t *LCSTable) Diff() []Item {
	return t.recursiveDiff(len(t.a), len(t.b))
}

func (t *LCSTable) recursiveDiff(i, j int) []Item {
	if i == 0 && j == 0 {
		return nil
	}

	var toAdd Item
	if i == 0 {
		toAdd.Type = Insertion
	} else if j == 0 {
		toAdd.Type = Deletion
	} else if t.a[i-1].Text == t.b[j-1].Text {
		toAdd.Type = Unchanged
	} else if t.getLength(i, j-1) > t.getLength(i-1, j) {
		toAdd.Type = Insertion
	} else {
		toAdd.Type = Deletion
	}

	switch toAdd.Type {
	case Insertion:
		toAdd.Line = t.b[j-1]
		j--
	case Unchanged:
		toAdd.Line = t.a[i-1]
		i--
		j--
	case Deletion:
		toAdd.Line = t.a[i-1]
		i--
	}

	return append(t.recursiveDiff(i, j), toAdd)
}
