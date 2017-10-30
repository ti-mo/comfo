package libcomfo

import (
	"bytes"
	"errors"
	"fmt"
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

func (sr mockSetReq) Type() setRequest { return sr.mockType }
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
		err  error
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
