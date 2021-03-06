package libcomfo

import (
	"encoding/binary"
	"strings"
)

var (
	// ResponseType maps incoming response types to
	// their internal structs. Call New() on a Response
	// to get a clean copy that can be safely modified.
	ResponseType = map[uint8]Response{
		0x0C: &Fans{},
		0x68: &BootInfo{},
		0x6A: &BootInfo{},
		0xCE: &FanProfiles{},
		0xD2: &Temps{},
		0xDA: &Errors{},
		0xDE: &Hours{},
		0xE0: &Bypass{},
	}
)

// BootInfo holds bootloader/firmware-related info.
type BootInfo struct {
	MajorVersion uint8  `json:"major_version"`
	MinorVersion uint8  `json:"minor_version"`
	BetaVersion  uint8  `json:"beta_version"`
	DeviceName   string `json:"device_name"`
}

// New returns a new instance of BootInfo.
func (bi *BootInfo) New() Response { return &BootInfo{} }

// UnmarshalBinary unmarshals the binary representation
// into a BootInfo structure. Whitespace is trimmed from DeviceName.
func (bi *BootInfo) UnmarshalBinary(in []byte) error {

	if len(in) != 13 {
		return errPktLen
	}

	bi.MajorVersion = in[0]
	bi.MinorVersion = in[1]
	bi.BetaVersion = in[2]
	bi.DeviceName = strings.TrimSpace(string(in[3:]))

	return nil
}

// Temps holds the various temperature readings
// from the ventilation unit.
type Temps struct {
	Comfort     temperature `json:"comfort"`
	OutsideAir  temperature `json:"outside_air"`
	SupplyAir   temperature `json:"supply_air"`
	InsideAir   temperature `json:"inside_air"`
	ExhaustAir  temperature `json:"exhaust_air"`
	GeoHeat     temperature `json:"geo_heat"`
	Reheating   temperature `json:"reheating"`
	KitchenHood temperature `json:"kitchen_hood"`
}

// New returns a new instance of Temps.
func (*Temps) New() Response { return &Temps{} }

// UnmarshalBinary unmarshals the binary representation
// into a Temps structure. Fixed length is 9 bytes.
func (t *Temps) UnmarshalBinary(in []byte) error {

	if len(in) != 9 {
		return errPktLen
	}

	t.Comfort.UnmarshalBinary(in[0])
	t.OutsideAir.UnmarshalBinary(in[1])
	t.SupplyAir.UnmarshalBinary(in[2])
	t.InsideAir.UnmarshalBinary(in[3])
	t.ExhaustAir.UnmarshalBinary(in[4])
	t.GeoHeat.UnmarshalBinary(in[6])
	t.Reheating.UnmarshalBinary(in[7])
	t.KitchenHood.UnmarshalBinary(in[8])

	return nil
}

// Bypass holds the information about
// the unit's heat exchanger bypass valve.
type Bypass struct {
	Factor     uint8 `json:"factor"`
	Level      uint8 `json:"level"`
	Correction uint8 `json:"correction"`
	SummerMode bool  `json:"summer_mode"`
}

// New returns a new instance of Bypass.
func (*Bypass) New() Response { return &Bypass{} }

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

// Hours holds the amount of working hours
// for every moving component in the unit.
type Hours struct {
	FanAway      uint32 `json:"fan_away"`
	FanLow       uint32 `json:"fan_low"`
	FanMid       uint32 `json:"fan_mid"`
	FanHigh      uint32 `json:"fan_high"`
	FrostProtect uint16 `json:"frost_protect"`
	Reheating    uint16 `json:"reheating"`
	BypassOpen   uint16 `json:"bypass_open"`
	Filter       uint16 `json:"filter"`
}

// New returns a new instance of Hours.
func (*Hours) New() Response { return &Hours{} }

