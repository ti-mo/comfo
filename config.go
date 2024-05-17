package main

import (
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var (
	// Configuration Keys
	configDebug  = "debug"
	configMode   = "mode"
	configTarget = "target"
	configListen = "listen"

	// Default Configuration Values
	defaultDebug  = false
	defaultMode   = "serial"
	defaultTarget = "/dev/ttyUSB0"
	defaultListen = "[::]:3094"
)

func init() {
	// Register Flags
	pflag.Bool(configDebug, defaultDebug, "Whether to enable debug logging")
	pflag.String(configMode, defaultMode, "The mode to connect to the unit (serial/tcp)")
	pflag.String(configTarget, defaultTarget, "The address or serial device of the unit")
	pflag.String(configListen, defaultListen, "Address to bind the API on.")

	pflag.Parse()
	if err := viper.BindPFlags(pflag.CommandLine); err != nil {
		panic(err)
	}
}
