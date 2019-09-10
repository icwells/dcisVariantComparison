// Stores variant data in struct

package main

import (
	"fmt"
	"strconv"
	"strings"
)

func setCoordinate(n string) int {
	// Removes decimal from coordinate number
	if strings.Contains(n, ".") {
		n = strings.Split(n, ".")[0]
	}
	ret, err := strconv.Atoi(n)
	if err != nil {
		ret = -1
	}
	return ret
}

func (v *variant) setAllele(val string) string {
	// Makes sure alleles are in the same format
	ret := strings.ToUpper(strings.TrimSpace(val))
	if ret == "." {
		ret = "-"
	}
	return ret
}

type variant struct {
	id      string
	name    string
	chr     string
	start   int
	end     int
	ref     string
	alt     string
	matches int
}

func newVariant(id, chr, start, end, ref, alt, name string) *variant {
	v := new(variant)
	v.id = id
	v.chr = chr
	v.start = setCoordinate(start)
	v.end = setCoordinate(end)
	v.ref = v.setAllele(ref)
	v.alt = v.setAllele(alt)
	v.name = strings.TrimSpace(name)
	return v
}

func (v *variant) String() string {
	// Returns formatted string for printing
	return fmt.Sprintf("%s,%s,%d,%d,%s,%d\n", v.id, v.chr, v.start, v.end, v.name, v.matches)
}

func (v *variant) equals(pos int, ref, alt string) bool {
	// Returns true if pos is inside v.start/end and ref == v.ref
	ref = v.setAllele(ref)
	alt = v.setAllele(alt)
	if v.start <= pos && v.end >= pos && ref == v.ref && alt == v.alt {
		v.matches++
		return true
	} else {
		return false
	}
}
