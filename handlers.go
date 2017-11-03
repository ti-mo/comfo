package main

import (
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/ti-mo/comfo/libcomfo"
	"net/http"
	"sync"
	"time"
)

var (
	fanLock = sync.Mutex{}
)

// TempHandlerGet displays the temperature cache object.
func TempHandlerGet(w http.ResponseWriter, r *http.Request) {

	jsonWrite(w, &tempCache)
}

// FanHandlerGet displays the fan speed cache object.
func FanHandlerGet(w http.ResponseWriter, r *http.Request) {

	jsonWrite(w, &fanCache)
}

// FanHandlerSet updates the fan speed on the unit and updates
// the fan speed and fan profile cache objects.
func FanHandlerSet(w http.ResponseWriter, r *http.Request) {

	start := time.Now()
	vars := mux.Vars(r)

	var err error
	var modified bool

	// Initialize speed values
	origSpeed := fanProfilesCache.CurrentLevel

	// Input validation
	actionStr, ok := vars["speed"]
	if !ok {
		jsonError(w, errors.New("missing request variable 'speed'"))
		return
	}

	// Apply action string to original speed
	tgtSpeed, err := modifySpeed(int(origSpeed), actionStr)
	if err != nil {
		jsonError(w, err)
		return
	}

	// Only send actions to the unit if speed needs to be modified
	if uint8(tgtSpeed) != origSpeed {

		// Lock the fan speed mutex
		fanLock.Lock()
		defer fanLock.Unlock()

		// Send target speed to the unit
		err = libcomfo.SetSpeed(uint8(tgtSpeed), comfo)
		if err != nil {
			jsonError(w, err)
			return
		}

		// Poll the unit for the target speed
		err = fanCache.UpdatePoll(10, uint8(tgtSpeed))
		if err != nil {
			jsonError(w, err)
			return
		}

		// Flush the fanProfilesCache when values have converged
		fanProfilesCache.Update(true)

		// Set modified flag
		modified = true
	}

	// Compose response
	resp := map[string]interface{}{
		"success":        true,
		"modified":       modified,
		"original_speed": origSpeed,
		"target_speed":   tgtSpeed,
		"fans":           &fanCache,
		"reqtime":        fmt.Sprint(time.Since(start)),
	}

	jsonWrite(w, &resp)
}

// FanProfilesHandlerGet displays the fan speed profiles cache object.
func FanProfilesHandlerGet(w http.ResponseWriter, r *http.Request) {

	jsonWrite(w, &fanProfilesCache)
}

// FlushCacheHandler synchronously updates all data caches of the unit.
func FlushCacheHandler(w http.ResponseWriter, r *http.Request) {

	start := time.Now()
	vars := mux.Vars(r)

	// Feed the 'cache' URI parameter to the flush worker
	flushCache <- vars["cache"]

	// Build response
	resp := map[string]interface{}{
		"success": <-flushSuccess,
		"reqtime": fmt.Sprint(time.Since(start)),
	}

	// Write response
	jsonWrite(w, resp)
}
