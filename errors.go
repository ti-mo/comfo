package main

import "errors"

var (
	errDelim           = errors.New("could not detect start and end delimiters of packet")
	errPktLen          = errors.New("unexpected packet length")
	errPayloadSize     = errors.New("detected payload size larger than packet allows")
	errVerifyChecksum  = errors.New("error verifying checksum of the payload")
	errChecksum        = errors.New("error calculating checksum")
	errTooShort        = errors.New("input was too short")
	errTooLong         = errors.New("input was too long")
	errMarshalPacket   = errors.New("encountered error marshaling packet")
	errTimeout         = errors.New("operation timed out")
	errReadByte        = errors.New("unable to read byte from connection")
	errScanInput       = errors.New("unexpected input while scanning for token")
	errWrite           = errors.New("error writing packet to connection")
	errAck             = errors.New("error writing ACK to connection")
	errInvalidResponse = errors.New("unexpected response type")
	errTooHigh         = errors.New("value too high")
	errTooLow          = errors.New("value is too low")
	errRequestType     = errors.New("unknown request type")
	errResponseType    = errors.New("unknown response type")
)
