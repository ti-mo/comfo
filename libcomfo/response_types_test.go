package libcomfo

import (
	"reflect"
	"testing"
)

// Unmarshal-specific corner cases and error checking.
func TestAll_UnmarshalBinary(t *testing.T) {
	tests := []struct {
		name  string
		b     []byte
		rtype uint8
		resp  Response
		err   error
	}{
		{
			name: "Fans",
			b:    []byte{33, 66, 0x04, 0xE2, 0x02, 0x71},
			resp: &Fans{
				InPercent:  33,
				OutPercent: 66,
				InSpeed:    1500,
				OutSpeed:   3000,
			},
			rtype: 0x0C,
		},
		{
			name:  "Fans zero fan speed",
			b:     []byte{33, 66, 0, 0, 0, 0},
			err:   errZeroValue,
			rtype: 0x0C,
		},
		{
			name:  "Fans incorrect length",
			b:     make([]byte, 0),
			err:   errPktLen,
			rtype: 0x0C,
		},
		{
			name: "BootloaderVersion",
			b:    append([]byte{1, 2, 3}, "abcdefghij"...),
			resp: &BootInfo{
				MajorVersion: 1,
				MinorVersion: 2,
				BetaVersion:  3,
				DeviceName:   "abcdefghij",
			},
			rtype: 0x68,
		},
		{
			name:  "BootloaderVersion incorrect length zero",
			b:     make([]byte, 0),
			rtype: 0x68,
			err:   errPktLen,
		},
		{
			name: "FirmwareVersion",
			b:    append([]byte{1, 2, 3}, "qrstuvwxyz"...),
			resp: &BootInfo{
				MajorVersion: 1,
				MinorVersion: 2,
				BetaVersion:  3,
				DeviceName:   "qrstuvwxyz",
			},
			rtype: 0x6A,
		},
		{
			name:  "FirmwareVersion incorrect length zero",
			b:     make([]byte, 0),
			rtype: 0x6A,
			err:   errPktLen,
		},
		{
			name: "FanProfiles",
			b:    []byte{5, 10, 15, 30, 35, 40, 50, 55, 60, 1, 20, 45, 0xFF, 0xFF},
			resp: &FanProfiles{
				OutAway:      5,
				OutLow:       10,
				OutMid:       15,
				OutHigh:      20,
				InFanActive:  true,
				InAway:       30,
				InLow:        35,
				InMid:        40,
				InHigh:       45,
				CurrentOut:   50,
				CurrentIn:    55,
				CurrentLevel: 60,
			},
			rtype: 0xCE,
		},
		{
			name:  "FanProfiles incorrect length",
			b:     make([]byte, 0),
			err:   errPktLen,
			rtype: 0xCE,
		},
		{
			name: "Temps",
			b:    []byte{66, 114, 66, 114, 66, 0xFF, 114, 66, 114},
			resp: &Temps{
				Comfort:     13,
				OutsideAir:  37,
				SupplyAir:   13,
				OutAir:      37,
				ExhaustAir:  13,
				GeoHeat:     37,
				Reheating:   13,
				KitchenHood: 37,
			},
			rtype: 0xD2,
		},
		{
			name:  "Temps incorrect length",
			b:     make([]byte, 0),
			err:   errPktLen,
			rtype: 0xD2,
		},
		{
			name: "Hours",
			b: []byte{
				0xAA, 0x55, 0x55, 0xAA, 0xA5, 0xA5, 0x5A, 0x5A,
				0xA5, 0x5A, 0x5A, 0xA5, 0xAA, 0x55, 0x55, 0xAA,
				0xBB, 0xCC, 0x11, 0x22,
			},
			resp: &Hours{
				FanAway:      11162965,
				FanLow:       11183525,
				FanMid:       5921445,
				FanHigh:      13373730,
				FrostProtect: 23130,
				Reheating:    42410,
				BypassOpen:   21845,
				Filter:       43707,
			},
			rtype: 0xDE,
		},
		{
			name:  "Hours incorrect length",
			b:     make([]byte, 0),
			err:   errPktLen,
			rtype: 0xDE,
		},
		{
			name: "Bypass",
			b:    []byte{0xFF, 0xFF, 0xF0, 0x90, 0x01, 0xFF, 1},
			resp: &Bypass{
				Factor:     240,
				Level:      144,
				Correction: 1,
				SummerMode: true,
			},
			rtype: 0xE0,
		},
		{
			name:  "Bypass incorrect length",
			b:     make([]byte, 0),
			err:   errPktLen,
			rtype: 0xE0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			// Get correct response type
			resp := ResponseType[tt.rtype].New()

			// Initialize Response field to the correct type so an empty struct
			// is not compared to an empty interface (<nil>) on error.
			if tt.resp == nil {
				tt.resp = resp.New()
			}

			if want, got := tt.err, resp.UnmarshalBinary(tt.b); want != got {
				t.Fatalf("unexpected error unmarshaling response:\n- want: %v\n-  got: %v",
					want, got)
			}

			if want, got := tt.resp, resp; !reflect.DeepEqual(want, got) {
				t.Fatalf("unexpected response unmarshal result:\n- want: %v\n-  got: %v",
					want, got)
			}
		})
	}
}
