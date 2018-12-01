package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/calebhiebert/gobbl/messenger"
	"github.com/matoous/go-nanoid"
)

var client *http.Client
var startTime int64

// TMessenger starts the test for a facebook messenger bot
func TMessenger(config *Config) {
	spacing := time.Duration(int(config.TestDuration) / config.Requests)
	resultChan := make(chan *SingleRequestResult)
	done := make(chan bool)
	client = &http.Client{
		Timeout: 30 * time.Second,
	}

	go readResults(resultChan, config.Requests, done)

	fmt.Printf("Starting test with %.1f/sec\n", float64(config.Requests)/float64(config.TestDuration*time.Second))
	startTime = time.Now().Unix()

	for i := 0; i < config.Requests; i++ {
		time.Sleep(spacing)
		go doRequest(config, resultChan)
	}

	<-done
}

// readResults will wait for the specified number of results and compile results
func readResults(resultsChan chan *SingleRequestResult, expectedResults int, done chan bool) {
	totalResponseTime := 0

	for i := 0; i < expectedResults; i++ {
		result := <-resultsChan

		totalResponseTime += int(result.TimeTaken)

		if !result.WasSuccessful {
			fmt.Printf("\rERR %d - %v\n", i, result.ErrorMessage)
		}

		fmt.Printf("\rRequest %d/%d - %.2f ms - %d",
			i+1, expectedResults,
			(float64(totalResponseTime)/float64(i))/float64(time.Millisecond),
			time.Now().Unix()-startTime)
	}

	done <- true
}

// doRequest will make a single messenger request
func doRequest(config *Config, res chan *SingleRequestResult) {

	req := generateMessengerRequest()
	jsonBytes, err := json.Marshal(&req)
	if err != nil {
		panic(err)
	}

	start := time.Now().UnixNano()

	_, err = client.Post(config.Messenger.Endpoint, "application/json", bytes.NewReader(jsonBytes))

	duration := time.Duration(time.Now().UnixNano() - start)

	res <- &SingleRequestResult{
		TimeTaken:     duration,
		WasSuccessful: err == nil,
		ErrorMessage:  err,
	}
}

func generateMessengerRequest() fb.WebhookRequest {
	mid, _ := gonanoid.Nanoid(64)

	return fb.WebhookRequest{
		Object: "page",
		Entry: []fb.WebhookEntry{
			fb.WebhookEntry{
				ID:   "123456789",
				Time: time.Now().UnixNano() * int64(time.Millisecond),
				Messaging: []fb.MessagingItem{
					fb.MessagingItem{
						Sender: fb.User{
							ID: "123456789",
						},

						Recipient: fb.User{
							ID: "123456789",
						},

						Timestamp: time.Now().UnixNano() * int64(time.Millisecond),

						Message: fb.WHMessage{
							MID:  mid,
							Seq:  1,
							Text: "example text",
						},
					},
				},
			},
		},
	}
}
