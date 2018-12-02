package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/calebhiebert/gobbl/messenger"
	"github.com/matoous/go-nanoid"
)

var client *http.Client

// TMessenger starts the test for a facebook messenger bot
func TMessenger(config *Config, results chan *SingleRequestResult) {
	spacing := time.Duration(int(config.TestDuration) / config.Requests)
	done := make(chan bool)
	client = &http.Client{
		Timeout: 30 * time.Second,
	}

	go readResults(results, config.Requests, done)

	fmt.Printf("Starting test with %d requests per second\n", config.Requests/int(config.TestDuration.Seconds()))
	startTime = time.Now().Unix()

	for i := 0; i < config.Requests; i++ {
		time.Sleep(spacing)
		go doRequest(config, results)
	}

	<-done
}

// doRequest will make a single messenger request
func doRequest(config *Config, res chan *SingleRequestResult) {

	req := generateMessengerRequest(config)
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

func generateMessengerRequest(config *Config) fb.WebhookRequest {
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
							ID: config.Messenger.PSIDList[rand.Intn(len(config.Messenger.PSIDList))],
						},

						Recipient: fb.User{
							ID: "123456789",
						},

						Timestamp: time.Now().UnixNano() * int64(time.Millisecond),

						Message: fb.WHMessage{
							MID:  mid,
							Seq:  1,
							Text: config.Messenger.Messages[rand.Intn(len(config.Messenger.Messages))],
						},
					},
				},
			},
		},
	}
}
