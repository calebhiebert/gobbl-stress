package main

import (
	"time"
)

func main() {
	conf := Config{
		TestDuration: time.Minute,
		Requests:     3800,
		Messenger: &MessengerConfig{
			Endpoint: "https://google.ca",
		},
	}

	TMessenger(&conf)
}
