package parse_args

import (
	"flag"
)

func Parse_args() (*bool, *bool, *bool, *string, *uint, *string) {
	ver := flag.Bool("version", false, "Current version number")
	intern := flag.Bool("intern", false, "Look for internship positions")
	fulltime := flag.Bool("fulltime", false, "Look for full time positions")
	ofile := flag.String("o", "", "Post scraping output file")
	limit := flag.Uint("limit", 15, "Number of positions to scrape. Default = 15")
	flags := flag.String("flag", "", "Comma separated flags to filter the results by")
	
	flag.Parse()

	return ver, intern, fulltime, ofile,  limit, flags
}

