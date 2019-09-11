// Tests variants struct

package main

import (
	"testing"
)

func getVariants() map[string][]*variant {
	// Returns map of variants for testing
	ret := make(map[string][]*variant)
	pid := "DCIS_1 "
	ret["1"] = []*variant{newVariant(pid, "1", "100.0", "200.0", "A", "t", "NA")}
	ret["1"] = append(ret["1"], newVariant(pid, "1", "1025", "1119", "G", "-", "NA"))
	ret["2"] = []*variant{newVariant(pid, "2", "25006", "25124", "C", "G", "NA")}
	ret["X"] = []*variant{newVariant(pid, "X", "90045", "90157.5", ".", "A", "NA")}
	return ret
}

func TestEquals(t *testing.T) {
	vars := getVariants()
	cases := []struct {
		pid string
		pos int
		ref string
		alt string
		exp bool
		mat int
	}{
		{"1", 155, "A", "T", true, 1},
		{"1", 1075, "G", "C", false, 0},
		{"2", 8875, "C", "G", false, 0},
		{"X", 90065, "-", "A", true, 1},
	}
	for _, i := range cases {
		actual := false
		match := 0
		for _, v := range vars[i.pid] {
			res := v.equals(i.pos, i.ref, i.alt)
			if res == true {
				actual = res
				match = v.matches
			}
		}
		if actual != i.exp {
			t.Errorf("Actual result %v does not equal expect at position %d", actual, i.pos)
		} else if match != i.mat {
			t.Errorf("Actual number of matches %d does not equal expected %d where position is %d", match, i.mat, i.pos)
		}
	}
}
