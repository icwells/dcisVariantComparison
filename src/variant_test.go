// Tests variant struct

package main

import (
	"sync"
	"testing"
)

func TestGetSampleID(t *testing.T) {
	var v variants
	v.vars = make(map[string]map[string][]*variant)
	v.vars["DCIS64"] = make(map[string][]*variant)
	v.vars["DCIS267"] = make(map[string][]*variant)
	v.vars["DCIS168_C4"] = make(map[string][]*variant)
	cases := map[string]string{
		"ampliseq2/vcfs/DCIS-064-A61.vcf":      "DCIS64",
		"/ampliseq2/vcfs/DCIS-064-A81-inv.vcf": "DCIS64",
		"ampliseq2/vcfs/DCIS-267-B1-node.vcf":  "DCIS267",
		"ampliseq2/vcfs/DCIS-168-C4-inv.vcf":   "DCIS168_C4",
		"ampliseq2/vcfs/DCIS-300-C4-inv.vcf":   "",
	}
	for k, val := range cases {
		act := v.getSampleID(k)
		if act != val {
			t.Errorf("Actual sample ID %s does not equal expected: %s", act, val)
		}
	}
}

func TestSetChromosome(t *testing.T) {
	var v variants
	cases := map[string]string{"1.0 ": "1", " 2 ": "2", "GL001.0": "GL001", " X": "X"}
	for k, val := range cases {
		act := v.setChromosome(k)
		if act != val {
			t.Errorf("Actual chromosome %s does not equal expected: %s", act, val)
		}
	}
}

func TestSetCoordinate(t *testing.T) {
	cases := []struct {
		input string
		exp   int
	}{
		{"100.0", 100},
		{"afd", -1},
		{"54.a", 54},
		{"21,900", 21900},
	}
	for _, i := range cases {
		act := setCoordinate(i.input)
		if act != i.exp {
			t.Errorf("Actual coordinate %d does not equal expected: %d", act, i.exp)
		}
	}
}

func TestSetAllele(t *testing.T) {
	cases := map[string]string{"a ": "A", " . ": "-", "TGCTGT": "TGCTGT", " ATcg": "ATCG"}
	for k, val := range cases {
		act := setAllele(k)
		if act != val {
			t.Errorf("Actual allele %s does not equal expected: %s", act, val)
		}
	}
}

func getVariants() map[string][]*variant {
	// Returns map of expected variants for testing (store expected true/false for id)
	ret := make(map[string][]*variant)
	ret["1"] = []*variant{newVariant("true", "1", "100.0", "100.0", "A", "t", "NA", "A")}
	ret["1"] = append(ret["1"], newVariant("false", "1", "1025", "1119", "G", "-", "NA", "A"))
	ret["2"] = []*variant{newVariant("true", "2", "25006", "25009", "CTCA", "GCAT", "NA", "A")}
	ret["X"] = []*variant{newVariant("true", "X", "90045", "90045.5", ".", "A", "NA", "A")}
	return ret
}

func getReadCounts() map[string]map[int]*variant {
	// Returns test cases for variant.evaluate
	ret := make(map[string]map[int]*variant)
	ret["1"] = make(map[int]*variant)
	ret["2"] = make(map[int]*variant)
	ret["X"] = make(map[int]*variant)
	ret["1"][155] = newReadCount("1", "A", 155, map[string]int{"A": 3, "T": 9, "G": 0, "C": 0})
	ret["1"][1075] = newReadCount("1", "G", 1075, map[string]int{"A": 10, "T": 0, "G": 0, "C": 15})
	ret["2"][25006] = newReadCount("2", "C", 25006, map[string]int{"A": 1, "T": 0, "G": 12, "C": 1})
	ret["2"][25007] = newReadCount("2", "T", 25007, map[string]int{"A": 1, "T": 0, "G": 2, "C": 9})
	ret["2"][25008] = newReadCount("2", "C", 25008, map[string]int{"A": 11, "T": 0, "G": 2, "C": 0})
	ret["2"][25009] = newReadCount("2", "A", 25009, map[string]int{"A": 1, "T": 6, "G": 2, "C": 0})
	ret["X"][90065] = newReadCount("X", "-", 90065, map[string]int{"A": 10, "T": 0, "G": 0, "C": 6})
	return ret
}

func TestEvaluate(t *testing.T) {
	var wg sync.WaitGroup
	cases := getVariants()
	vars := getReadCounts()
	for k, v := range cases {
		for _, i := range v {
			wg.Add(1)
			match := i.matches
			i.evaluate(&wg, false, vars[k])
			wg.Wait()
			if i.id == "true" && i.matches == match {
				t.Errorf("No match where position is %s:%d-%d", i.chr, i.start, i.end)
			} else if i.matches > match {
				t.Errorf("False match where position is %s:%d-%d", i.chr, i.start, i.end)
			}
		}
	}
}
