package libcomfo

import "fmt"

var (
	esc   byte = 0x07
	ack   byte = 0xF3
	start byte = 0xF0
	end   byte = 0x0F

	pktStart = []byte{esc, start}
	pktEnd   = []byte{esc, end}
	pktAck   = []byte{esc, ack}

	// byte 1-2 - command
	cmdOffset = 0
	// byte 3 - payload length
	lenOffset = cmdOffset + 2
	// byte 4-x - payload
	dataOffset = lenOffset + 1
)

type Packet struct {
	// Command type or response code
	Command uint8

	// Payload byte slice
	Data []byte

	// Whether to expect a response
	Expect bool
}

func (p Packet) String() string {
	return fmt.Sprintf("<Packet: %x %v, %v>", p.Command, p.Data, p.Expect)
}

// MarshalPacket takes a Packet structure and converts it into wire protocol.
func MarshalPacket(in Packet) (out []byte, err error) {

	// Save the original, unescaped length of the payload
	dataLen := len(in.Data)

	// Packet length specifier is one byte long, so 255 is the maximum packet size.
	if dataLen > 255 {
		return nil, errTooLong
	}

	// Escape the payload and re-calculate (wire) length
	data := escapeData(in.Data)
	dataLenWire := len(data)

	// Make output slice with payload length + 4 (all other bytes)
	out = make([]byte, dataLenWire+4)
	outLen := dataLenWire + 4

	copy(out[cmdOffset:lenOffset], []byte{0x00, in.Command}) // Command
	out[lenOffset] = uint8(dataLen)                          // Payload length
	copy(out[dataOffset:dataOffset+dataLenWire], data)       // Payload

	// Set checksum (last byte)
	out[outLen-1] = calculateChecksum(in.Command, uint8(dataLen), in.Data)

	return out, nil
}

// UnmarshalPacket parses a byte slice of wire protocol into a Packet structure.
func UnmarshalPacket(in []byte) (out Packet, err error) {
	if len(in) < 4 {
		return out, errPktLen
	}

	// The command
	var cmd uint8 = in[cmdOffset+1]

	// The size of the payload
	var dataLen uint8 = in[lenOffset]

	// Get checksum from packet (last byte)
	var cksum uint8 = in[len(in)-1]

	data := unescapeData(in[dataOffset : len(in)-1])

	// Compare the expected size of the payload to the size of the packet.
	// There are 4 non-payload bytes in a packet, so we can reliably calculate
	// how many bytes the payload *should* be.
	// Escaped 0x07's do not count towards dataLen, so needs to be unescaped first.
	if int(dataLen) != len(data) {
		return out, errPayloadSize
	}

	if !verifyChecksum(cksum, cmd, dataLen, data) {
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
		if in[i] == 0x07 {
			// 0x07 lookahead with bounds check
			if i+1 < len(in) && in[i+1] == 0x07 {
				// Only append a single 0x07
				out = append(out, 0x07)

				// Advance the window twice when successfully unescaping
				// to prevent re-evaluating on the second 0x07.
				i++
			}
		} else {
			out = append(out, in[i])
		}
	}

	return out
}

// calculateChecksum calculates the checksum of a Packet command,
// length and payload without (!) escaped 7s.
func calculateChecksum(cmd uint8, dataLen uint8, data []byte) uint8 {

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
	return uint8(tempSum)
}

// verifyChecksum computes a checksum over the packet's cmd, dataLen field and payload.
func verifyChecksum(in uint8, cmd uint8, dataLen uint8, data []byte) bool {

	return in == calculateChecksum(cmd, dataLen, data)
}
