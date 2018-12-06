package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"sync"
	"time"
)

var totalResponseTime int64
var totalResponses int64
var expectedResponses int64
var totalErrors int64
var resultsCache []*SingleRequestResult
var resultsCacheMutex *sync.Mutex
var startTime int64

func main() {
	resultsCache = []*SingleRequestResult{}
	resultsCacheMutex = &sync.Mutex{}

	b, err := ioutil.ReadFile("./stress.json")
	if err != nil {
		panic(err)
	}

	var conf Config

	err = json.Unmarshal(b, &conf)
	if err != nil {
		panic(err)
	}

	results := make(chan *SingleRequestResult)

	go printResults()

	TMessenger(&conf, results)
}

func readResults(results chan *SingleRequestResult, expectedResultsCount int, done chan bool) {
	expectedResponses = int64(expectedResultsCount)

	for i := 0; i < expectedResultsCount; i++ {
		result := <-results

		if result.ErrorMessage != nil {
			totalErrors++
		}

		totalResponses = int64(i)

		resultsCacheMutex.Lock()
		resultsCache = append(resultsCache, result)
		totalResponseTime += int64(result.TimeTaken)
		resultsCacheMutex.Unlock()
	}

	done <- true
}

func printResults() {
	for {
		resultsCacheMutex.Lock()

		for _, r := range resultsCache {
			if !r.WasSuccessful {
				fmt.Printf("\rERR %v\n", r.ErrorMessage)
			}
		}

		resultsCache = []*SingleRequestResult{}

		resultsCacheMutex.Unlock()

		fmt.Printf("\rRequest %d/%d - %.2f ms - %d - E %d",
			totalResponses+1, expectedResponses,
			(float64(totalResponseTime)/float64(totalResponses+1))/float64(time.Millisecond),
			time.Now().Unix()-startTime, totalErrors)

		time.Sleep(50 * time.Millisecond)

		if totalResponses == expectedResponses {
			break
		}
	}
}
