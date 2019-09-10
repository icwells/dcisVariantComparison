// This script defines a struct for storing a covb/nan output

package main

import (
	"fmt"
	"github.com/icwells/go-tools/iotools"
	"strings"
	"sync"
)

func (v *variants) writeOutput() {
	// Write results to output file
	fmt.Println("\tWriting results to file...")
	out := iotools.CreateFile(v.outfile)
	defer out.Close()
	out.WriteString("Patient,Chr,Start,End,Name,Coverage\n")
	for _, val := range v.vars {
		for _, v := range val {
			for _, i := range v {
				out.WriteString(i.String())
			}
		}
	}
}

func (v *variants) examineVariant(id string, h map[string]int, row []string) {
	// Compares variant to v.variants
	chr := v.setChromosome(row[h["CHROM"]])
	pos := setCoordinate(row[h["POS"]])
	ref := row[h["REF"]]
	alt := row[h["ALT"]]
	if _, ex := v.vars[id][chr]; ex == true {
		for _, i := range v.vars[id][chr] {
			if i.equals(pos, ref, alt) {
				// Equals method records hits if true
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
	defer wg.Done()
	f := iotools.OpenFile(infile)
	defer f.Close()
	input := iotools.GetScanner(f)
	for input.Scan() {
		line := string(input.Text())
		if head == false {
			v.examineVariant(id, h, strings.Split(line, d))
		} else if line[0] == '#' && line[1] != '#' {
			// Skip over vcf header
			line = strings.Replace(line, "#", "", 1)
			d = iotools.GetDelim(line)
			h = iotools.GetHeader(strings.Split(line, d))
			head = false
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
