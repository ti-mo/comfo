package libcomfo

import (
	"reflect"
	"testing"
)

// Unmarshal-specific corner cases and error checking.
func TestSetRequests_MarshalBinary(t *testing.T) {
	tests := []struct {
		name  string
		rtype uint8
		req   SetRequest
		b     []byte
		err   error
	}{
		{
			name:  "setSpeed",
			rtype: 0x99,
			req: setSpeed{
				Speed: 2,
			},
			b: []byte{0x02},
		},
		{
			name:  "setSpeed too high",
			rtype: 0x99,
			req: setSpeed{
				Speed: 5,
			},
			err: errTooHigh,
		},
		{
			name:  "setComfort",
			rtype: 0xD3,
			req: setComfort{
				Temperature: 21,
			},
			b: []byte{0x52},
		},
		{
			name:  "setComfort too high",
			rtype: 0xD3,
			req: setComfort{
				Temperature: 108,
			},
			err: errTooHigh,
		},
		{
			name:  "setComfort too low",
			rtype: 0xD3,
			req: setComfort{
				Temperature: -21,
			},
			err: errTooLow,
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
