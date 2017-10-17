package main

var (
	// Map incoming response types to their internal structs
	ResponseType = map[uint8]Response{
		0xD2: &Temps{},
		0xE0: &Bypass{},
	}
)

// Temps holds the various temperature readings
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

// UnmarshalBinary unmarshals a Temps structure to the binary
// representation. Fixed length is 9 bytes.
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

type Bypass struct {
	Factor     uint8
	Level      uint8
	Correction uint8
	SummerMode bool
}

func (b *Bypass) UnmarshalBinary(in []byte) error {

	b.Factor = in[2]
	b.Level = in[3]
	b.Correction = in[4]
	b.SummerMode = in[6] == 1

	return nil
}
