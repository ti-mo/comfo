package main

import (
	"io"
)

// SetComfort sets the comfort temperature on the unit.
func SetComfort(val uint8, conn io.ReadWriter) (err error) {

	return setQuery(setComfort{Temperature: temperature(val)}, conn)
}

// GetTemperatures gets the temperature readings from the unit.
func GetTemperatures(conn io.ReadWriter) (temps Temps, err error) {

	resp, err := getQuery(getTemps, conn)

	return *resp.(*Temps), err
}

// GetBypass gets the bypass information from the unit.
func GetBypass(conn io.ReadWriter) (bypass Bypass, err error) {

	resp, err := getQuery(getBypass, conn)

	return *resp.(*Bypass), err
}
