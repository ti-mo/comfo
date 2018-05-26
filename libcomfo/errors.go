package libcomfo

import "errors"

var (
	errPktLen          = errors.New("unexpected packet length")
	errPayloadSize     = errors.New("detected payload size larger than packet allows")
	errVerifyChecksum  = errors.New("error verifying checksum of the payload")
	errTooShort        = errors.New("input was too short")
	errTooLong         = errors.New("input was too long")
	errTimeout         = errors.New("operation timed out")
	errScanInput       = errors.New("unexpected input while scanning for token")
	errAck             = errors.New("error writing ACK to connection")
	errInvalidResponse = errors.New("unexpected response type")
	errTooHigh         = errors.New("value too high")
	errTooLow          = errors.New("value is too low")
	errRequestType     = errors.New("unknown request type")
	errResponseType    = errors.New("unknown response type")
	errDecodeNil       = errors.New("attempting to decode nil Packet")
	errZeroValue       = errors.New("unexpected zero value (risk of dividing by zero)")
	errNotExist        = errors.New("item does not exist")
	errUnknownLevel    = errors.New("unknown speed level")
)
