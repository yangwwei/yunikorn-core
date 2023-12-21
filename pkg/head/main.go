package main

import (
	"flag"
	"os"
)

var endpoint = flag.String("endpoint", "tcp://localhost:3333", "head endpoint")

func main() {
	flag.Parse()
	Run(*endpoint)
	os.Exit(0)
}
