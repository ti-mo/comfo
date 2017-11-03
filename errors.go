package main

import "errors"

var (
	errSetPollFailed = errors.New("failure polling result after set operation")
	errOutOfRange    = errors.New("value out of range")
)
