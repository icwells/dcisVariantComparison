// Tests variant struct

package main

import (
	"testing"
)

func TestGetAlleleFrequencyFromVCF(t *testing.T) {
	var v variants
	info := []string{"AB=0.454545;ABP=3.20771;AC=1;AF=0.5", "AB=0;ABP=0;AC=2;AF=1", "AB=0.454545;AF=0.5;ABP=3.20771;AC=1", "AB=0;ABP=0;AF=0.25;AC=2"}
	exp := []string{"0.5", "1", "0.5", "0.25"}
	for idx, i := range info {
		act := v.getAlleleFrequencyFromVCF(i)
		if act != exp[idx] {
			t.Errorf("Actual allele frequency %s does not equal expected: %s.", act, exp[idx])
		}
	}
}

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
	var v variant
	cases := map[string]string{"a ": "A", " . ": "-", "TGCTGT": "TGCTGT", " ATcg": "ATCG"}
	for k, val := range cases {
		act := v.setAllele(k)
		if act != val {
			t.Errorf("Actual allele %s does not equal expected: %s", act, val)
		}
	}
}

func getVariants() map[string][]*variant {
	// Returns map of variants for testing
	ret := make(map[string][]*variant)
	pid := "DCIS_1 "
	ret["1"] = []*variant{newVariant(pid, "1", "100.0", "200.0", "A", "t", "NA", "A")}
	ret["1"] = append(ret["1"], newVariant(pid, "1", "1025", "1119", "G", "-", "NA", "A"))
	ret["2"] = []*variant{newVariant(pid, "2", "25006", "25124", "C", "G", "NA", "A")}
	ret["X"] = []*variant{newVariant(pid, "X", "90045", "90157.5", ".", "A", "NA", "A")}
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
