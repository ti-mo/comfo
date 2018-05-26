package libcomfo

import (
	"reflect"
	"testing"
)

// Unmarshal-specific corner cases and error checking.
func TestSetRequests_MarshalBinary(t *testing.T) {
	tests := []struct {
		name  string
		rtype setRequest
		req   SetRequest
		b     []byte
		err   error
	}{
		{
			name:  "setSpeedT",
			rtype: 0x99,
			req: setSpeedT{
				Speed: 2,
			},
			b: []byte{0x02},
		},
		{
			name:  "setSpeedT too high",
			rtype: 0x99,
			req: setSpeedT{
				Speed: 5,
			},
			err: errTooHigh,
		},
		{
			name:  "setComfortT",
			rtype: 0xD3,
			req: setComfortT{
				Temperature: 21,
			},
			b: []byte{0x52},
		},
		{
			name:  "setComfortT too high",
			rtype: 0xD3,
			req: setComfortT{
				Temperature: 108,
			},
			err: errTooHigh,
		},
		{
			name:  "setComfortT too low",
			rtype: 0xD3,
			req: setComfortT{
				Temperature: -21,
			},
			err: errTooLow,
		},
		{
			name:  "FanProfiles (bulk)",
			rtype: 0xCF,
			req: FanProfiles{
				OutAway: 10,
				InAway:  10,
				OutLow:  20,
				InLow:   20,
				OutMid:  30,
				InMid:   30,
				OutHigh: 40,
				InHigh:  40,
			},
			b: []byte{10, 20, 30, 10, 20, 30, 40, 40},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			// Marshal the request
			breq, err := tt.req.MarshalBinary()

			if want, got := tt.rtype, tt.req.Type(); want != got {
				t.Fatalf("unexpected request type:\n- want: %v\n-  got: %v",
					want, got)
			}

			if want, got := tt.err, err; want != got {
				t.Fatalf("unexpected error marshalling request:\n- want: %v\n-  got: %v",
					want, got)
			}

			if want, got := tt.b, breq; !reflect.DeepEqual(want, got) {
				t.Fatalf("unexpected request marshal result:\n- want: %v\n-  got: %v",
					want, got)
			}
		})
	}
}
