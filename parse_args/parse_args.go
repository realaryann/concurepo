package parse_args

import (
	"flag"
)

func Parse_args() (*bool, *string, *uint, *string) {
	ver := flag.Bool("version", false, "Current version number")
	ofile := flag.String("o", "", "Post scraping output file")
	spawn := flag.Uint("spawn", 1, "Number of concurrent goroutines to scrape, positive integer")
	flags := flag.String("flag", "", "Comma separated flags to filter the results by")
	
	flag.Parse()

	return ver, ofile, spawn, flags
}

