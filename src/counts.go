// Defines count struct to record counts for normal/tumor files

package main

import (
	"fmt"
	"strconv"
)

type counts struct {
	ref   int
	alt   int
	freq  string
	bases map[string]int
}

func newCounts() *counts {
	// Initializes struct
	c := new(counts)
	c.freq = "NA"
	c.bases = map[string]int{"A": 0, "T": 0, "G": 0, "C": 0, "-": 0}
	return c
}

func (c *counts) String() string {
	// Returns formatted string
	return fmt.Sprintf("%d,%d,%s,%d,%d,%d,%d", c.ref, c.alt, c.freq, c.bases["A"], c.bases["T"], c.bases["G"], c.bases["C"])
}

func (c *counts) getAlternate(ref string) string {
	// Returns most common alternate base
	ret := ref
	max := 0
	for k, v := range c.bases {
		if k != ref && v > max {
			ret = k
			max = v
		}
	}
	return ret
}

func (c *counts) calculateAlleleFrequency() {
	// Calculates variant allele frequency from bam-readcount data
	if c.ref > 0 && c.alt > 0 {
		f := float64(c.alt) / float64(c.alt+c.ref)
		c.freq = strconv.FormatFloat(f, 'f', 4, 64)
	}
}

func (c *counts) addCounts(ref string, bases map[string]int) {
	// Adds number of reads with ref/alt alleles
	for k, val := range bases {
		if k == ref {
			c.ref += val
		} else {
			c.alt += val
		}
		c.bases[k] += val
	}
	c.calculateAlleleFrequency()
}
