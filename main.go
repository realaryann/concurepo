package main

import (
	"fmt"
	"concurjob/parse_args"
	"concurjob/version"
)

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
	fmt.Println(*ofile, *spawn)

}