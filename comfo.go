package main

import (
	"fmt"
	"io"
	"log/slog"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/ti-mo/comfo/comfoserver"

	"github.com/spf13/viper"
	"github.com/tarm/serial"
)

var (
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

	// Configure Viper
	viper.SetEnvPrefix("comfo")
	viper.AutomaticEnv()

	// Configure logging.
	level := slog.LevelInfo
	if viper.GetBool(configDebug) {
		level = slog.LevelDebug
	}
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: level}))
	slog.SetDefault(logger)
	slog.Debug("Debug logging enabled")

	// Open connection to unit
	c, err := ConnectUnit(viper.GetString(configMode), viper.GetString(configTarget))
	if err != nil {
		slog.Error("Error connecting to unit", "error", err)
		os.Exit(1)
	}
	defer c.Close()

	// Initialize and start cache timers
	comfoserver.StartCaches(c)

	// Initialize router and listen for connections
	router := NewRouter()

	srv := &http.Server{
		Handler:      router,
		Addr:         viper.GetString(configListen),
		WriteTimeout: 5 * time.Second,
		ReadTimeout:  5 * time.Second,
	}

	slog.Info("API listening", "address", viper.GetString(configListen))

	if err := srv.ListenAndServe(); err != nil {
		slog.Error("ListenAndServe", "error", err)
		os.Exit(1)
	}
}

// ConnectUnit sets up a connection to the unit over TCP or Serial.
func ConnectUnit(mode string, unit string) (conn io.ReadWriteCloser, err error) {
	switch mode {
	case "tcp":
		// Establish TCP connection
		slog.Info("Connecting to the unit over tcp", "address", unit)

		conn, err = net.Dial("tcp", unit)
		if err != nil {
			return nil, fmt.Errorf("unable to dial the unit at %s: %w", unit, err)
		}

		slog.Info("Connection established!", "address", unit)

	case "serial":
		// Establish serial connection
		slog.Info("Opening serial device", "device", unit)

		conn, err = serial.OpenPort(&serial.Config{Name: unit, Baud: 9600})
		if err != nil {
			return nil, fmt.Errorf("unable to open serial device at %s: %w", unit, err)
		}

		slog.Info("Opened device", "device", unit)

	default:
		return nil, fmt.Errorf("unsupported unit mode %s", mode)
	}

	return
}
