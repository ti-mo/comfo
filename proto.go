package main

var (
	esc   byte = 0x07
	ack   byte = 0xF3
	start byte = 0xF0
	end   byte = 0x0F

	pktStart = []byte{esc, start}
	pktEnd   = []byte{esc, end}
	pktAck   = []byte{esc, ack}

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
	// Command type or response code
	Command uint8

	// Payload byte slice
	Data []byte

	// Whether to expect a response
	Expect bool
}

// MarshalPacket takes a Packet structure and converts it into wire protocol.
func MarshalPacket(in Packet) (out []byte, err error) {

	// Save the original, unescaped length of the payload
	dataLen := len(in.Data)

	// Escape the payload and re-calculate (wire) length
	data := escapeData(in.Data)
	dataLenWire := len(data)

	// Make output slice with payload length + 8
	out = make([]byte, dataLenWire+8)
	outLen := dataLenWire + 8

	copy(out[:pktStartLen], pktStart)                        // Preamble
	copy(out[cmdOffset:lenOffset], []byte{0x00, in.Command}) // Command
	out[lenOffset] = uint8(dataLen)                          // Payload length
	copy(out[dataOffset:dataOffset+dataLenWire], data)       // Payload

	// Calculate checksum
	cksum, err := calculateChecksum(in.Command, uint8(dataLen), in.Data)
	if err != nil {
		return out, err
	}

	out[outLen-cksumOffset] = cksum      // Set checksum
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

	// Detect packet start and end
	if !byteCmp(in[:cmdOffset], pktStart) || !byteCmp(in[len(in)-pktEndLen:], pktEnd) {
		return out, errDelim
	}

	// The command
	var cmd uint8 = in[cmdOffset+1]

	// The size of the payload
	var dataLen uint8 = in[lenOffset]

	// Offset of the check
	var endOffset int = len(in) - pktEndLen

	// Get checksum from packet
	var cksum uint8 = in[endOffset-1]

	data := unescapeData(in[dataOffset : endOffset-1])

	// Compare the expected size of the payload to the size of the packet.
	// There are 8 non-payload bytes in a packet, so we can reliably calculate
	// how many bytes the payload *should* be.
	// Escaped 0x07's do not count towards dataLen, so needs to be unescaped first.
	if int(dataLen) != len(data) {
		return out, errPayloadSize
	}

	check, err := verifyChecksum(cksum, cmd, dataLen, data)
	if err != nil {
		return out, errChecksum
	} else if !check {
		return out, errVerifyChecksum
	}

	// Extract command and payload
	out.Command = cmd // Actual command is 2nd byte of command field
	out.Data = data

	return out, nil
}

// escapeData escapes 0x07 characters with 0x07 in a payload.
func escapeData(in []byte) []byte {

	out := make([]byte, 0)

	for _, v := range in {
		// Append an extra 0x07 when a 0x07 is read
		if v == 0x07 {
			out = append(out, 0x07)
		}

		out = append(out, v)
	}

	return out
}

// unescapeData unescapes escape sequences of a wire-format payload.
func unescapeData(in []byte) []byte {

	out := make([]byte, 0)

	for i := 0; i < len(in); i++ {
		// Detect two successive 0x07
		if in[i] == 0x07 && in[i+1] == 0x07 {
			// Only append a single 0x07
			out = append(out, 0x07)

			// Advance the window twice when successfully unescaping
			// to prevent re-evaluating on the second 0x07.
			i++
		} else {
			out = append(out, in[i])
		}
	}

	return out
}

// calculateChecksum calculates the checksum of a Packet command,
// length and payload without (!) escaped 7s.
func calculateChecksum(cmd uint8, dataLen uint8, data []byte) (uint8, error) {

	// Allocate large enough int to hold our calculation
	var tempSum uint64

	tempSum += uint64(cmd)
	tempSum += uint64(dataLen)

	// Add all bytes together
	for _, v := range data {
		tempSum += uint64(v)
	}

	// Add magic value
	tempSum += 173

	// Truncate the sum to a single byte
	return uint8(tempSum), nil
}

// verifyChecksum computes a checksum over the packet's cmd, dataLen field and payload.
func verifyChecksum(in uint8, cmd uint8, dataLen uint8, data []byte) (bool, error) {

	cksum, err := calculateChecksum(cmd, dataLen, data)
	if err != nil {
		return false, err
	}

	return cksum == in, nil
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
