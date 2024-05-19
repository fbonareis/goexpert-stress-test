package main

import (
	"flag"
	"fmt"
	"net/http"
	"sync"
	"time"
)

type Result struct {
	HttpStatus int
	Duration   time.Duration
}

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

	results := []Result{}
	chanResults := make(chan Result)
	wgResults := sync.WaitGroup{}
	wgResults.Add(1)
	go func() {
		for result := range chanResults {
			results = append(results, result)
		}
		wgResults.Done()
	}()

	chanRequests := make(chan int)
	wgRequests := sync.WaitGroup{}
	wgRequests.Add(requests)

	for i := 1; i <= concurrency; i++ {
		go worker(&url, chanRequests, chanResults, &wgRequests)
	}
	for i := 0; i < requests; i++ {
		chanRequests <- i
	}

	wgRequests.Wait()
	close(chanResults)
	wgResults.Wait()

	printReport(results)
}

func worker(url *string, requests chan int, results chan Result, wg *sync.WaitGroup) {
	for range requests {
		if r, err := doRequest(url); err == nil {
			results <- *r
		}
		wg.Done()
	}
}

func doRequest(url *string) (*Result, error) {
	start := time.Now()
	req, err := http.NewRequest("GET", *url, nil)
	if err != nil {
		return nil, err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	return &Result{
		HttpStatus: resp.StatusCode,
		Duration:   time.Since(start),
	}, nil
}

func printReport(results []Result) {
	httpStatusMap := make(map[int]int)
	totalTime := time.Time{}

	for _, r := range results {
		httpStatusMap[r.HttpStatus] = httpStatusMap[r.HttpStatus] + 1
		totalTime = totalTime.Add(r.Duration)
	}

	fmt.Println("---------------- STRESS TEST RESULT ----------------")
	fmt.Printf("Total time spent executing: %s\n", time.Duration(totalTime.Nanosecond()).String())
	fmt.Printf("Total number of requests made: %d\n", len(results))
	for k, v := range httpStatusMap {
		fmt.Printf("Number of requests with HTTP status %d: %d\n", k, v)
	}
	fmt.Println("----------------------------------------------------")
}
