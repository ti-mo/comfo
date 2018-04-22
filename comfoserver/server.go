package comfoserver

import (
	"context"
	"fmt"
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

// GetTemps displays the temperature cache object.
func (s *Server) GetTemps(context.Context, *rpc.Noop) (*rpc.Temps, error) {

	return tempCache.Protobuf(), nil
}

// GetFans returns the fan speed cache object in protobuf format.
func (s *Server) GetFans(context.Context, *rpc.Noop) (*rpc.Fans, error) {

	return fanCache.Protobuf(), nil
}

// SetFanSpeed updates the fan speed on the unit and updates
// the fan speed and fan profile cache objects.
func (s *Server) SetFanSpeed(ctx context.Context, fst *rpc.FanSpeedTarget) (*rpc.FanSpeedModified, error) {

	start := time.Now()

	var err error
	var modified bool

	// Initialize speed values
	origSpeed := fanProfilesCache.CurrentLevel

	// Apply action string to original speed
	tgtSpeed, err := modifySpeed(origSpeed, fst)
	if err != nil {
		return nil, err
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

// GetFanProfiles displays the fan speed profiles cache object.
func (s *Server) GetFanProfiles(context.Context, *rpc.Noop) (*rpc.FanProfiles, error) {

	return fanProfilesCache.Protobuf(), nil
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
