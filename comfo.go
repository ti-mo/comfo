package main

import (
	"fmt"
	"github.com/tarm/serial"
	"io"
	"log"
	"net"
	"net/http"
	"time"
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

	// API Server Configuration

	// ListenAddress is the address the API should serve on
	ListenAddress = "0.0.0.0"
	// ListenPort is the port number the API should serve on
	ListenPort = 3094

	// Unit configuration

	// UnitMode is the connection mode of the unit, serial or tcp
	UnitMode = "serial"
	// UnitTCPAddress is the IP address:port of the serial socket forwarder
	UnitTCPAddress = "10.1.1.5:1234"
	// UnitSerialPort is the identifier of the serial port to open.
	UnitSerialPort = "/dev/ttyUSB0"
)

func main() {
	fmt.Printf("Comfo API %v - home automation endpoint for ComfoAir-based ventilation units\n\n", Version)
	fmt.Printf("Git Revision: %v\nBuild time: %v, with %v\n\n", GitRev, BuildTime, GoVersion)

	var err error
	var conn io.ReadWriteCloser

	if UnitMode == "tcp" {
		// Establish TCP connection
		log.Printf("Connecting to the unit over tcp at %v ..", UnitTCPAddress)
		conn, err = net.Dial("tcp", UnitTCPAddress)
		if err != nil {
			log.Fatalf("unable to dial the unit at %v: %v", UnitTCPAddress, err)
		}
		log.Printf("Connection to %v established!", UnitTCPAddress)
	} else if UnitMode == "serial" {
		log.Printf("Opening serial device %v ..", UnitSerialPort)
		conn, err = serial.OpenPort(&serial.Config{Name: UnitSerialPort, Baud: 9600})
		if err != nil {
			log.Fatalf("unable to open serial device at %v: %v", UnitSerialPort, err)
		}
		log.Printf("Opened device %v!", UnitSerialPort)
	} else {
		log.Fatalf("unsupported unit mode %v", UnitMode)
	}

	// Set connection to package-wide variable
	comfo = conn
	defer comfo.Close()

	// Initialize and start cache timers
	StartCaches()

	// Initialize router and listen for connections
	router := NewRouter()

	srv := &http.Server{
		Handler:      router,
		Addr:         fmt.Sprintf("%s:%d", ListenAddress, ListenPort),
		WriteTimeout: 5 * time.Second,
		ReadTimeout:  5 * time.Second,
	}

	// Print API endpoint info
	printEndpoints(ListenPort)

	log.Fatal(srv.ListenAndServe())
}
