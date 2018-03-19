package main

import (
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var (
	// Configuration Keys

	ConfigMode   = "mode"
	ConfigTarget = "target"
	ConfigListen = "listen"

	// Default Configuration Values

	DefaultMode   = "serial"
	DefaultTarget = "/dev/ttyUSB0"
	DefaultListen = "[::]:3094"
)

func init() {

	// Register Flags
	pflag.String(ConfigMode, DefaultMode, "The mode to connect to the unit (serial/tcp)")
	pflag.String(ConfigTarget, DefaultTarget, "The address or serial device of the unit")
	pflag.String(ConfigListen, DefaultListen, "Address to bind the API on.")

	pflag.Parse()
	viper.BindPFlags(pflag.CommandLine)

}
