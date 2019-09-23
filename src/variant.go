// Stores variant data in struct

package main

import (
	"fmt"
	"strconv"
	"strings"
	"sync"
)

func setCoordinate(n string) int {
	// Removes decimal from coordinate number
	n = strings.Replace(n, ",", "", -1)
	if strings.Contains(n, ".") {
		n = strings.Split(n, ".")[0]
	}
	ret, err := strconv.Atoi(n)
	if err != nil {
		ret = -1
	}
	return ret
}

func setAllele(val string) string {
	// Makes sure alleles are in the same format
	ret := strings.ToUpper(strings.TrimSpace(val))
	if ret == "." || ret == "=" {
		ret = "-"
	}
	return ret
}

type variant struct {
	id        string
	name      string
	shared    string
	chr       string
	start     int
	end       int
	ref       string
	alt       string
	deletion  bool
	insertion bool
	matches   int
	normal    *counts
	tumor     *counts
}

func newReadCount(chr, ref string, pos int, bases map[string]int) *variant {
	// Returns initialized variant struct for bam-readcount data
	v := new(variant)
	v.chr = chr
	v.ref = ref
	v.start = pos
	v.tumor = newCounts()
	v.tumor.addCounts(v.ref, bases)
	return v
}

func (v *variant) setType() {
	// Determines mutation type
	if v.ref == "-" {
		v.insertion = true
	} else if v.alt == "-" {
		v.deletion = true
	}
}

func newVariant(id, chr, start, end, ref, alt, name, shared string) *variant {
	v := new(variant)
	v.id = id
	v.chr = chr
	v.start = setCoordinate(start)
	v.end = setCoordinate(end)
	v.ref = setAllele(ref)
	v.alt = setAllele(alt)
	v.name = strings.TrimSpace(name)
	v.shared = strings.TrimSpace(shared)
	v.normal = newCounts()
	v.tumor = newCounts()
	v.setType()
	return v
}

func (v *variant) String() string {
	// Returns formatted string for printing
	n := v.normal.String()
	t := v.tumor.String()
	return fmt.Sprintf("%s,%s,%s,%d,%d,%s,%s,%s,%d,%s,%s\n", v.id, v.shared, v.chr, v.start, v.end, v.ref, v.alt, v.name, v.matches, t, n)
}

func (v *variant) findInsertion(row map[int]*variant) (bool, map[string]int) {
	// Assembles variant from bam-readcount data
	var found bool
	a := new(variant)
	a.tumor = newCounts()
	for i := 0; i <= len(v.alt); i++ {
		// Attempt to contruct reference and alternate variants from readcount data
		r, ex := row[i+v.start]
		if ex == false {
			return false, a.tumor.bases
		}
		a.ref += r.ref
		a.alt += r.tumor.getAlternate(r.ref)
		a.tumor.addCounts(r.ref, r.tumor.bases)
	}
	if a.ref != v.alt && a.alt == v.alt {
		found = true
	}
	return found, a.tumor.bases
}

func (v *variant) findDeletion(row map[int]*variant) (bool, map[string]int) {
	// Assembles variant from bam-readcount data
	var found bool
	a := new(variant)
	a.tumor = newCounts()
	for i := 0; i <= len(v.ref); i++ {
		// Attempt to contruct reference and alternate variants from readcount data
		r, ex := row[i+v.start]
		if ex == false {
			return false, a.tumor.bases
		}
		a.ref += r.ref
		a.alt += r.tumor.getAlternate(r.ref)
		a.tumor.addCounts(r.ref, r.tumor.bases)
	}
	if a.ref == v.ref && a.alt != v.ref {
		found = true
	}
	return found, a.tumor.bases
}

func (v *variant) findSNP(row map[int]*variant) (bool, map[string]int) {
	// Finds matching SNPs
	var found bool
	var ret map[string]int
	a, ex := row[v.start]
	if ex == true && a.ref == v.ref && a.alt == v.alt {
		found = true
		ret = a.tumor.bases
	}
	return found, ret
}

func (v *variant) evaluate(wg *sync.WaitGroup, normal bool, row map[int]*variant) {
	// Identifies matching variants from bam-readcount data
	defer wg.Done()
	var bases map[string]int
	var found bool
	if v.insertion {

	} else if v.deletion {

	} else {
		found, bases = v.findSNP(row)
	}
	if found {
		v.matches++
		if normal == true {
			v.normal.addCounts(v.ref, bases)
		} else {
			v.tumor.addCounts(v.ref, bases)
		}
	}
}
