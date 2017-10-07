package main

import (
	"reflect"
	"testing"
)

func TestPacket_Unmarshal(t *testing.T) {
	tests := []struct {
		name string
		pkt  Packet
		b    []byte
		err  error
	}{
		{
			name: "get rf status (empty payload)",
			b:    []byte{0x07, 0xF0, 0x00, 0xE5, 0x00, 0x92, 0x07, 0x0F},
			pkt:  Packet{Command: 0xE5, Data: []byte{}},
		},
		{
			name: "escape double 07",
			b:    []byte{0x07, 0xF0, 0x00, 0xAB, 0x03, 0x7, 0x7, 0x1, 0x2, 0xB8, 0x07, 0x0F},
			pkt:  Packet{Command: 0xAB, Data: []byte{0x7, 0x1, 0x2}},
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
