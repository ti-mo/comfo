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
	getFanProfiles getRequest = 0xCD // getFanProfiles gets the RPM profiles for every speed mode.
	getTemps       getRequest = 0xD1 // getTemps gets all available temperatures.
	getErrors      getRequest = 0xD9 // getErrors get a list of error slots and their values.
	getBypass      getRequest = 0xDF // getBypass gets heat exchanger information.
	getHours       getRequest = 0xDD // getHours gets the working hours for moving parts.
)

//
// set* requests
// Alter unit state or configuration
//

const (
	setSpeed       setRequest = 0x99
	setComfort     setRequest = 0xD3
	setFanProfiles setRequest = 0xCF
)

// SetRequest setSpeedT controls the fan speed of the unit.
// The Speed member can range from 0 (away) to 4 (highest).
type setSpeedT struct {
	Speed uint8
}

func (q setSpeedT) requestType() setRequest { return setSpeed }
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

func (q setComfortT) requestType() setRequest { return setComfort }
func (q setComfortT) MarshalBinary() (out []byte, err error) {

	t, err := q.Temperature.MarshalBinary()
	if err != nil {
		return
	}

	return []byte{t}, nil
}

// Type returns the FanProfiles set request message code.
func (q FanProfiles) requestType() setRequest { return setFanProfiles }

// MarshalBinary marshals a FanProfiles into a byte representation
// to be used as a setRequest.
func (q FanProfiles) MarshalBinary() (out []byte, err error) {

	out = make([]byte, 8)

	out[0] = q.OutAway
	out[1] = q.OutLow
	out[2] = q.OutMid
	out[3] = q.InAway
	out[4] = q.InLow
	out[5] = q.InMid
	out[6] = q.OutHigh
	out[7] = q.InHigh

	return
}
