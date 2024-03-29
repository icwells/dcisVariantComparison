// This script defines a struct for storing a covb/nan output

package main

import (
	"fmt"
	"github.com/icwells/go-tools/iotools"
	"strconv"
	"strings"
	"sync"
)

func (v *variants) getAverage(total, matched int) string {
	// Returns formatted average coverage
	ret := strconv.FormatFloat(float64(total)/float64(matched), 'f', 4, 64)
	ret = strings.TrimRight(ret, "0")
	if ret[len(ret)-1] == '.' {
		// Remove trailing decimal
		ret = strings.Replace(ret, ".", "", 1)
	}
	return ret
}

func (v *variants) writeOutput() {
	// Write results to output file
	var matched, total int
	fmt.Println("\tWriting results to file...")
	out := iotools.CreateFile(v.outfile)
	defer out.Close()
	head := "Patient,Shared,Chr,Start,End,REF,ALT,Name,Coverage,"
	head += "TReferenceReads,TVariantReads,TAlleleFrequency,A,T,G,C,"
	head += "NReferenceReads,NVariantReads,NAlleleFrequency,A,T,G,C\n"
	out.WriteString(head)
	for _, val := range v.vars {
		for _, v := range val {
			for _, i := range v {
				// Write string and record number of matches
				out.WriteString(i.String())
				total += i.matches
				if i.matches > 0 {
					matched++
				}
			}
		}
	}
	fmt.Printf("\n\tVerified %d of %d variants.\n", matched, v.total)
	fmt.Printf("\tIdentified %d variants with an average coverage of %s.\n", total, v.getAverage(total, matched))
}

func (v *variants) identifyVariants(normal bool, ids []string, target map[string]map[int]*variant) {
	// Compares variants map against target map
	var wg sync.WaitGroup
	for _, id := range ids {
		for k := range v.vars[id] {
			row, ex := target[k]
			if ex {
				for idx := range v.vars[id][k] {
					wg.Add(1)
					go v.vars[id][k][idx].evaluate(&wg, normal, row)
				}
			}
		}
		wg.Wait()
	}
}

func (v *variants) getAlleles(ref string, row []string) map[string]int {
	// Extracts alternate alleles from row
	ret := make(map[string]int)
	if len(row) >= 5 {
		for _, i := range row[4:] {
			s := strings.Split(i, ":")
			b := setAllele(s[0])
			count, err := strconv.Atoi(s[1])
			if err == nil && count > 0 {
				ret[b] = count
			}
		}
	}
	return ret
}

func (v *variants) examineBamReadcount(id string, row []string) *variant {
	// Compares variant from bam-readcount to v.vars
	ret := new(variant)
	chr := v.setChromosome(row[0])
	if _, ex := v.vars[id][chr]; ex == true {
		pos := setCoordinate(row[1])
		ref := row[2]
		bases := v.getAlleles(ref, row)
		if len(bases) >= 1 {
			ret = newReadCount(chr, ref, pos, bases)
		}
	}
	return ret
}

func (v *variants) getNormalStatus(infile string) bool {
	// Returns true if infile is of normal tissue
	ret := false
	f := strings.ToLower(infile)
	if strings.Contains(f, "node") || strings.Contains(f, "benign") {
		ret = true
	}
	return ret
}

func (v *variants) findIDs(normal bool, id string) []string {
	// Finds all mathcin sample ids
	var ret []string
	if normal {
		id = strings.Split(id, "_")[0]
		for k := range v.vars {
			if strings.Contains(k, id) {
				ret = append(ret, k)
			}
		}
	} else if _, ex := v.vars[id]; ex {
		ret = append(ret, id)
	}
	return ret
}

func (v *variants) readVCF(wg *sync.WaitGroup, id, infile string) {
	// Reads in infile as a dictionary stored by chromosome
	var d string
	defer wg.Done()
	first := true
	target := make(map[string]map[int]*variant)
	normal := v.getNormalStatus(infile)
	ids := v.findIDs(normal, id)
	if len(ids) > 0 {
		f := iotools.OpenFile(infile)
		defer f.Close()
		input := iotools.GetScanner(f)
		for input.Scan() {
			line := strings.TrimSpace(string(input.Text()))
			if first == true {
				d = iotools.GetDelim(line)
				first = false
			}
			rc := v.examineBamReadcount(id, strings.Split(line, d))
			if len(rc.chr) > 0 {
				if _, ex := target[rc.chr]; ex == false {
					target[rc.chr] = make(map[int]*variant)
				}
				target[rc.chr][rc.start] = rc
			}
		}
		v.identifyVariants(normal, ids, target)
	}
}

func (v *variants) compareVariants() {
	// Compares input vcfs against variants file
	count := 1
	var wg sync.WaitGroup
	for k, vals := range v.vcfs {
		for _, i := range vals {
			wg.Add(1)
			go v.readVCF(&wg, k, i)
			fmt.Printf("\r\tDispatched %d of %d vcf files.", count, v.files)
			count++
		}
	}
	fmt.Print("\n")
	wg.Wait()
}
