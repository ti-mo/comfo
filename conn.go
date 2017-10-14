package main

import (
	"bufio"
	"io"
	"sync"
	"time"
)

var (
	// Serial connection Mutex
	pm = &sync.Mutex{}
)

func WaitTimeout(wait <-chan bool, timeout *time.Timer) (success bool, err error) {
	select {
	case result := <-wait:
		// Operation succeeded
		return result, nil
	case <-timeout.C:
		// Time out the transaction
		return false, errTimeout
	}
}

// ScanTimeout advances the given scanner until the timer expires.
// All output found on the way is returned.
func ScanTimeout(scanner *bufio.Scanner, timer *time.Timer) (out []byte, err error) {

	done := make(chan bool)

	// Advance the scanner past one delimiter
	go func() { done <- scanner.Scan() }()

	// Wait for delimiter to appear, or time out
	if _, err := WaitTimeout(done, timer); err != nil {
		return out, err
	}

	// Return all characters found on the way to the delimiter
	out = scanner.Bytes()

	return
}

func ReadPacket(conn io.Reader, timer *time.Timer, expect bool) (out Packet, err error) {

	// Read the first byte in the stream
	scanner := bufio.NewScanner(conn)

	// Split function with configurable
	splitFunc := func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		for i := 0; i < len(data); i++ {
			// Look for escape characters
			if data[i] == esc {
				// Look for the given byte following the escape character within bounds
				if i+1 < len(data) {
					char := data[i+1]
					switch char {
					case ack, start, end:
						// Advance the pointer past the 2nd command character
						return i + 2, data[:i], nil
					case esc:
						// Ignore double 0x07
						return i + 2, nil, nil
					default:
					}
				}
			}
		}
		return
	}

	scanner.Split(splitFunc)

	// Read until first escape sequence (ACK)
	pr, err := ScanTimeout(scanner, timer)
	if err != nil {
		return out, err
	}

	if len(pr) != 0 {
		return out, errScanInput
	}

	// Return if we're not expecting data (only ACK)
	if !expect {
		return
	}
	// Read second escape sequence (start)
	pr, err = ScanTimeout(scanner, timer)
	if err != nil {
		return out, err
	}

	if len(pr) != 0 {
		return out, errScanInput
	}

	// Read second escape sequence (data + stop)
	pr, err = ScanTimeout(scanner, timer)
	if err != nil {
		return out, err
	}

	// Abort on invalid input
	if len(pr) > 0 && len(pr) < 4 {
		return out, errTooShort
	}

	// Decode the payload
	if len(pr) >= 4 {
		out, err = UnmarshalPacket(pr)
		if err != nil {
			return out, err
		}
	}

	return
}

// WritePacket serializes a Packet into wire representation, wrapped
// in start and end sequences, and writes to the given connection.
func WritePacket(in Packet, conn io.Writer) (out bool, err error) {

	// Get wire representation of the given Packet
	pb, err := MarshalPacket(in)
	if err != nil {
		return out, errMarshalPacket
	}

	// Wrap start and end escape sequences
	wr := append(pktStart, pb...)
	wr = append(wr, pktEnd...)

	// Write slice to connection
	_, err = conn.Write(wr)
	if err != nil {
		return false, errWrite
	}

	return true, nil
}

// WriteAck writes an ACK response to a connection.
func WriteAck(conn io.Writer) (out bool, err error) {
	num, err := conn.Write(pktAck)
	if err != nil {
		return false, err
	}

	if num != len(pktAck) {
		return false, errWrite
	}

	return
}

func Query(in Packet, conn io.ReadWriter) (out Packet, err error) {

	// Take out lock - start critical section
	pm.Lock()
	defer pm.Unlock()

	// Write the Packet to the connection
	if _, err := WritePacket(in, conn); err != nil {
		return out, err
	}

	// Start a query-scoped timeout
	timer := time.NewTimer(time.Second)

	// Read packet from the remote
	out, err = ReadPacket(conn, timer, in.Expect)
	if err != nil {
		return out, err
	}

	// Send ACK
	WriteAck(conn)

	// Return
	return out, nil
}