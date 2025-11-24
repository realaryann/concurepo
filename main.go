package main

import (
	"log"
	"sync"
	"os"
	"concurjob/parse_args"
	"concurjob/version"
	"concurjob/scrape"
	"github.com/jedib0t/go-pretty/v6/table"
	"strings"
)

func main() {
	ver, intern, fulltime, company, ofile, limit, flags  := parse_args.Parse_args()
	// Print concurepo version
	if *ver {
		
		version.Version()
		os.Exit(0)
	}
	
	// Redirect output to a specifed output file
	if *ofile != "" {
		
		f, err := os.Create(*ofile)
		orig_stdout := os.Stdout

		if err != nil {
			log.Fatalf("Unable to create output file")
		}

		// Defer the closing of the file
		defer func() {
			os.Stdout = orig_stdout
			f.Close()
		}()

		os.Stdout = f
	}

	// Filter flags to apply to scraped HTML
	flag_s := strings.Split(*flags, ",")
	company_s := strings.Split(*company, ",")

	flag_set := make(map[string]struct{})
	company_set := make(map[string]struct{})
	if flag_s[0] != "" {
		for _,v := range(flag_s) {
			flag_set[strings.ToLower(v)] = struct{}{}
		}
	}

	if company_s[0] != "" {
		for _,v := range(company_s) { 
			company_set[strings.ToLower(v)] = struct{}{}
		}
	}

	tab := table.NewWriter()
	tab.SetOutputMirror(os.Stdout)
	tab.SetColumnConfigs([]table.ColumnConfig{
        {
            Name:     "Link", 
            WidthMax: 80,            
        },
    })
	tab.AppendHeader(table.Row{"Company", "Role", "Link","Date"})

	// Websites to scrape jobs from
	websites := []string{"https://github.com/SimplifyJobs/Summer2026-Internships", "https://github.com/SimplifyJobs/New-Grad-Positions"}
	// Waitgroup to wait for all scraping goroutines
	var wg sync.WaitGroup

	if *intern {
		wg.Add(1)
		go scrape.Scraper(websites[0], &wg, *limit, flag_set, company_set, tab)
	} else if *fulltime {
		wg.Add(1)
		go scrape.Scraper(websites[1], &wg, *limit, flag_set, company_set, tab)
	} else {
		wg.Add(1)
		go scrape.Scraper(websites[0], &wg, *limit, flag_set, company_set, tab)
		wg.Add(1)
		go scrape.Scraper(websites[1], &wg, *limit, flag_set, company_set, tab)
	}

	wg.Wait()

	tab.Render()

}