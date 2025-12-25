package scrape

import (
	"fmt"
	"net/http"
	"log"
	"sync"
	"context"
	"regexp"
	"github.com/PuerkitoBio/goquery"
	"github.com/google/go-github/v76/github"
	"github.com/jedib0t/go-pretty/v6/table"
	"strings"
)

func Hyperlink(url string) string {
	osc := "\033]8;;"
	st := "\033\\"
	return fmt.Sprintf("%s%s%s%s%s%s", osc, url, st, "link", osc, st)
}

func Scraper(website string, wg *sync.WaitGroup, limit uint, flag_set, company_set map[string]struct{}, tab table.Writer) {
	/*
	Scrape all the HTML data from a website, filter it, and print the desired output.
	website: target website
	wg: waitgroup
	limit: number of positions to print
	flags: desired flags to filter positions by
	*/
	
	defer wg.Done()
	// Http Get request to get the data from webpage and err code 
	data, err := http.Get(website)

	re_simplify := regexp.MustCompile(`https://simplify\.jobs`)
	re_http := regexp.MustCompile(`^https`)
	re_clean_company := regexp.MustCompile(`[^\x00-\x7F]+`)

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

	var company_flag bool = len(company_set) != 0
	var role_flag bool = len(flag_set) != 0

	var rows uint = 1
	doc.Find("tr").Each(func(i int, tr *goquery.Selection) {
		var rowdata []string
		tr.Find("td").Each(func(j int, td *goquery.Selection) {
			// j == 0 -> company name
			// j == 1 -> role name
			// j == 2 -> Location
			// j == 3 -> Application button ( Look for a tags)
			// j == 4 -> date posted
			if j == 0 {
				company_text := strings.ToLower(strings.TrimSpace(re_clean_company.ReplaceAllString(td.Text(), "")))
				if _,v := company_set[company_text]; (company_flag && v) || !company_flag {
					rowdata = append(rowdata, company_text)
				} 
			}
			if j == 1 {
				role_text := strings.ToLower(strings.TrimSpace(re_clean_company.ReplaceAllString(td.Text(), "")))
				if _,v := flag_set[role_text]; (role_flag && v) || !role_flag {
					rowdata = append(rowdata, role_text)
				}
			} else if j == 3 {
				td.Find("a").Each(func(k int, a *goquery.Selection) {
					link, exists := a.Attr("href")
					if exists && !re_simplify.MatchString(link) && re_http.MatchString(link) {
						rowdata = append(rowdata, Hyperlink(link))
					} 
				})
			} else if j == 4 {
				scraped_date := strings.ToLower(strings.TrimSpace(td.Text()))
				rowdata = append(rowdata, scraped_date)
			}
		})
		if (len(rowdata) == 4) && (rows <= limit) {
			tab.AppendRow(slice_to_row(rowdata))
			rows++
		} else {
			return
		}
	})

}

func slice_to_row(strlist []string) table.Row {
	row := make(table.Row, len(strlist))
	for i,v := range(strlist) {
		row[i] = v 
	}
	return row
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