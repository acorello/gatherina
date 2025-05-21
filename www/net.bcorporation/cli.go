package main

import (
	"flag"
	"math"
)

var cli struct {
	maxPages   int
	outputPath string
}

func init() {
	flag.IntVar(&cli.maxPages, "max-pages", math.MaxInt, "maximum number of pages to fetch")
	flag.StringVar(&cli.outputPath, "output-path", "", "path to the file where output will be written (output format is JSONL); if omitted a temporary file is created")
	flag.Parse()
}
