package main

import (
	"fmt"
	"net/http"
	"log"
	"sync"
	"context"
	"regexp"
	"os"
	"concurjob/parse_args"
	"concurjob/version"
	"github.com/PuerkitoBio/goquery"
	"github.com/google/go-github/v76/github"
	"strings"
)

func scrape(website string, wg *sync.WaitGroup, flags []string) {
	defer wg.Done()
	// Http Get request to get the data from webpage and err code 
	data, err := http.Get(website)
	re_simplify := regexp.MustCompile(`https://simplify\.jobs`)
	re_http := regexp.MustCompile(`^https`)
	if err != nil {
		log.Fatal(err)
	}

	// If the status code is not HTTP 200..
	if data.StatusCode != http.StatusOK {
		log.Fatalf("Received a non-200 HTTP GET code: %d", data)
	}

	doc, err := goquery.NewDocumentFromReader(data.Body)

	if err != nil {
		log.Fatal(err)
	}

	data.Body.Close()

	doc.Find("tr").Each(func(i int, tr *goquery.Selection) {
		var rowdata string
		tr.Find("td").Each(func(j int, td *goquery.Selection) {
			// j == 0 -> company name
			// j == 1 -> role name
			// j == 2 -> Location
			// j == 3 -> Application button ( Look for a tags)
			// j == 4 -> date posted
			if j == 1 {
				rowdata = rowdata + strings.TrimSpace(td.Text()) + " | " 
			}
			if j == 3 {
				td.Find("a").Each(func(k int, a *goquery.Selection) {
					link, exists := a.Attr("href")
					if exists && !re_simplify.MatchString(link) && re_http.MatchString(link) {
						rowdata = rowdata + "Link: " + link + " | "
					} 
				})
			}
		})
		if strings.Count(rowdata,"|") == 2  {
			fmt.Printf("%s\n\n", rowdata)
			fmt.Printf("\n------------------------------------------------------------------------------------------------------------------\n")
		}
	})

}

func github_go_api(flags []string, wg *sync.WaitGroup) {
	defer wg.Done()
	var q string
	for i := 0; i<len(flags)-1; i++ {
		q += flags[i] + " OR"
	}
	q = q + " " + flags[len(flags)-1]

	ctx := context.Background()
	client := github.NewClient(nil)

	opt := &github.SearchOptions{
		Sort:        "stars",
		Order:       "desc",
		ListOptions: github.ListOptions{PerPage: 10},
	}

	query := q + " in:name,description"
	results, _, err := client.Search.Repositories(ctx, query, opt)

	if err != nil {
		log.Fatalf("Ran into error while using github-go search")
	}
      
	for _, repo := range results.Repositories {
		fmt.Printf("[ %-32s ] Link=[ %s ]\n", repo.GetFullName(), repo.GetHTMLURL()) 
	}

}

func main() {
	ver, intern, fulltime, ofile, spawn, flags  := parse_args.Parse_args()
	_ = spawn 
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
	// Websites to scrape repositories from
	websites := []string{"https://github.com/SimplifyJobs/Summer2026-Internships", "https://github.com/SimplifyJobs/New-Grad-Positions"}
	// Waitgroup to wait for all scraping goroutines
	var wg sync.WaitGroup

	if *intern {
		wg.Add(1)
		go scrape(websites[0], &wg, flag_s)
	} else if *fulltime {
		wg.Add(1)
		go scrape(websites[1], &wg, flag_s)
	} else {
		wg.Add(1)
		go scrape(websites[0], &wg, flag_s)
		wg.Add(1)
		go scrape(websites[1], &wg, flag_s)
	}

	wg.Wait()

}