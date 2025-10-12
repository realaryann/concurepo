package main

import (
	"fmt"
	"net/http"
	"log"
	"sync"
	"io/ioutil"
	"concurjob/parse_args"
	"concurjob/version"
)

func scrape(websites []string, wg *sync.WaitGroup) {
	defer wg.Done()
	// Http Get request to get the data from webpage and err code 
	for _,v := range(websites) {
		data, err := http.Get(v)

		if err != nil {
			log.Fatal(err)
		}

		// If the status code is not HTTP 200..
		if data.StatusCode != http.StatusOK {
			log.Fatalf("Received a non-200 HTTP GET code: %d", data)
		}

		// Actually read the data from the response

		body, err := ioutil.ReadAll(data.Body)

		data.Body.Close()

		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(string(body))
	}
}

func main() {
	/*
	CLI Args
	-spawn x: Spawn x goroutines to search for repositories to contribute to
	-o file.txt: Output the links to the repositories 
	-version: Output the version of concurjob
	*/
	ver, ofile, spawn := parse_args.Parse_args()

	if *ver {
		version.Version()
	}

	websites := []string{"https://github.com/trending"}

	var wg sync.WaitGroup

	for i := uint(0); i<*spawn; i++ {
		wg.Add(1)
		go scrape(websites, &wg)
	}

	wg.Wait()

	fmt.Println(*ofile, *spawn)

}