package parse_args

import (
	"flag"
)

func Parse_args() (*bool, *string, *uint) {
	ver := flag.Bool("version", false, "Current version number")
	ofile := flag.String("o", "a.txt", "Post scraping output file")
	spawn := flag.Uint("spawn", 1, "Number of concurrent goroutines to scrape, positive integer")
	
	flag.Parse()

	return ver, ofile, spawn
}

