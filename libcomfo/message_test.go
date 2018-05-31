package libcomfo

import (
	"bytes"
	"errors"
	"reflect"
	"testing"
)

// mockSetReq implements the SetRequest interface
// and is used to raise errors during testing.
type mockSetReq struct {
	mockType setRequest
	mockData []byte
	mockErr  error
}

func (sr mockSetReq) requestType() setRequest { return sr.mockType }
func (sr mockSetReq) MarshalBinary() (out []byte, err error) {
	return sr.mockData, sr.mockErr
}

func TestTemperature_MarshalBinary(t *testing.T) {
	tests := []struct {
		name string
		temp temperature
		b    byte
		err  error
	}{
		{
			name: "average value",
			temp: temperature(25),
			b:    90,
		},
		{
			name: "high value",
			temp: temperature(107),
			b:    254,
		},
		{
			name: "low value",
			temp: temperature(-20),
			b:    0,
		},
		{
			name: "temperature too low",
			temp: temperature(-21),
			err:  errTooLow,
		},
		{
			name: "temperature too high",
			temp: temperature(108),
			err:  errTooHigh,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			// Marshal value into binary representation
			b, err := tt.temp.MarshalBinary()

			if want, got := tt.err, err; want != got {
				t.Fatalf("unexpected error marshaling temperature:\n- want: %v\n-  got: %v",
					want, got)
			}

			if want, got := tt.b, b; !reflect.DeepEqual(want, got) {
				t.Fatalf("unexpected temperature marshal:\n- want: %v\n-  got: %v",
					want, got)
			}
		})
	}
}

func TestTemperature_UnmarshalBinary(t *testing.T) {
	tests := []struct {
		name string
		temp temperature
		b    byte
	}{
		{
			name: "average value",
			b:    90,
			temp: temperature(25),
		},
		{
			name: "high value",
			b:    255,
			temp: temperature(107),
		},
		{
			name: "low value",
			b:    0,
			temp: temperature(-20),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			// Unmarshal binary representation into temperature value
			var temp temperature
			temp.UnmarshalBinary(tt.b)

			if want, got := tt.temp, temp; !reflect.DeepEqual(want, got) {
				t.Fatalf("unexpected temperature unmarshal:\n- want: %v\n-  got: %v",
					want, got)
			}
		})
	}
}

// Simple test that calls EncodeGetRequest on a getRequest,
// which wraps an int in a struct and sets Expect to true.
func TestEncodeGetRequest(t *testing.T) {

	pkt := EncodeGetRequest(getBootloader)

	if want, got := getBootloader, getRequest(pkt.Command); want != got {
		t.Fatalf("unexpected Packet command:\n- want: %v\n-  got: %v", want, got)
	}

	if want, got := true, pkt.Expect; want != got {
		t.Fatalf("unexpected Packet.Expect:\n- want: %v\n-  got: %v", want, got)
	}
}

func TestEncodeSetRequest(t *testing.T) {

	errTestMarshal := errors.New("error during marshal")

	tests := []struct {
		name string
		sr   SetRequest
		pkt  Packet
		err  error
	}{
		{
			name: "setSpeedT",
			sr:   setSpeedT{Speed: 3},
			pkt:  Packet{Command: 0x99, Data: []byte{3}, Expect: false},
		},
		{
			name: "mockSetReq unknown type",
			sr:   mockSetReq{mockType: 0},
			err:  errRequestType,
		},
		{
			name: "mockSetReq marshal error",
			sr: mockSetReq{
				mockType: 0xff,
				mockErr:  errTestMarshal,
			},
			err: errTestMarshal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			pkt, err := EncodeSetRequest(tt.sr)

			if want, got := tt.err, err; want != got {
				t.Fatalf("unexpected error encoding request:\n- want: %v\n-  got: %v",
					want, got)
			}

			if want, got := tt.pkt, pkt; !reflect.DeepEqual(want, got) {
				t.Fatalf("unexpected request encoding result:\n- want: %v\n-  got: %v",
					want, got)
			}
		})
	}
}

func TestDecodeResponse(t *testing.T) {

	tests := []struct {
		name string
		pkt  Packet
		resp Response
		err  error
	}{
		{
			name: "Bypass",
			pkt:  Packet{Command: 0xE0, Data: []byte{0xFF, 0xFF, 0xF0, 0x90, 0x01, 0xFF, 1}},
			resp: &Bypass{
				Factor:     240,
				Level:      144,
				Correction: 1,
				SummerMode: true,
			},
		},
		{
			name: "response with nil command",
			pkt:  Packet{},
			err:  errDecodeNil,
		},
		{
			name: "non-existing response type",
			pkt:  Packet{Command: 0xFF},
			err:  errResponseType,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			resp, err := DecodeResponse(tt.pkt)

			if want, got := tt.err, err; want != got {
				t.Fatalf("unexpected error decoding response:\n- want: %v\n-  got: %v",
					want, got)
			}

			if want, got := tt.resp, resp; !reflect.DeepEqual(want, got) {
				t.Fatalf("unexpected response decoding result:\n- want: %v\n-  got: %v",
					want, got)
			}
		})
	}
}

