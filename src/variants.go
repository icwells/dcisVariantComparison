// This script defines a struct for storing a covb/nan output

package main

import (
	"fmt"
	"github.com/icwells/go-tools/iotools"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
)

type variants struct {
	infile  string
	vcfs    map[string][]string
	files   int
	outfile string
	vars    map[string]map[string][]*variant
}

func (v *variants) getSampleID(filename string) string {
	// Attempts to resolve result of getSampleName with variants keys
	n := strings.Split(iotools.GetFileName(filename), "-")
	// Add DCIS and Sample name
	ret := n[0] + n[1]
	if _, err := strconv.Atoi(string(n[2][1])); err == nil {
		// Add alpha-numeric codes
		ret = fmt.Sprintf("%s_%s", ret, n[2])
	}
	if _, ex := v.vars[ret]; ex == false {
		ret = strings.Split(ret, "_")[0]
		ret = strings.Replace(ret, "0", "", -1)
		if _, ex := v.vars[ret]; ex == false {
			ret = ""
		}
	}
	return ret
}

func (v *variants) setVCFs(vcfs string) {
	// Reads in map of vcf files
	v.vcfs = make(map[string][]string)
	files, err := filepath.Glob(path.Join(vcfs, "*.vcf"))
	if err == nil {
		for _, i := range files {
			n := v.getSampleID(i)
			if n != "" {
				v.vcfs[n] = append(v.vcfs[n], i)
				v.files++
			}
		}
	}
	if v.files < 1 {
		fmt.Print("\n\t[Error] No matching vcfs files found. Exiting.\n\n")
		os.Exit(1)
	}
}

func (v *variants) setChromosome(val string) string {
	// Removes decimal from chromosome number
	if strings.Contains(val, ".0") {
		val = strings.Split(val, ".")[0]
	}
	return strings.TrimSpace(val)
}

func (v *variants) setVariant(h map[string]int, row []string) {
	// Reads variant from row and stores in map
	id := strings.TrimSpace(row[h["Patient"]])
	chr := v.setChromosome(row[h["Chr"]])
	start := row[h["Start"]]
	end := row[h["End"]]
	ref := row[h["REF"]]
	alt := row[h["ALT"]]
	name := row[h["Name"]]
	if _, ex := v.vars[id]; ex == false {
		v.vars[id] = make(map[string][]*variant)
	}
	newvar := newVariant(id, chr, start, end, ref, alt, name)
	v.vars[id][chr] = append(v.vars[id][chr], newvar)
}

func (v *variants) setVariants() {
	// Reads interval file into struct
	var h map[string]int
	var d string
	first := true
	v.vars = make(map[string]map[string][]*variant)
	fmt.Println("\n\tReading input variants file...")
	f := iotools.OpenFile(v.infile)
	defer f.Close()
	input := iotools.GetScanner(f)
	for input.Scan() {
		line := string(input.Text())
		if first == false {
			v.setVariant(h, strings.Split(line, d))
		} else {
			d = iotools.GetDelim(line)
			h = iotools.GetHeader(strings.Split(line, d))
			first = false
		}
	}
}

func newVariants(infile, vcfs, outfile string) *variants {
	// Initializes new variants struct and reads in input file
	v := new(variants)
	v.infile = infile
	v.outfile = outfile
	v.setVariants()
	v.setVCFs(vcfs)
	return v
}
