package comfoserver

import "errors"

var (
	errSetPollFailed = errors.New("failure polling result after set operation")
	errBothAbsRel    = errors.New("Abs and Rel are mutually exclusive")
	errNoneAbsRel    = errors.New("one of Abs and Rel is required")
)
