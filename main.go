package main

import (
	"fmt"
	"net/http"
	"log"
	"sync"
	"io/ioutil"
	"concurjob/parse_args"
	"concurjob/version"
	"strings"
)

func scrape(website string, id uint, wg *sync.WaitGroup) {
	defer wg.Done()
	// Http Get request to get the data from webpage and err code 
	data, err := http.Get(website)

	if err != nil {
		log.Fatal(err)
	}

	// If the status code is not HTTP 200..
	if data.StatusCode != http.StatusOK {
		log.Fatalf("Received a non-200 HTTP GET code: %d", data)
	}

	// Actually read the data from the response

	body, err := ioutil.ReadAll(data.Body)
	
	_ = body

	data.Body.Close()

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("id = ", id, "\n")
}

func main() {
	/*
	CLI Args
	-spawn x: Spawn x goroutines to search for repositories to contribute to
	-o file.txt: Output the links to the repositories 
	-version: Output the version of concurjob
	*/
	ver, ofile, spawn, flags := parse_args.Parse_args()

	if *ver {
		version.Version()
	}

	flag_s := strings.Split(*flags, " ")

	websites := []string{"https://github.com/trending"}

	var wg sync.WaitGroup

	for i := uint(0); i<*spawn; i++ {
		wg.Add(1)
		go scrape(websites[0], i, &wg)
	}

	wg.Wait()
	fmt.Println(flag_s)
	fmt.Println(*ofile, *spawn)

}