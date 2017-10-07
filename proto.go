package main

var (
	ack      = []byte{0x07, 0xF3}
	pktStart = []byte{0x07, 0xF0}
	pktEnd   = []byte{0x07, 0x0F}

	// byte 1-2 - start
	pktStartLen = len(pktStart)
	// byte -1, -2 - end
	pktEndLen = len(pktEnd)

	// byte 2-3 - command
	cmdOffset = pktStartLen
	// byte 4 - payload length
	lenOffset = cmdOffset + 2
	// byte 5-x - payload
	dataOffset = lenOffset + 1
	// byte -3 (from end!) - checksum
	cksumOffset = 3
)

type Packet struct {
	Command uint8
	Data    []byte
}

// MarshalPacket takes a Packet structure and converts it into wire protocol.
func MarshalPacket(in Packet) (out []byte, err error) {

	// TODO: detect duplicate 0x7's before reallocating target slice
	dataLen := len(in.Data)
	out = make([]byte, dataLen+8)
	outLen := len(out)

	copy(out[:pktStartLen], pktStart)                        // Preamble
	copy(out[cmdOffset:lenOffset], []byte{0x00, in.Command}) // Command
	out[lenOffset] = uint8(dataLen)                          // Payload length
	copy(out[dataOffset:dataOffset+dataLen], in.Data)        // Payload

	// Calculate checksum
	cksum, err := calculateChecksum(out[cmdOffset : outLen-cksumOffset])
	if err != nil {
		return out, err
	}

	out[outLen-cksumOffset] = cksum
	copy(out[outLen-pktEndLen:], pktEnd) // End

	return out, nil
}

// UnmarshalPacket parses a byte slice of wire protocol into a Packet structure.
// The byte slice must start with pktStart and end with pktEnd,
// or the operation will fail.
func UnmarshalPacket(in []byte) (out Packet, err error) {
	if len(in) < 8 {
		return out, errPktLen
	}

	// Offsets of this packet's start and end delimiters
	var endOffset int = len(in) - pktEndLen

	// Detect packet start and end
	if !byteCmp(in[:cmdOffset], pktStart) || !byteCmp(in[endOffset:], pktEnd) {
		return out, errDelim
	}

	// The size of the payload
	// TODO: Bump this by 1 for every double 0x07 BEFORE copying payload
	var dataLen uint8 = in[lenOffset]

	// Compare the expected size of the payload to the size of the packet.
	// There are 8 non-payload bytes in a packet, so the size cannot exceed this.
	if int(dataLen) > len(in)-8 {
		return out, errPayloadSize
	}

	// Verify checksum, extract command and payload
	if cksum, err := verifyChecksum(in); cksum && err == nil {
		out.Command = in[cmdOffset+1] // Actual command is 2nd byte of command field
		out.Data = in[dataOffset : dataOffset+int(dataLen)]
	} else if !cksum {
		return out, errVerifyChecksum
	} else if err != nil {
		return out, errChecksum
	}

	return out, nil
}

// calculateChecksum calculates the checksum of a Packet byte string
// excluding start and end.
// TODO: Double 0x7
func calculateChecksum(in []byte) (uint8, error) {

	// Empty packet without start, end and checksum has at least 3 bytes
	if len(in) < 3 {
		return 0, errTooShort
	}

	// Allocate large enough int to hold our calculation
	var tempSum uint64

	// Add all bytes together
	for _, v := range in {
		tempSum += uint64(v)
	}

	// Add magic value
	tempSum += 173

	// Truncate the sum to a single byte
	return uint8(tempSum), nil
}

// verifyChecksum computes a checksum over the packet excluding start and end.
// TODO: Implement this
func verifyChecksum(in []byte) (bool, error) {
	return true, nil
}

// byteCmp compares the value and length of two byte slices.
func byteCmp(a []byte, b []byte) bool {
	if len(a) != len(b) {
		return false
	}

	for i, v := range a {
		if v != b[i] {
			return false
		}
	}

	return true
}
