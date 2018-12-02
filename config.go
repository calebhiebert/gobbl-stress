package main

import (
	"time"
)

// Config represents a gobbl-stress config. This config tells the program
// how to stress test the bot
type Config struct {
	TestDuration time.Duration `json:"testDuration"`
	Requests     int           `json:"requests"`

	Messenger *MessengerConfig `json:"messenger"`
}

// SingleRequestResult contains the results of a single request
type SingleRequestResult struct {
	TimeTaken     time.Duration
	WasSuccessful bool
	ErrorMessage  interface{}
}

// MessengerConfig holds test configuration specific to facebook messenger
type MessengerConfig struct {
	PSIDList []string `json:"psids"`
	Messages []string `json:"messages"`
	Endpoint string   `json:"endpoint"`
}
