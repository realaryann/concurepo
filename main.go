package main

import (
	"fmt"
	"net/http"
	"log"
	"time"
	"io/ioutil"
	"concurjob/parse_args"
	"concurjob/version"
)

func scrape() {
	url := "https://realaryann.github.io/"
	
	// Http Get request to get the data from webpage and err code 
	data, err := http.Get(url)

	if err != nil {
		log.Fatal(err)
	}

	// If the status code is not HTTP 200..
	if data.StatusCode != http.StatusOK {
		log.Fatalf("Received a non-200 HTTP GET code: %d", data)
	}

	// Must close the data's body after computing
	defer data.Body.Close()

	// Actually read the data from the response
	body, err := ioutil.ReadAll(data.Body)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(body))
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

	for i := uint(0); i<*spawn; i++ {
		go scrape()
		time.Sleep(100 * time.Millisecond)
	}

	fmt.Println(*ofile, *spawn)

}