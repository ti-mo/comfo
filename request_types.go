package main

// Type getRequest is an enum holding the command types
// of all Get operations.
type getRequest uint8

//
// get* requests
// Get information from the unit, no state is altered.
//

const (
	getTemps  getRequest = 0xD1 // getTemps gets all available temperatures from the unit.
	getBypass            = 0xDF // getBypass gets heat exchanger information from the unit.
)

//
// set* requests
// Alter unit state or configuration
//

// SetRequest setSpeed controls the fan speed of the unit.
// The Speed member can range from 0 (away) to 4 (highest).
type setSpeed struct {
	Speed uint8
}

func (q setSpeed) Type() uint8 { return 0x99 }
func (q setSpeed) MarshalBinary() (out []byte, err error) {

	if q.Speed > 4 {
		return nil, errTooHigh
	}

	return []byte{q.Speed}, nil
}

// SetRequest setComfort controls the comfort temperature
// of the unit. This is the point at which the heat exchanger
// will recover outgoing heat into the incoming air stream.
type setComfort struct {
	Temperature temperature
}

func (q setComfort) Type() uint8 { return 0xD3 }
func (q setComfort) MarshalBinary() (out []byte, err error) {

	t, err := q.Temperature.MarshalBinary()
	if err != nil {
		return
	}

	return []byte{t}, nil
}
