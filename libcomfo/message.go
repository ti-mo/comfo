package libcomfo

import "io"

// SetRequest is the interface that set requests to the unit
// must satisfy, as they need to be binary-encoded over the wire.
type SetRequest interface {
	MarshalBinary() ([]byte, error)
	Type() setRequest
}

// Response is the interface that response command
// implementations need to satisfy. They only need to be decoded
// on the client and don't need to be marshaled again.
type Response interface {
	UnmarshalBinary(in []byte) error // Interprets the Packet payload into a struct
	New() Response                   // Yields a freshly-allocated instance of the Response type
}

// EncodeGetRequest generates a Packet with an empty payload
// given the getRequest command type. A response payload is expected.
func EncodeGetRequest(gr getRequest) (out Packet) {

	out.Command = uint8(gr)

	// A get request does not have a request payload,
	// so we expect a response payload.
	out.Expect = true

	return
}

// EncodeSetRequest generates a Packet from a structure implementing
// the SetRequest interface. The Packet.Expect bool is only enabled
// when the request is empty (expects a response payload).
func EncodeSetRequest(in SetRequest) (out Packet, err error) {

	// Get the type from the SetRequest structure
	cmd := in.Type()
	if cmd == 0 {
		return out, errRequestType
	}

	// Get the binary representation of the request
	bin, err := in.MarshalBinary()
	if err != nil {
		return out, err
	}

	// Populate output structure
	out.Command = uint8(cmd)
	out.Data = bin

	// All empty requests expect a response payload (get operations),
	// but if the __request__ has a payload (a set operation),
	// we stop reading the response after the first ACK. (Expect false)
	out.Expect = false

	return
}

// DecodeResponse generates a Response from an incoming Packet
// from the unit. The ResponseType map is used as a type generator/translator
// for incoming command types.
func DecodeResponse(in Packet) (out Response, err error) {

	// Error out if we know the response cannot be decoded
	if in.Command == 0 {
		return nil, errDecodeNil
	}

	// Look up the response type in the ResponseType map,
	// do not unmarshal if the entry does not exist.
	out = ResponseType[in.Command]
	if out == nil {
		return nil, errResponseType
	}

	// Replace looked-up struct with a new instance
	out = out.New()

	err = out.UnmarshalBinary(in.Data)

	return
}

// setQuery is the end-to-end glue method for encoding a setRequest,
// sending it and returning the success status.
func setQuery(in SetRequest, conn io.ReadWriter) (err error) {

	sr, err := EncodeSetRequest(in)
	if err != nil {
		return
	}

	_, err = QueryPacket(sr, conn)

	return
}

// getQuery is the end-to-end glue method for encoding a getRequest,
// sending it and returning a Response.
func getQuery(in getRequest, conn io.ReadWriter) (resp Response, err error) {

	gr := EncodeGetRequest(in)

	qp, err := QueryPacket(gr, conn)
	if err != nil {
		return
	}

	resp, err = DecodeResponse(qp)

	return
}

// Type temperature represents a temperature data point.
// The binary representation is not a float, so an algorithm
// is in place to make the conversion.
type temperature float32

// MarshalBinary marshals the float32-derived temperature type
// into its binary representation. (after calculations)
func (t temperature) MarshalBinary() (out byte, err error) {

	// Wire format is (temperature + 20) * 2,
	// so 107 is the highest value that fits in one byte.
	// -20 is the lowest value that does not wrap.
	if int(t) > 107 {
		return out, errTooHigh
	}
	if int(t) < -20 {
		return out, errTooLow
	}

	return byte((t + 20) * 2), nil
}

// UnmarshalBinary unmarshals a binary temperature representation
// into a float32-derived temperature type.
func (t *temperature) UnmarshalBinary(in byte) {
	*t = temperature(in/2) - 20
	return
}

// leftPad32 pads a byte sequence to 32 bytes on the left.
func leftPad32(in []byte) (out []byte) {

	if len(in) > 4 {
		panic("input length to leftPad32 exceeds 4 bytes")
	}

	left := make([]byte, 4-len(in))
	out = append(left, in...)

	return
}
