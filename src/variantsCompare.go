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
	out.WriteString("Patient,Chr,Start,End,REF,ALT,Name,Coverage,ReferenceReads,VariantReads,AlleleFrequency,A,T,G,C\n")
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
	fmt.Printf("\n\tFound %d new variants.\n", v.neu)
	fmt.Printf("\tVerified %d of %d variants.\n", matched, v.total)
	fmt.Printf("\tIdentified %d variants with an average coverage of %s.\n", total, v.getAverage(total, matched))
}

func (v *variants) getAlleles(ref string, row []string) map[string]int {
	// Extracts alternate alleles from row
	ret := make(map[string]int)
	if len(row) >= 5 {
		for _, i := range row[4:] {
			s := strings.Split(i, ":")
			b := strings.ToUpper(s[0])
			if b != "=" {
				count, err := strconv.Atoi(s[1])
				if err == nil && count > 0 {
					ret[b] = count
				}
			}
		}
	}
	return ret
}

func (v *variants) examineBamReadcount(id string, h map[string]int, row []string) {
	// Compares variant from bam-readcount to v.vars
	chr := v.setChromosome(row[0])
	if _, ex := v.vars[id][chr]; ex == true {
		pos := setCoordinate(row[1])
		ref := row[2]
		bases := v.getAlleles(ref, row)
		if len(bases) >= 1 {
			for k := range bases {
				if k != ref {
					v.neu++
					for _, i := range v.vars[id][chr] {
						if i.equals(pos, ref, k) {
							// Equals method records hits if true
							i.addCounts(bases)
							break
						}
					}
				}
			}
		}
	}
}

func (v *variants) getAlleleFrequencyFromVCF(s string) string {
	// Subsets variant allele frequency from vcf info section
	ret := "NA"
	s = s[strings.Index(s, "AF="):]
	s = s[strings.Index(s, "=")+1:]
	if strings.Contains(s, ";") {
		s = s[:strings.Index(s, ";")]
	}
	if _, err := strconv.ParseFloat(s, 64); err == nil {
		ret = s
	}
	return ret
}

func (v *variants) examineVCF(id string, h map[string]int, row []string) {
	// Compares variant from vcf to v.vars
	chr := v.setChromosome(row[h["CHROM"]])
	if _, ex := v.vars[id][chr]; ex == true {
		pos := setCoordinate(row[h["POS"]])
		ref := row[h["REF"]]
		alt := row[h["ALT"]]
		v.neu++
		for _, i := range v.vars[id][chr] {
			if i.equals(pos, ref, alt) {
				// Equals method will record matches
				freq := v.getAlleleFrequencyFromVCF(row[h["INFO"]])
				i.appendFrequency(freq)
				break
			}
		}
	}
}

func (v *variants) readVCF(wg *sync.WaitGroup, id, infile string) {
	// Reads in infile as a dictionary stored by chromosome
	var h map[string]int
	var d string
	head := true
	vcf := false
	defer wg.Done()
	f := iotools.OpenFile(infile)
	defer f.Close()
	input := iotools.GetScanner(f)
	for input.Scan() {
		line := strings.TrimSpace(string(input.Text()))
		if head == true {
			if line[0] == '#' && line[1] != '#' {
				// Skip over remaining vcf header
				line = strings.Replace(line, "#", "", 1)
				d = iotools.GetDelim(line)
				h = iotools.GetHeader(strings.Split(line, d))
				head = false
			} else if strings.Contains(line, "##") && vcf == false {
				// Record file format
				vcf = true
			} else if vcf == false {
				// Read first line of bam-readcount output
				d = iotools.GetDelim(line)
				v.examineBamReadcount(id, h, strings.Split(line, d))
				head = false
			}
		} else if vcf == true {
			v.examineVCF(id, h, strings.Split(line, d))
		} else {
			v.examineBamReadcount(id, h, strings.Split(line, d))
		}
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
