package libcomfo

import (
	"bytes"
	"fmt"
	"io"
	"reflect"
	"testing"
	"time"
)

// TestConn is a structure that implements io.ReadWriter and can be used
// to mock a connection with the unit. HangStart makes initial Read and Write
// calls hang. HangEOF makes Read hang on the second call when a Scanner
// would expect an EOF error to be returned.
// Not thread-safe.
type TestConn struct {
	Receives []byte
	Emits    []byte

	HangTime  time.Duration
	HangStart bool
	HangEOF   bool

	Limit int

	emitted  int
	received int
}

func (tr *TestConn) Read(p []byte) (n int, err error) {

	// Simulate read delay before returning output or EOF
	if tr.HangTime > 0 && tr.HangStart {
		time.Sleep(tr.HangTime)
	}

	// Return EOF when buffer is fully copied
	if tr.emitted == len(tr.Emits) {

		if tr.HangTime > 0 && tr.HangEOF {
			time.Sleep(tr.HangTime)
		}

		return 0, io.EOF
	}

	// Return the data in the connection's Emits field
	n = copy(p, tr.Emits)

	// Keep track of how many bytes we've been able to copy
	tr.emitted = tr.emitted + n

	// Return mocked read bytes total
	if tr.Limit > 0 {
		n = tr.Limit
	}

	return
}

func (tr *TestConn) Write(p []byte) (n int, err error) {

	// Simulate write delay
	if tr.HangTime > 0 {
		time.Sleep(tr.HangTime)
	}

	// Bounds check for Receives
	if tr.received+len(p) > len(tr.Receives) {
		return n,
			fmt.Errorf("write to TestConn exceeded expected input specified in Receives: %v",
				// This formula determines the offset of the slice
				// of surplus characters over the struct's 'Receives' field.
				p[len(p)-(tr.received+len(p)-len(tr.Receives)):])
	}

	// Expect the written data to correspond to what's in the 'Receives' field.
	// Offset the comparison against the write index (received).
	if want, got := tr.Receives[tr.received:tr.received+len(p)], p; !bytes.Equal(want, got) {
		err = fmt.Errorf(
			"unexpected write to TestConn:\n  - want: %v\n  -  got: %v", want, got)
	}

	// Keep a record of the bytes we successfully received
	// and compared against the 'Receives' buffer.
	tr.received = tr.received + len(p)

	// Return mocked written bytes total
	if tr.Limit > 0 {
		n = tr.Limit
	} else {
		// Return the length of the 'written' slice
		n = len(p)
	}

	return
}

func TestWaitTimeout(t *testing.T) {

	returnChan := make(chan bool)

	// Return timer at 2 milliseconds
	returnTimer := func() {
		time.Sleep(time.Millisecond * 2)
		returnChan <- true
	}

	// Run test 10 times back-to-back
	for i := 0; i < 10; i++ {
		t.Run(t.Name(), func(t *testing.T) {

			// Start a timeout timer with a timeout higher than the return timer
			timeOutTimer := time.NewTimer(time.Millisecond * 3)
			go returnTimer() // Start return timer

			// Expect returnChan to unblock before timeOutTimer
			_, err := WaitTimeout(returnChan, timeOutTimer)
			if err != nil {
				t.Fatal("returnChan did not unblock before timeOutTimer")
			}

			// Start a timeout timer with a lower timeout than return timer
			timeOutTimer = time.NewTimer(time.Millisecond * 1)
			go returnTimer() // Start return timer

			// Expect returnChan to unblock before timeOutTimer
			_, err = WaitTimeout(returnChan, timeOutTimer)
			if err != errTimeout {
				t.Fatal("timeOutTimer did not expire before returnChan unblock")
			}

			// Wait for the return timer to send on channel
			// This value needs to be read to prevent it from interfering other tests
			<-returnChan
		})
	}
}

func TestReadPacket(t *testing.T) {

	rt := []struct {
		name string
		tc   TestConn
		pkt  Packet
		err  error
	}{
		{
			name: "single and double escaped 0x07s in payload",
			tc: TestConn{
				Emits: []byte{ // Response
					esc, ack,
					0x07, 0xF0, // Frame Start
					0x00, uint8(getFans + 1), // response type
					0x06,       // Length
					0xAA, 0xBB, // in/out percents
					0x11, 0x07, 0x07, // in speed (one seven)
					0x07, 0x07, 0x07, 0x07, // out speed (two sevens)
					0x4A, // Checksum
					esc, end,
					esc, ack,
				},
			},
			pkt: Packet{
				Command: uint8(getFans + 1),
				Data:    []byte{0xAA, 0xBB, 0x11, 0x7, 0x7, 0x7},
			},
		},
		{
			name: "timeout reading first response ACK",
			tc: TestConn{
				HangStart: true,
				HangTime:  2 * time.Millisecond,
			},
			err: errTimeout,
		},
		{
			name: "timeout reading start sequence",
			tc: TestConn{
				Emits:    []byte{esc, ack},
				HangEOF:  true,
				HangTime: 2 * time.Millisecond,
			},
			err: errTimeout,
		},
		{
			name: "timeout reading end sequence",
			tc: TestConn{
				Emits:    []byte{esc, ack, esc, start},
				HangEOF:  true,
				HangTime: 2 * time.Millisecond,
			},
			err: errTimeout,
		},
		{
			name: "garbage before ACK",
			tc: TestConn{
				Emits: []byte{esc, ack, 0xF0, 0x0B, esc, start},
			},
			err: errScanInput,
		},
		{
			name: "impossible packet length (need at least 4 bytes)",
			tc: TestConn{
				Emits: []byte{esc, ack, esc, start, 0xF0, 0x0B, esc, end},
			},
			err: errTooShort,
		},
		{
			name: "packet unmarshal error (any)",
			tc: TestConn{
				Emits: []byte{esc, ack, esc, start, 0xF0, 0x0B, 0xBA, 0xAA, esc, end},
			},
			err: errPayloadSize,
		},
	}

	for _, tt := range rt {
		t.Run(tt.name, func(t *testing.T) {

			// Copy the TestConn structure of this test cycle into
			// the goroutine to prevent the test engine from modifying
			// the structure before it is read by the library.
			tc := tt.tc

			// Start timeout timer
			to := time.NewTimer(time.Millisecond * 1)

			// Attempt to read packet from mock connection
			pkt, err := ReadPacket(&tc, to, true)

			if want, got := tt.err, err; want != got {
				t.Fatalf("unexpected error reading packet:\n- want: %v\n-  got: %v",
					want, got)
			}

			if want, got := tt.pkt, pkt; !reflect.DeepEqual(want, got) {
				t.Fatalf("unexpected packet read from connection:\n- want: %v\n-  got: %v",
					want, got)
			}
		})
	}
}
