package main

import (
	"flag"

	"github.com/fbonareis/goexpert-stress-test/pkg/stresstest"
)

var (
	url         string
	requests    int
	concurrency int
)

func main() {
	flag.StringVar(&url, "url", "", "URL of the service to be tested")
	flag.IntVar(&requests, "requests", 0, "Total number of requests")
	flag.IntVar(&concurrency, "concurrency", 0, "Number of simultaneous calls")
	flag.Parse()

	st := stresstest.New(url, requests, concurrency)
	st.Start()
}
