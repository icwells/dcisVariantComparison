// Stores variant data in struct

package main

import (
	"fmt"
	"strconv"
	"strings"
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

func (v *variant) setAllele(val string) string {
	// Makes sure alleles are in the same format
	ret := strings.ToUpper(strings.TrimSpace(val))
	if ret == "." {
		ret = "-"
	}
	return ret
}

func newVariant(id, chr, start, end, ref, alt, name, shared string) *variant {
	v := new(variant)
	v.id = id
	v.chr = chr
	v.start = setCoordinate(start)
	v.end = setCoordinate(end)
	v.ref = v.setAllele(ref)
	v.alt = v.setAllele(alt)
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

func (v *variant) evaluate(normal bool, pos int, ref string, bases map[string]int) bool {
	// Returns true if pos is inside v.start/end and ref == v.ref
	ret := false
	ref = v.setAllele(ref)
	if v.start <= pos && v.end >= pos && ref == v.ref {
		// Only proceed for potential match
		for k := range bases {
			k = v.setAllele(k)
			if k == v.alt {
				v.matches++
				ret = true
				if normal == true {
					v.normal.addCounts(v.ref, bases)
				} else {
					v.tumor.addCounts(v.ref, bases)
				}
			}
		}
	}
	return ret
}
