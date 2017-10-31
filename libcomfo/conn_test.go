package libcomfo

import (
	"bytes"
	"fmt"
	"testing"
	"time"
)

// TestConn is a structure that implements io.ReadWriter
// and can be used to mock a connection with the unit.
type TestConn struct {
	Receives []byte
	Emits    []byte

	received int
}

func (tr *TestConn) Read(p []byte) (n int, err error) {

	// Return the data in the connection's Emits field
	n = copy(p, tr.Emits)

	return
}

func (tr *TestConn) Write(p []byte) (n int, err error) {

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

	// Return the length of the 'written' slice
	return len(p), err
}

func TestWaitTimeout(t *testing.T) {

	returnChan := make(chan bool)

	// Return timer at 2 milliseconds
	returnTimer := func() {
		<-time.NewTimer(time.Millisecond * 2).C
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
