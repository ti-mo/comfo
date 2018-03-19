package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/spf13/viper"
	"github.com/tarm/serial"
)

var (
	comfo io.ReadWriteCloser

	// Compile-time injected variables

	// Version is the Comfo API version
	Version string
	// GitRev is the git revision the binary was built with
	GitRev string
	// BuildTime is the binary build timestamp
	BuildTime string
	// GoVersion is the go compiler version the binary was built with
	GoVersion string
)

func main() {
	fmt.Printf("Comfo API %v - home automation endpoint for ComfoAir-based ventilation units\n\n", Version)
	fmt.Printf("Git Revision: %v\nBuild time: %v, with %v\n\n", GitRev, BuildTime, GoVersion)

	// Open connection to unit
	c, err := ConnectUnit(viper.GetString(ConfigMode), viper.GetString(ConfigTarget))
	if err != nil {
		log.Fatalln("Error connecting to unit:", err)
	}

	// Save connection reference to package
	comfo = c
	defer comfo.Close()

	// Initialize and start cache timers
	StartCaches()

	// Initialize router and listen for connections
	router := NewRouter()

	srv := &http.Server{
		Handler:      router,
		Addr:         viper.GetString(ConfigListen),
		WriteTimeout: 5 * time.Second,
		ReadTimeout:  5 * time.Second,
	}

	log.Println("API listening on", viper.GetString(ConfigListen))

	log.Fatal(srv.ListenAndServe())
}

// ConnectUnit sets up a connection to the unit over TCP or Serial.
func ConnectUnit(mode string, unit string) (conn io.ReadWriteCloser, err error) {

	switch mode {
	case "tcp":
		// Establish TCP connection
		log.Printf("Connecting to the unit over tcp at %v ..", unit)

		conn, err = net.Dial("tcp", unit)
		if err != nil {
			log.Fatalf("unable to dial the unit at %v: %v", unit, err)
		}

		log.Printf("Connection to %v established!", unit)

	case "serial":
		// Establish serial connection
		log.Printf("Opening serial device %v ..", unit)

		conn, err = serial.OpenPort(&serial.Config{Name: unit, Baud: 9600})
		if err != nil {
			log.Fatalf("unable to open serial device at %v: %v", unit, err)
		}

		log.Printf("Opened device %v!", unit)

	default:
		log.Fatalf("unsupported unit mode %v", mode)
	}

	return
}
