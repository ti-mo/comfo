package main

var (
	ack         = []byte{0x07, 0xF3}
	pktStart    = []byte{0x07, 0xF0}
	pktEnd      = []byte{0x07, 0x0F}
	pktStartLen = len(pktStart)
	pktEndLen   = len(pktEnd)
)

type Packet struct {
	Command uint8
	Data    []byte
}

// MarshalPacket takes a Packet structure and converts it into wire protocol.
func MarshalPacket(in Packet) (out []byte, err error) {

	return
}

// UnmarshalPacket parses a byte slice of wire protocol into a Packet structure.
// The byte slice must start with pktStart and end with pktEnd,
// or the operation will fail.
func UnmarshalPacket(in []byte) (out Packet, err error) {
	if len(in) < 8 {
		return out, errPktLen
	}

	// Offsets of this packet's start and end delimiters
	var startOffset int = pktStartLen
	var endOffset int = len(in) - pktEndLen

	// Detect packet start and end
	if !byteCmp(in[:startOffset], pktStart) || !byteCmp(in[endOffset:], pktEnd) {
		return out, errDelim
	}

	// The size of the payload
	// TODO: Bump this by 1 for every double 0x07 BEFORE copying payload
	var dataLen uint8 = in[4]

	// Compare the expected size of the payload to the size of the packet.
	// There are 8 non-payload bytes in a packet, so the size cannot exceed this.
	if int(dataLen) > len(in)-8 {
		return out, errPayloadSize
	}

	// Verify checksum, extract command and payload
	if cksum, err := verifyChecksum(in); cksum && err != nil {
		out.Command = in[3]
		out.Data = in[5 : 5+dataLen]
	} else if !cksum {
		return out, errVerifyChecksum
	} else if err != nil {
		return out, errChecksum
	}

	return out, nil
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
