package libcomfo

// getRequest is an enum holding the command types
// of all Get operations.
type getRequest uint8

// setRequest is an enum holding the command types
// of all Set operations.
type setRequest uint8

//
// get* requests
// Get information from the unit, no state is altered.
//

const (
	getFans        getRequest = 0x0B // getFans gets ventilator speed % and RPM
	getBootloader  getRequest = 0x67 // getBootloader gets bootloader info and name from the unit.
	getFirmware    getRequest = 0x69 // getFirmware gets firmware info and device name.
	getFanProfiles getRequest = 0xCD // getFanProfiles gets the RPM profiles for every speed level.
	getTemps       getRequest = 0xD1 // getTemps gets all available temperatures.
	getBypass      getRequest = 0xDF // getBypass gets heat exchanger information.
	getHours       getRequest = 0xDD // getHours gets the working hours for moving parts.
)

//
// set* requests
// Alter unit state or configuration
//

const (
	setSpeed   setRequest = 0x99
	setComfort setRequest = 0xD3
)

// SetRequest setSpeedT controls the fan speed of the unit.
// The Speed member can range from 0 (away) to 4 (highest).
type setSpeedT struct {
	Speed uint8
}

func (q setSpeedT) Type() setRequest { return setSpeed }
func (q setSpeedT) MarshalBinary() (out []byte, err error) {

	if q.Speed > 4 {
		return nil, errTooHigh
	}

	return []byte{q.Speed}, nil
}

// SetRequest setComfortT controls the comfort temperature
// of the unit. This is the point at which the heat exchanger
// will recover outgoing heat into the incoming air stream.
type setComfortT struct {
	Temperature temperature
}

func (q setComfortT) Type() setRequest { return setComfort }
func (q setComfortT) MarshalBinary() (out []byte, err error) {

	t, err := q.Temperature.MarshalBinary()
	if err != nil {
		return
	}

	return []byte{t}, nil
}
