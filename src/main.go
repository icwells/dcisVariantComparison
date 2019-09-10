// Compares vcf from ampliseq samples to variants summary file

package main

import (
	"fmt"
	"github.com/icwells/go-tools/iotools"
	"gopkg.in/alecthomas/kingpin.v2"
	"os"
	"time"
)

var (
	app     = kingpin.New("dcisVariantComparison", "Compares vcfs from ampliseq samples to variants summary file.")
	infile  = kingpin.Flag("infile", "Path to input variants file.").Required().Short('i').String()
	vcfs    = kingpin.Flag("vcfs", "Path to directory of input vcf files.").Required().Short('v').String()
	outfile = kingpin.Flag("outfile", "Path to output file.").Required().Short('o').String()
)

func checkArgs() {
	// Exists if input is missing
	for _, i := range []string{*infile, *vcfs} {
		if iotools.Exists(i) == false {
			fmt.Printf("\n\t[Error] %s not found. Exiting.\n", i)
			os.Exit(1)
		}
	}
}

func main() {
	start := time.Now()
	kingpin.Parse()
	checkArgs()
	v := newVariants(*infile, *vcfs, *outfile)
	v.compareVariants()
	v.writeOutput()
	fmt.Printf("\tFinished. Runtime: %s\n\n", time.Since(start))
}
