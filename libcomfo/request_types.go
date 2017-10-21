package libcomfo

// Type getRequest is an enum holding the command types
// of all Get operations.
type getRequest uint8

//
// get* requests
// Get information from the unit, no state is altered.
//

const (
	getVentilators  getRequest = 0x0B // getVentilators gets ventilator speed % and RPM
	getBootloader              = 0x67 // getBootloader gets bootloader info and name from the unit.
	getFirmware                = 0x69 // getFirmware gets firmware info and device name.
	getVentProfiles            = 0xCD // getVentLevels gets the RPM profiles for every speed level.
	getTemps                   = 0xD1 // getTemps gets all available temperatures.
	getBypass                  = 0xDF // getBypass gets heat exchanger information.
	getHours                   = 0xDD // getHours gets the working hours for moving parts.
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
