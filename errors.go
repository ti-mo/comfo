package main

import "errors"

var (
	errDelim          = errors.New("could not detect start and end delimiters of packet")
	errPktLen         = errors.New("unexpected packet length")
	errPayloadSize    = errors.New("detected payload size larger than packet allows")
	errVerifyChecksum = errors.New("error verifying checksum of the payload")
	errChecksum       = errors.New("error calculating checksum")
	errTooShort       = errors.New("input was too short")
)
