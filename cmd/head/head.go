package main

import (
	"flag"
	"github.com/apache/yunikorn-core/pkg/head"
)

var endpoint = flag.String("endpoint", "tcp://localhost:3333", "head endpoint")

func main() {
	flag.Parse()
	service := head.NewFleetHttpService(*endpoint)
	service.Start()
	service.Wait()
}
