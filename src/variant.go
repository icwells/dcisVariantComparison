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
	id      string
	name    string
	shared  string
	chr     string
	start   int
	end     int
	ref     string
	alt     string
	matches int
	normal  *counts
	tumor   *counts
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
	return v
}

func (v *variant) String() string {
	// Returns formatted string for printing
	n := v.normal.String()
	t := v.tumor.String()
	return fmt.Sprintf("%s,%s,%s,%d,%d,%s,%s,%s,%d,%s,%s\n", v.id, v.shared, v.chr, v.start, v.end, v.ref, v.alt, v.name, v.matches, t, n)
}

func (v *variant) evaluate(wg *sync.WaitGroup, normal bool, row map[int]*variant) {
	// Assembles variant from bam-readcount data
	defer wg.Done()
	a := new(variant)
	a.tumor = newCounts()
	for i := 0; i <= v.end-v.start; i++ {
		// Attempt to contruct reference and alternate variants from readcount data
		idx := i + v.start
		r, ex := row[idx]
		if ex == false {
			break
		}
		a.ref += r.ref
		a.alt += r.tumor.getAlternate(r.ref)
		a.tumor.addCounts(r.ref, r.tumor.bases)
	}
	if a.ref == v.ref && a.alt == v.alt {
		v.matches++
		if normal == true {
			v.normal.addCounts(v.ref, a.tumor.bases)
		} else {
			v.tumor.addCounts(v.ref, a.tumor.bases)
		}
	}
}
