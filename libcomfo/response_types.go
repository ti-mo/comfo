package libcomfo

var (
	// Map incoming response types to their internal structs
	ResponseType = map[uint8]Response{
		0xD2: &Temps{},
		0xE0: &Bypass{},
	}
)

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