func TestConnSetRequest(t *testing.T) {

	srtests := []struct {
		name string
		sr   SetRequest
		tc   TestConn
		err  error
	}{
		{
			name: "set comfort temperature",
			sr:   setComfortT{Temperature: 42},
			tc: TestConn{
				Receives: []byte{
					0x07, 0xF0, // Frame Start
					0x00, uint8(setComfort), // request type
					0x01, 0x7C, // Length + payload (calculated temp)
					0xFD,     // Checksum
					esc, end, // Frame End
					esc, ack, // Ack
				},
				Emits: []byte{esc, ack}, // Read an ACK back from the connection
			},
		},
		{
			name: "set fan speed mode",
			sr:   setSpeedT{Speed: 3},
			tc: TestConn{
				Receives: []byte{
					0x07, 0xF0, // Start
					0x00, uint8(setSpeed), // request type
					0x01, 0x03, // Length + payload
					0x4A, // Checksum
					esc, end,
					esc, ack,
				},
				Emits: []byte{esc, ack},
			},
		},
		{
			name: "garbage before response ACK",
			sr:   setComfortT{Temperature: 42},
			tc: TestConn{
				Receives: []byte{
					0x07, 0xF0, // Start
					0x00, uint8(setComfort), // request type
					0x01, 0x7C, // Length + payload
					0xFD, // Checksum
					esc, end,
					esc, ack,
				},
				Emits: []byte{0xF0, 0x0B, 0xA1, 0x2, esc, ack},
			},
			err: errScanInput,
		},
		{
			name: "error from EncodeSetRequest",
			sr:   mockSetReq{mockType: 0x0},
			tc:   TestConn{},
			err:  errRequestType,
		},
	}

	for _, tt := range srtests {
		t.Run(tt.name, func(t *testing.T) {

			// Run setQuery against a test connection
			err := setQuery(tt.sr, &tt.tc)

			if want, got := tt.err, err; want != got {
				t.Fatalf("unexpected error during request:\n- want: %v\n-  got: %v",
					want, got)
			}
		})
	}
}

func TestConnGetRequest(t *testing.T) {

	grtests := []struct {
		name string
		gr   getRequest // Command type only
		resp Response
		tc   TestConn
		err  bool
	}{
		{
			name: "get fan speeds",
			gr:   getFans,
			tc: TestConn{
				Receives: []byte{ // Request
					0x07, 0xF0, // Frame Start
					0x00, uint8(getFans), // request type
					0x00,     // Length
					0xB8,     // Checksum
					esc, end, // Frame End
					esc, ack, // Response ACK
				},
				Emits: []byte{ // Response
					esc, ack,
					0x07, 0xF0, // Frame Start
					0x00, uint8(getFans + 1), // response type
					0x06,       // Length
					0xAA, 0xBB, // in/out percents
					0x11, 0x22, // in speed
					0x33, 0x44, // out speed
					0xCE, // Checksum
					esc, end,
					esc, ack,
				},
			},
			resp: &Fans{
				InPercent:  170,
				OutPercent: 187,
				InSpeed:    427,
				OutSpeed:   142,
			},
		},
		{
			name: "connection write error",
			gr:   getFans,
			tc:   TestConn{},
			err:  true,
		},
	}

	for _, tt := range grtests {
		t.Run(tt.name, func(t *testing.T) {

			// Run setQuery against a test connection
			resp, err := getQuery(tt.gr, &tt.tc)

			// Generic error checking
			if want, got := tt.err, err; !tt.err && err != nil {
				t.Fatalf("unexpected error during request:\n- want: %v\n-  got: %v",
					want, got)
			}

			if want, got := tt.resp, resp; !reflect.DeepEqual(want, got) {
				t.Fatalf("unexpected response:\n- want: %v\n-  got: %v",
					want, got)
			}
		})
	}
}

func TestLeftPad32(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("unable to recover from panic in leftPad32()")
		}
	}()

	want := []byte{0, 0, 0x1, 0x2}
	got := leftPad32([]byte{0x01, 0x02})

	if !bytes.Equal(want, got) {
		t.Fatalf("unexpected error in leftPad32t:\n- want: %v\n-  got: %v",
			want, got)
	}

	// Make leftPad32 panic with oversized input
	leftPad32([]byte{1, 2, 3, 4, 5})
}
