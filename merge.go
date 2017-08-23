package diff

import "strings"

func dropWhile(d []Item, fn func(Item) bool) []Item {
	for i, x := range d {
		if !fn(x) {
			return d[i:]
		}
	}
	return nil
}

// Merge takes a slice of original lines, as well as any number of other
// "versions". A Patience diff from the original is computed for each
// version, the changesets are merged together, and the resulting lines are
// returned.
//
// The algorithm used for merging the resulting diffs, is:
//	while there are non-empty diffs,
//	    if the first element in each diff is Unchanged, remove a line from the
//	      front of the original set of lines, and append it to the result lines,
//	      removing the Unchanged first element from each diff
//	    else,
//	        if a diff begins with an insertion,
//	            while the diff's first item is an insertion, remove it and add
//	              the line to a temporary block
//	            join the temporary block with newlines, and if the last line in
//	              the results does not have the temporary block as a prefix,
//	              append it to the results
//	        else if a diff begins with a deletion,
//	            if the first line of the original set of lines matches the deletion,
//	                remove it
//	            for each other non-empty diff,
//	                remove the first element if it is Unchanged and matches the deletion
func Merge(a []string, others ...[]string) []string {
	diffs := make([][]Item, len(others))
	for i, l := range others {
		diffs[i] = Patience(a, l)
	}

	merged := make([]string, 0)
	// Continue until all diffs have been exhausted
	for empty := false; !empty; {
		changed := false
		empty = true
		for i, d := range diffs {
			if len(d) == 0 {
				continue
			}
			empty = false
			switch d[0].Type {
			case Insertion:
				changed = true
				// accumulate blocks of added lines
				toInsert := make([]string, 0)
				diffs[i] = dropWhile(d, func(c Item) (drop bool) {
					if drop = (c.Type == Insertion); drop {
						toInsert = append(toInsert, c.Text)
					}
					return
				})
				// skip lines that were just added (completely or partially)
				// by another diff
				joined := strings.Join(toInsert, "\n")
				if len(merged) == 0 || !strings.HasPrefix(merged[len(merged)-1], joined) {
					merged = append(merged, joined)
				}
			case Deletion:
				changed = true
				diffs[i] = dropWhile(d, func(c Item) (drop bool) {
					if drop = (c.Type == Deletion); drop {
						if len(a) > 0 && a[0] == c.Text {
							a = a[1:]
						}
						// delete lines from other diffs
						for j := 0; j < len(diffs); j++ {
							if i == j {
								continue
							}
							diffs[j] = dropWhile(diffs[j], func(x Item) bool {
								return x.Type == Unchanged && x.Text == c.Text
							})
						}
					}
					return
				})
			}
		}
		if !changed {
			// "pop" the unchanged line from all diffs
			for i, d := range diffs {
				if len(d) > 0 {
					diffs[i] = d[1:]
				}
			}
			if len(a) > 0 {
				merged = append(merged, a[0])
				a = a[1:]
			}
		}
	}
	return merged
}