// UnmarshalBinary unmarshals the binary representation
// into an Hours structure.
func (h *Hours) UnmarshalBinary(in []byte) error {

	if len(in) != 20 {
		return errPktLen
	}

	h.FanAway = binary.BigEndian.Uint32(leftPad32(in[0:3]))
	h.FanLow = binary.BigEndian.Uint32(leftPad32(in[3:6]))
	h.FanMid = binary.BigEndian.Uint32(leftPad32(in[6:9]))
	h.FrostProtect = binary.BigEndian.Uint16(in[9:11])
	h.Reheating = binary.BigEndian.Uint16(in[11:13])
	h.BypassOpen = binary.BigEndian.Uint16(in[13:15])
	h.Filter = binary.BigEndian.Uint16(in[15:17])
	h.FanHigh = binary.BigEndian.Uint32(leftPad32(in[17:20]))

	return nil
}

// Fans holds the unit's fan percentage and speeds.
type Fans struct {
	InPercent  uint8  `json:"in_percent"`
	OutPercent uint8  `json:"out_percent"`
	InSpeed    uint16 `json:"in_speed"`
	OutSpeed   uint16 `json:"out_speed"`
}

// New returns a new instance of Fans.
func (*Fans) New() Response { return &Fans{} }

// UnmarshalBinary unmarshals the binary representation
// into a Fans structure.
func (v *Fans) UnmarshalBinary(in []byte) error {

	if len(in) != 6 {
		return errPktLen
	}

	// Sanity check both values to catch divide by zero
	inSpeed := binary.BigEndian.Uint16(in[2:4])
	outSpeed := binary.BigEndian.Uint16(in[4:6])

	if inSpeed == 0 || outSpeed == 0 {
		return errZeroValue
	}

	v.InPercent = uint8(in[0])
	v.OutPercent = uint8(in[1])

	v.InSpeed = uint16(1875000 / uint32(inSpeed))
	v.OutSpeed = uint16(1875000 / uint32(outSpeed))

	return nil
}

// FanProfiles holds the fan profiles (in percent)
// for every ventilation mode.
type FanProfiles struct {
	OutAway uint8 `json:"out_away"`
	OutLow  uint8 `json:"out_low"`
	OutMid  uint8 `json:"out_mid"`
	OutHigh uint8 `json:"out_high"`

	InFanActive bool  `json:"in_fan_active"`
	InAway      uint8 `json:"in_away"`
	InLow       uint8 `json:"in_low"`
	InMid       uint8 `json:"in_mid"`
	InHigh      uint8 `json:"in_high"`

	CurrentOut  uint8 `json:"current_out"`
	CurrentIn   uint8 `json:"current_in"`
	CurrentMode uint8 `json:"current_mode"`
}

// New returns a new instance of FanProfiles.
func (*FanProfiles) New() Response { return &FanProfiles{} }

// UnmarshalBinary unmarshals the binary representation
// into a FanProfiles structure.
func (vp *FanProfiles) UnmarshalBinary(in []byte) error {

	if len(in) != 14 {
		return errPktLen
	}

	vp.OutAway = in[0]
	vp.OutLow = in[1]
	vp.OutMid = in[2]
	vp.OutHigh = in[10]

	vp.InAway = in[3]
	vp.InLow = in[4]
	vp.InMid = in[5]
	vp.InHigh = in[11]
	vp.InFanActive = in[9] == 1

	vp.CurrentOut = in[6]
	vp.CurrentIn = in[7]
	vp.CurrentMode = in[8]

	return nil
}

// Lookup returns the in and out airflow percentages associated with the given
// speed profile (integer). Only values 1-4 are allowed.
func (vp FanProfiles) Lookup(speed uint8) (inPercent uint8, outPercent uint8, err error) {

	switch speed {
	case 1:
		return vp.InAway, vp.OutAway, nil
	case 2:
		return vp.InLow, vp.OutLow, nil
	case 3:
		return vp.InMid, vp.OutMid, nil
	case 4:
		return vp.InHigh, vp.OutHigh, nil
	default:
		return 0, 0, errNotExist
	}
}

// Errors holds a list of errors that can occur
// within the unit and their status.
type Errors struct {
	Filter bool `json:"filter"`
}

// New returns a new instance of Errors.
func (e *Errors) New() Response { return &Errors{} }

// UnmarshalBinary unmarshals the binary representation
// into a Errors structure.
func (e *Errors) UnmarshalBinary(in []byte) error {

	if len(in) != 17 {
		return errPktLen
	}

	e.Filter = in[8] == 1

	return nil
}
