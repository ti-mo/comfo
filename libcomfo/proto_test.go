package libcomfo

import (
	"reflect"
	"testing"
)

func genBytes(len int) (out []byte) {
	out = make([]byte, len)

	for i := 0; i < len; i++ {
		out[i] = uint8(i)
	}

	return
}

// End-to-end unmarshal-marshal tests.
func TestFull(t *testing.T) {
	tests := []struct {
		name string
		pkt  Packet
		b    []byte
		err  error
	}{
		{
			name: "get rf status (empty payload)",
			b:    []byte{0x00, 0xE5, 0x00, 0x92},
			pkt:  Packet{Command: 0xE5, Data: []byte{}},
		},
		{
			name: "escape 07",
			b:    []byte{0x00, 0xAB, 0x03, 0x7, 0x7, 0x1, 0x2, 0x65},
			pkt:  Packet{Command: 0xAB, Data: []byte{0x7, 0x1, 0x2}},
		},
		{
			name: "escape double 07",
			b:    []byte{0x00, 0xAB, 0x04, 0x7, 0x7, 0x7, 0x7, 0x1, 0x2, 0x6D},
			pkt:  Packet{Command: 0xAB, Data: []byte{0x7, 0x7, 0x1, 0x2}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			// Unmarshal binary content into nested structures
			pkt, err := UnmarshalPacket(tt.b)

			if want, got := tt.err, err; want != got {
				t.Fatalf("unexpected error unmarshaling packet:\n- want: %v\n-  got: %v",
					want, got)
			}

			if want, got := tt.pkt, pkt; !reflect.DeepEqual(want, got) {
				t.Fatalf("unexpected packet unmarshal result:\n- want: %v\n-  got: %v",
					want, got)
			}

			var b []byte

			// Attempt re-marshal into binary form
			b, err = MarshalPacket(pkt)

			if want, got := tt.err, err; want != got {
				t.Fatalf("unexpected error marshaling packet:\n- want: %v\n-  got: %v",
					want, got)
			}

			if want, got := tt.b, b; !reflect.DeepEqual(want, got) {
				t.Fatalf("unexpected packet marshal result:\n- want: %v\n-  got: %v",
					want, got)
			}
		})
	}
}

// Unmarshal-specific corner cases and error checking.
func TestUnmarshalPacket(t *testing.T) {
	tests := []struct {
		name string
		pkt  Packet
		b    []byte
		err  error
	}{
		{
			name: "input too short",
			b:    []byte{0x00, 0x00, 0x00},
			err:  errPktLen,
		},
		{
			name: "input empty",
			b:    []byte{0x00, 0x00, 0x00},
			err:  errPktLen,
		},
		{
			name: "incorrect length specifier (larger than payload)",
			b:    []byte{0x00, 0x00, 0xff, 0x00},
			err:  errPayloadSize,
		},
		{
			name: "incorrect length specifier (includes escaped 7s)",
			b:    []byte{0x00, 0x00, 0x02, 0x07, 0x07, 0x00},
			err:  errPayloadSize,
		},
		{
			name: "correct length specifier and checksum (escaped 7s)",
			b:    []byte{0x00, 0x00, 0x01, 0x07, 0x07, 0xB5},
			pkt:  Packet{Command: 0, Data: []byte{0x07}},
		},
		{
			name: "incorrect checksum",
			b:    []byte{0x00, 0x00, 0x00, 0x00},
			err:  errVerifyChecksum,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			// Unmarshal binary content into nested structures
			pkt, err := UnmarshalPacket(tt.b)

			if want, got := tt.err, err; want != got {
				t.Fatalf("unexpected error unmarshaling packet:\n- want: %v\n-  got: %v",
					want, got)
			}

			if want, got := tt.pkt, pkt; !reflect.DeepEqual(want, got) {
				t.Fatalf("unexpected packet unmarshal result:\n- want: %v\n-  got: %v",
					want, got)
			}
		})
	}
}

// Unmarshal-specific corner cases and error checking.
func TestMarshalPacket(t *testing.T) {
	tests := []struct {
		name string
		pkt  Packet
		b    []byte
		err  error
	}{
		{
			name: "input too long",
			pkt:  Packet{Data: genBytes(256)},
			err:  errTooLong,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			// Attempt re-marshal into binary form
			b, err := MarshalPacket(tt.pkt)

			if want, got := tt.err, err; want != got {
				t.Fatalf("unexpected error marshaling packet:\n- want: %v\n-  got: %v",
					want, got)
			}

			if want, got := tt.b, b; !reflect.DeepEqual(want, got) {
				t.Fatalf("unexpected packet marshal result:\n- want: %v\n-  got: %v",
					want, got)
			}
		})
	}
}

// Make sure string representation displays
// all fields in a clear, readable manner.
func TestPacket_String(t *testing.T) {
	pkt := Packet{Command: 0x42, Expect: true, Data: []byte{3, 6, 9, 12, 15, 18}}.String()

	exp := "<Packet {Command: 42, Data: [3 6 9 12 15 18], Expect: true}>"

	if want, got := exp, pkt; pkt != exp {
		t.Fatalf("unexpected packet string:\n- want: %v\n-  got: %v",
			want, got)
	}
}
