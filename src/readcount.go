// Defines readcount struct

package main

import (
	""
)

type readcount struct {
	chr string
	pos	int
	ref	string
	bases map[string]int
}

func newReadCount(chr, ref string, pos int, bases map[string]int) *readcount {
	// Returns initialized readcount struct
	r := new(readcount)
	r.chr = chr
	r.pos = po
	r.ref = ref
	r.bases = bases
	return r
}

