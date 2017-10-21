package libcomfo

import (
	"io"
)

// SetComfort sets the comfort temperature on the unit.
func SetComfort(val uint8, conn io.ReadWriter) (err error) {

	return setQuery(setComfort{Temperature: temperature(val)}, conn)
}

// SetSpeed sets the fan speed of the unit.
func SetSpeed(val uint8, conn io.ReadWriter) (err error) {

	return setQuery(setSpeed{Speed: val}, conn)
}

// GetTemperatures gets the temperature readings from the unit.
func GetTemperatures(conn io.ReadWriter) (temps Temps, err error) {

	resp, err := getQuery(getTemps, conn)
	if err != nil {
		return Temps{}, err
	}

	return *resp.(*Temps), err
}

// GetBypass gets the bypass information from the unit.
func GetBypass(conn io.ReadWriter) (bypass Bypass, err error) {

	resp, err := getQuery(getBypass, conn)
	if err != nil {
		return Bypass{}, err
	}

	return *resp.(*Bypass), err
}

// GetHours gets the bypass information from the unit.
func GetHours(conn io.ReadWriter) (hours Hours, err error) {

	resp, err := getQuery(getHours, conn)
	if err != nil {
		return Hours{}, err
	}

	return *resp.(*Hours), err
}

func GetBootloader(conn io.ReadWriter) (bi BootInfo, err error) {
	resp, err := getQuery(getBootloader, conn)
	if err != nil {
		return BootInfo{}, err
	}

	return *resp.(*BootInfo), err
}

func GetFirmware(conn io.ReadWriter) (bi BootInfo, err error) {

	resp, err := getQuery(getFirmware, conn)
	if err != nil {
		return BootInfo{}, err
	}

	return *resp.(*BootInfo), err
}
