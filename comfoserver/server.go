package comfoserver

import (
	"context"
	"fmt"
	"math"
	"sync"
	"time"

	"github.com/twitchtv/twirp"

	"github.com/ti-mo/comfo/libcomfo"

	rpc "github.com/ti-mo/comfo/rpc/comfo"
)

var (
	fanLock = sync.Mutex{}
)

// Server implements the Comfo RPC server.
type Server struct{}

//
// Getters return information about the unit.
//

// GetErrors returns a unit errors cache protobuf.
func (s *Server) GetErrors(context.Context, *rpc.Noop) (*rpc.Errors, error) {

	return errorsCache.Protobuf(), nil
}

// GetFans returns a fan speed cache protobuf.
func (s *Server) GetFans(context.Context, *rpc.Noop) (*rpc.Fans, error) {

	return fanCache.Protobuf(), nil
}

// GetFanProfiles returns a fan speed profiles cache protobuf.
func (s *Server) GetFanProfiles(context.Context, *rpc.Noop) (*rpc.FanProfiles, error) {

	return fanProfilesCache.Protobuf(), nil
}

// GetTemps returns a temperature cache protobuf.
func (s *Server) GetTemps(context.Context, *rpc.Noop) (*rpc.Temps, error) {

	return tempCache.Protobuf(), nil
}

//
// Setters that modify the state of the unit.
//

// SetComfortTemp updates the comfort temperature on the unit,
// causing it to recover the heat from inside the house up to this point.
func (s *Server) SetComfortTemp(ctx context.Context, ct *rpc.ComfortTarget) (*rpc.ComfortModified, error) {

	start := time.Now()

	var err error
	var modified bool

	origTemp := tempCache.Comfort
	targetTemp := uint8(ct.ComfortTemp)

	// Detect truncation
	if ct.ComfortTemp > math.MaxUint8 {
		return nil, twirp.InvalidArgumentError("ComfortTemp", "is out of range")
	}

	if uint8(origTemp) != targetTemp {

		// Lock the fan speed mutex
		fanLock.Lock()
		defer fanLock.Unlock()

		err = libcomfo.SetComfort(targetTemp, comfoConn)
		if err != nil {
			return nil, twirp.InternalError(err.Error())
		}

		modified = true
	}

	return &rpc.ComfortModified{
		Modified:     modified,
		OriginalTemp: uint32(origTemp),
		TargetTemp:   ct.ComfortTemp,
		ReqTime:      fmt.Sprint(time.Since(start)),
	}, nil
}

// SetFanSpeed updates the fan speed on the unit and updates
// the fan speed and fan profile cache objects.
func (s *Server) SetFanSpeed(ctx context.Context, fst *rpc.FanSpeedTarget) (*rpc.FanSpeedModified, error) {

	start := time.Now()

	var err error
	var modified bool

	// Initialize speed values
	origSpeed := fanProfilesCache.CurrentMode

	// Apply action string to original speed
	tgtSpeed, err := modifySpeed(origSpeed, fst)
	if err != nil {
		return nil, twirp.InternalError(err.Error())
	}

	// Only send actions to the unit if speed needs to be modified
	if uint8(tgtSpeed) != origSpeed {

		// Lock the fan speed mutex
		fanLock.Lock()
		defer fanLock.Unlock()

		// Send target speed to the unit
		err = libcomfo.SetSpeed(uint8(tgtSpeed), comfoConn)
		if err != nil {
			return nil, twirp.InternalError(err.Error())
		}

		// Poll the unit for the target speed
		err = fanCache.UpdatePoll(10, uint8(tgtSpeed))
		if err != nil {
			return nil, twirp.InternalError(err.Error())
		}

		// Flush the fanProfilesCache when values have converged
		fanProfilesCache.Update(true)

		// Set modified flag
		modified = true
	}

	return &rpc.FanSpeedModified{
		Modified:      modified,
		OriginalSpeed: uint32(origSpeed),
		TargetSpeed:   uint32(tgtSpeed),
		ReqTime:       fmt.Sprint(time.Since(start)),
	}, nil
}

// FlushCache synchronously updates all data caches of the unit.
func (s *Server) FlushCache(ctx context.Context, fcr *rpc.FlushCacheRequest) (*rpc.FlushCacheResponse, error) {

	start := time.Now()

	// Feed the 'cache' URI parameter to the flush worker
	flushCache <- fcr.Type

	// Build response
	return &rpc.FlushCacheResponse{
		Success: <-flushSuccess,
		ReqTime: fmt.Sprint(time.Since(start)),
	}, nil
}
