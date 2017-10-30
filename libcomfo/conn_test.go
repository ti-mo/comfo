package libcomfo

import (
	"bytes"
	"fmt"
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
