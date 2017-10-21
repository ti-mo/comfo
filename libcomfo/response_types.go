package libcomfo

import (
	"encoding/binary"
	"strings"
)

var (
	// Map incoming response types to their internal structs
	ResponseType = map[uint8]Response{
		0x0C: &Ventilators{},
		0x68: &BootInfo{},
		0x6A: &BootInfo{},
		0xD2: &Temps{},
		0xDE: &Hours{},
		0xE0: &Bypass{},
	}
)

// Type BootInfo holds bootloader/firmware-related info.
type BootInfo struct {
	MajorVersion uint8
	MinorVersion uint8
	BetaVersion  uint8
	DeviceName   string
}

// UnmarshalBinary unmarshals the binary representation
// into a BootInfo structure. Whitespace is trimmed from DeviceName.
func (bi *BootInfo) UnmarshalBinary(in []byte) error {
	if len(in) != 13 {
		return errPktLen
	}

	bi.MajorVersion = in[1]
	bi.MinorVersion = in[2]
	bi.BetaVersion = in[3]
	bi.DeviceName = strings.TrimSpace(string(in[3:13]))

	return nil
}

// Type Temps holds the various temperature readings
// from the ventilation unit.
type Temps struct {
	Comfort     temperature
	OutsideAir  temperature
	SupplyAir   temperature
	OutAir      temperature
	ExhaustAir  temperature
	GeoHeat     temperature
	Reheating   temperature
	KitchenHood temperature
}

// UnmarshalBinary unmarshals the binary representation
// into a Temps structure. Fixed length is 9 bytes.
func (t *Temps) UnmarshalBinary(in []byte) error {

	if len(in) != 9 {
		return errPktLen
	}

	t.Comfort.UnmarshalBinary(in[0])
	t.OutsideAir.UnmarshalBinary(in[1])
	t.SupplyAir.UnmarshalBinary(in[2])
	t.OutAir.UnmarshalBinary(in[3])
	t.ExhaustAir.UnmarshalBinary(in[4])
	t.GeoHeat.UnmarshalBinary(in[6])
	t.Reheating.UnmarshalBinary(in[7])
	t.KitchenHood.UnmarshalBinary(in[8])

	return nil
}

// Type Bypass holds the information about
// the unit's heat exchanger bypass valve.
type Bypass struct {
	Factor     uint8
	Level      uint8
	Correction uint8
	SummerMode bool
}

// UnmarshalBinary unmarshals the binary representation
// into a Bypass structure.
func (b *Bypass) UnmarshalBinary(in []byte) error {

	if len(in) != 7 {
		return errPktLen
	}

	b.Factor = in[2]
	b.Level = in[3]
	b.Correction = in[4]
	b.SummerMode = in[6] == 1

	return nil
}

// Type Hours holds the amount of working hours
// for every moving component in the unit.
type Hours struct {
	FanAway      uint32
	FanLow       uint32
	FanMid       uint32
	FanHigh      uint32
	FrostProtect uint16
	Reheating    uint16
	BypassOpen   uint16
	Filter       uint16
}

// UnmarshalBinary unmarshals the binary representation
// into an Hours structure.
func (h *Hours) UnmarshalBinary(in []byte) error {

	if len(in) != 20 {
		return errPktLen
	}

	h.FanAway = binary.BigEndian.Uint32(LeftPad32(in[0:3]))
	h.FanLow = binary.BigEndian.Uint32(LeftPad32(in[3:6]))
	h.FanMid = binary.BigEndian.Uint32(LeftPad32(in[6:9]))
	h.FrostProtect = binary.BigEndian.Uint16(in[9:11])
	h.Reheating = binary.BigEndian.Uint16(in[11:13])
	h.BypassOpen = binary.BigEndian.Uint16(in[13:15])
	h.Filter = binary.BigEndian.Uint16(in[15:17])
	h.FanHigh = binary.BigEndian.Uint32(LeftPad32(in[17:20]))

	return nil
}

// Type Ventilators holds the unit's fan percentage and speeds.
type Ventilators struct {
	InPercent  uint8
	OutPercent uint8
	InSpeed    uint16
	OutSpeed   uint16
}

// UnmarshalBinary unmarshals the binary representation
// into a Ventilators structure.
func (v *Ventilators) UnmarshalBinary(in []byte) error {

	if len(in) != 6 {
		return errPktLen
	}

	v.InPercent = uint8(in[0])
	v.OutPercent = uint8(in[1])

	v.InSpeed = uint16(1857000 / uint32(binary.BigEndian.Uint16(in[2:4])))
	v.OutSpeed = uint16(1857000 / uint32(binary.BigEndian.Uint16(in[4:6])))

	return nil
}
