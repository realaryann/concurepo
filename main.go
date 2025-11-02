package main

import (
	"fmt"
	"net/http"
	"log"
	"sync"
	"context"
	"regexp"
	"os"
	"concurepo/parse_args"
	"concurepo/version"
	"github.com/PuerkitoBio/goquery"
	"github.com/google/go-github/v76/github"
	"strings"
)

func scrape(website string, wg *sync.WaitGroup, flags []string) {
	defer wg.Done()
	// Http Get request to get the data from webpage and err code 
	data, err := http.Get(website)
	re_repo := regexp.MustCompile(`^\/[^\/]+\/[^\/]+\/[^\/]+$`)
	re_user := regexp.MustCompile(`^\/[^\/]+$`)
	re_gh := regexp.MustCompile(`^\/[^\/]+\/[^\/]+$`)
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

	link_select := doc.Find("a")
	link_select.Each(func(i int, element *goquery.Selection) {

		text := element.Text()
		text = strings.TrimSpace(text)
		for _,v := range(flags) {
			if strings.Contains(text, v) {
				link, exists := element.Attr("href")
				if exists {
					if re_repo.MatchString(link) || re_user.MatchString(link) || re_gh.MatchString(link) {
						link = "https://github.com"+link
						fmt.Printf("[ %-32s ] Link=[ %s ]\n", text, link)
	
					} else {
						fmt.Printf("[ %-32s ] Link=[ %s ]\n", text, link)
					}
				}
				break
			}
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
	ver, ofile, spawn, flags  := parse_args.Parse_args()
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
	websites := []string{"https://github.com/trending"}
	// Waitgroup to wait for all scraping goroutines
	var wg sync.WaitGroup

	wg.Add(1)
	go scrape(websites[0], &wg, flag_s)

	if flag_s[0] != "" {
		wg.Add(1)
		go github_go_api(flag_s, &wg)
	}

	wg.Wait()

}