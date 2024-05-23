package stresstest

import (
	"fmt"
	"net/http"
	"sync"
	"time"
)

type RequestResult struct {
	HttpStatus int
	HasError   bool
}

type StressTest struct {
	url           string
	requests      int
	concurrency   int
	results       []RequestResult
	totalDuration time.Duration
	mu            sync.Mutex
}

func New(url string, requests, concurrency int) *StressTest {
	return &StressTest{
		url:         url,
		requests:    requests,
		concurrency: concurrency,
	}
}

func (st *StressTest) Start() {
	start := time.Now()
	ch := make(chan struct{}, st.requests)
	wg := sync.WaitGroup{}
	wg.Add(st.requests)

	for i := 0; i < st.concurrency; i++ {
		go st.worker(ch, &wg)
	}
	for i := 0; i < st.requests; i++ {
		ch <- struct{}{}
	}
	close(ch)

	wg.Wait()
	st.totalDuration = time.Since(start)
	st.printReport()
}

func (st *StressTest) doRequest() {
	var rr RequestResult
	req, _ := http.NewRequest("GET", st.url, nil)
	resp, err := http.DefaultClient.Do(req)

	rr.HasError = err != nil
	rr.HttpStatus = resp.StatusCode

	st.mu.Lock()
	defer st.mu.Unlock()
	st.results = append(st.results, rr)
}

func (st *StressTest) worker(ch chan struct{}, wg *sync.WaitGroup) {
	for range ch {
		st.doRequest()
		wg.Done()
	}
}

func (st *StressTest) printReport() {
	httpStatusMap := make(map[int]int)
	hasErrorMap := make(map[bool]int)

	for _, r := range st.results {
		httpStatusMap[r.HttpStatus] += 1
		hasErrorMap[r.HasError] += 1
	}

	fmt.Println("---------------- STRESS TEST RESULT ----------------")
	fmt.Printf("Total time spent executing: %s\n", st.totalDuration.String())
	fmt.Printf("Total number of requests made: %d\n", len(st.results))
	fmt.Printf("Total number of requests with error: %d\n", hasErrorMap[true])
	for k, v := range httpStatusMap {
		fmt.Printf("Number of requests with HTTP status %d: %d\n", k, v)
	}
	fmt.Println("----------------------------------------------------")
}
