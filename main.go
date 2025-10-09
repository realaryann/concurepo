package main

import (
	//"fmt"
	"flag"
	"concurjob/version"
)

func main() {
	/*
	CLI Args
	-spawn x: Spawn x goroutines to search for repositories to contribute to
	-o file.txt: Output the links to the repositories 
	-version: Output the version of concurjob
	*/
	ver := flag.Bool("version", false, "Current version number")

	flag.Parse()

	if *ver {
		version.Version()
	}

}