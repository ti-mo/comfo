package libcomfo

import (
	"io"
	"sync"
)

var (
	// API mutex for calls that perform a potential dirty r/w
	am = &sync.RWMutex{}
)

// SetComfort sets the comfort temperature on the unit.
func SetComfort(val uint8, conn io.ReadWriter) (err error) {

	return setQuery(setComfortT{Temperature: temperature(val)}, conn)
}

// SetSpeed sets the fan speed of the unit.
func SetSpeed(val uint8, conn io.ReadWriter) (err error) {

	return setQuery(setSpeedT{Speed: val}, conn)
}

// SetFanProfile sets a fan percentage of a single speed level.
func SetFanProfile(level uint8, speed uint8, conn io.ReadWriter) error {

	// Make sure speed level is in range
	if level > 3 {
		return errUnknownLevel
	}

	// Don't accept any values over 100 (percent)
	if speed > 100 {
		return errTooHigh
	}

	// Lock API mutex for r/w transaction
	am.Lock()
	defer am.Unlock()

	cfp, err := GetFanProfiles(conn)
	if err != nil {
		return err
	}

	switch level {
	case 0:
		cfp.InAway, cfp.OutAway = speed, speed
	case 1:
		cfp.InLow, cfp.OutLow = speed, speed
	case 2:
		cfp.InMid, cfp.OutMid = speed, speed
	case 3:
		cfp.InHigh, cfp.OutHigh = speed, speed
	}

	return setQuery(cfp, conn)
}

// SetFanProfiles sets the fan percentages associated with every speed level.
// Only away/low/mid/high's in and out fields are sent to the unit.
func SetFanProfiles(fp *FanProfiles, conn io.ReadWriter) error {

	return setQuery(fp, conn)
}

// GetTemperatures gets the temperature readings from the unit.
func GetTemperatures(conn io.ReadWriter) (temps Temps, err error) {

	resp, err := getQuery(getTemps, conn)
	if err != nil {
		return
	}

	return *resp.(*Temps), err
}

// GetBypass gets the bypass information from the unit.
func GetBypass(conn io.ReadWriter) (bypass Bypass, err error) {

	resp, err := getQuery(getBypass, conn)
	if err != nil {
		return
	}

	return *resp.(*Bypass), err
}

// GetFans gets the speeds of the unit's ventilators.
func GetFans(conn io.ReadWriter) (fans Fans, err error) {

	resp, err := getQuery(getFans, conn)
	if err != nil {
		return
	}

	return *resp.(*Fans), err
}

// GetHours gets the operating hours for all moving parts in the unit.
func GetHours(conn io.ReadWriter) (hours Hours, err error) {

	resp, err := getQuery(getHours, conn)
	if err != nil {
		return
	}

	return *resp.(*Hours), err
}

// GetBootloader gets the bootloader information from the unit.
func GetBootloader(conn io.ReadWriter) (bi BootInfo, err error) {

	resp, err := getQuery(getBootloader, conn)
	if err != nil {
		return
	}

	return *resp.(*BootInfo), err
}

// GetFirmware gets the firmware information from the unit.
func GetFirmware(conn io.ReadWriter) (bi BootInfo, err error) {

	resp, err := getQuery(getFirmware, conn)
	if err != nil {
		return
	}

	return *resp.(*BootInfo), err
}

// GetFanProfiles gets the fan profiles for each ventilation level.
func GetFanProfiles(conn io.ReadWriter) (fp FanProfiles, err error) {

	resp, err := getQuery(getFanProfiles, conn)
	if err != nil {
		return
	}

	return *resp.(*FanProfiles), err
}

// GetErrors gets the unit's error statuses and their values.
func GetErrors(conn io.ReadWriter) (e Errors, err error) {

	resp, err := getQuery(getErrors, conn)
	if err != nil {
		return
	}

	return *resp.(*Errors), err
}
