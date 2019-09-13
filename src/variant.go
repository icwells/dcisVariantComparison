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
	chr     string
	start   int
	end     int
	ref     string
	rcount  int
	alt     string
	acount  int
	freq    string
	matches int
}

func (v *variant) setAllele(val string) string {
	// Makes sure alleles are in the same format
	ret := strings.ToUpper(strings.TrimSpace(val))
	if ret == "." {
		ret = "-"
	}
	return ret
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
	v.freq = "NA"
	return v
}

func (v *variant) String() string {
	// Returns formatted string for printing
	return fmt.Sprintf("%s,%s,%d,%d,%s,%s,%s,%d,%d,%d,%s\n", v.id, v.chr, v.start, v.end, v.ref, v.alt, v.name, v.matches, v.rcount, v.acount, v.freq)
}

func (v *variant) addCounts(a, r int) {
	// Adds number of reads with ref/alt alleles
	v.acount = a
	v.rcount = r
}

func (v *variant) appendFrequency(f string) {
	// Stores variant allele frequency
	if v.freq == "NA" {
		v.freq = f
	}
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
