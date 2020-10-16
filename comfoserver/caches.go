package comfoserver

import (
	"io"
	"log"
	"sync"
	"time"

	"github.com/ti-mo/comfo/libcomfo"
)

var (
	comfoConn io.ReadWriteCloser

	bootInfoCache    BootInfoCache
	tempCache        TempCache
	fanCache         FanCache
	fanProfilesCache FanProfilesCache
	errorsCache      ErrorsCache

	flushCache   = make(chan CacheType)
	flushSuccess = make(chan bool)
)

// BootInfoCache wraps libcomfo's BootInfo structure with caching data.
type BootInfoCache struct {
	libcomfo.BootInfo
	LastUpdated time.Time  `json:"last_updated"`
	CacheLock   sync.Mutex `json:"-"`
}

// TempCache wraps libcomfo's Temps structure with caching data.
type TempCache struct {
	libcomfo.Temps
	LastUpdated time.Time  `json:"last_updated"`
	CacheLock   sync.Mutex `json:"-"`
}

// FanCache wraps libcomfo's Fans structure with caching data.
type FanCache struct {
	libcomfo.Fans
	LastUpdated time.Time  `json:"last_updated"`
	CacheLock   sync.Mutex `json:"-"`
}

// FanProfilesCache wraps libcomfo's FanProfiles structure with caching data.
type FanProfilesCache struct {
	libcomfo.FanProfiles
	LastUpdated time.Time    `json:"last_updated"`
	CacheLock   sync.RWMutex `json:"-"`
}

// ErrorsCache wraps libcomfo's Errors structure with caching data.
type ErrorsCache struct {
	libcomfo.Errors
	LastUpdated time.Time    `json:"last_updated"`
	CacheLock   sync.RWMutex `json:"-"`
}

// UpdateCaches is a macro method to update all data caches of the unit.
// This is run at startup and must populate fields that are not meant to be
// refreshed continuously (eg. BootInfo).
func UpdateCaches(force bool) {
	bootInfoCache.Update(force)
	tempCache.Update(force)
	fanCache.Update(force)
	fanProfilesCache.Update(force)
	errorsCache.Update(force)
}

// FlushCaches manages forced cache flushes for all components and
// sends messages down the flushSuccess channel.
func FlushCaches(cache CacheType) {
	switch cache {
	case BootInfo:
		bootInfoCache.Update(true)
		flushSuccess <- true
	case Fans:
		fanCache.Update(true)
		flushSuccess <- true
	case Temps:
		tempCache.Update(true)
		flushSuccess <- true
	case Profiles:
		fanProfilesCache.Update(true)
		flushSuccess <- true
	case All:
		UpdateCaches(true)
		flushSuccess <- true
	default:
		flushSuccess <- false
	}
}

// Update executes a libcomfo query to fetch boot (firmware) information
// from the unit and sets LastUpdated on the cache object.
// The force parameter ignores the staleness check and updates anyway.
func (bc *BootInfoCache) Update(force bool) {

	// Lock the cache object.
	bc.CacheLock.Lock()
	defer bc.CacheLock.Unlock()

	// Freeze transaction time to start of method.
	now := time.Now()

	// Do not update cache if we're not forced to
	// and if the update is not due yet.
	if !force && !isStale(bc.LastUpdated, now) {
		return
	}

	// Call out to the unit and update object.
	if gb, err := libcomfo.GetBootloader(comfoConn); err == nil {
		bc.BootInfo = gb
		bc.LastUpdated = now
	} else {
		log.Printf("BootInfoCache.Update() - Error updating boot info cache: %s", err)
	}
}

// Update executes a libcomfo query to fetch temperature
// data from the unit and sets LastUpdated on the cache object.
// The force parameter ignores the staleness check and updates anyway.
func (tc *TempCache) Update(force bool) {

	// Lock the cache object
	tc.CacheLock.Lock()
	defer tc.CacheLock.Unlock()

	// Freeze transaction time to start of method
	now := time.Now()

	// Do not update cache if we're not forced to
	// and if the update is not due yet.
	if !force && !isStale(tc.LastUpdated, now) {
		return
	}

	// Call out to the unit and update object
	if gt, err := libcomfo.GetTemperatures(comfoConn); err == nil {
		tc.Temps = gt
		tc.LastUpdated = now
	} else {
		log.Printf("TempCache.Update() - Error updating temperature cache: %s", err)
	}
}

// Update executes a libcomfo query to fetch fan speed
// data from the unit and sets LastUpdated on the cache object.
// The force parameter ignores the staleness check and updates anyway.
func (fc *FanCache) Update(force bool) {

	// Lock the cache object
	fc.CacheLock.Lock()
	defer fc.CacheLock.Unlock()

	// Freeze transaction time to start of method
	now := time.Now()

	// Do not update cache if we're not forced to
	// and if the update is not due yet
	if !force && !isStale(fc.LastUpdated, now) {
		return
	}

	// Call out to the unit and update object
	if gf, err := libcomfo.GetFans(comfoConn); err == nil {
		fc.Fans = gf
		fc.LastUpdated = now
	} else {
		log.Printf("FanCache.Update() - Error updating fan cache: %s", err)
	}
}

// UpdatePoll updates the cache for a maximum amount of `count` tries
// or until the current fan percentage is equal to the desired percentage.
// Takes out a read lock on the profiles cache
func (fc *FanCache) UpdatePoll(count int, tgtSpeed uint8) error {

	// Look up the fan profile matching the target speed
	fanProfilesCache.CacheLock.RLock() // Wait for any changes to the profile
	inPct, outPct, err := fanProfilesCache.Lookup(tgtSpeed)
	if err != nil {
		return err
	}
	fanProfilesCache.CacheLock.RUnlock()

	// Poll updates from the fan cache until desired result
	for i := 0; i < count; i++ {

		fanCache.Update(true)

		// Return when target fan profile is in desired state
		if fanCache.InPercent == inPct &&
			fanCache.OutPercent == outPct {
			return nil
		}

		// On my unit, I'm seeing multiple hundreds of ms of delay
		// of the speed readings on the fan unit, so wait between polls
		time.Sleep(time.Millisecond * 400)
	}

	return errSetPollFailed
}

// Update executes a libcomfo query to fetch fan profiles
// data from the unit and sets LastUpdated on the cache object.
// The force parameter ignores the staleness check and updates anyway.
func (fpc *FanProfilesCache) Update(force bool) {

	// Lock the cache object
	fpc.CacheLock.Lock()
	defer fpc.CacheLock.Unlock()

	// Freeze transaction time to start of method
	now := time.Now()

	// Do not update cache if we're not forced to
	// and if the update is not due yet
	if !force && !isStale(fpc.LastUpdated, now) {
		return
	}

	// Call out to the unit and update object
	if fp, err := libcomfo.GetFanProfiles(comfoConn); err == nil {
		fpc.FanProfiles = fp
		fpc.LastUpdated = now
	} else {
		log.Printf("FanProfilesCache.Update() - Error updating fan profiles cache: %s", err)
	}
}

// Update executes a libcomfo query to fetch error states
// from the unit and sets LastUpdated on the cache object.
// The force parameter ignores the staleness check and updates anyway.
func (ec *ErrorsCache) Update(force bool) {

	// Lock the cache object
	ec.CacheLock.Lock()
	defer ec.CacheLock.Unlock()

	// Freeze transaction time to start of method
	now := time.Now()

	// Do not update cache if we're not forced to
	// and if the update is not due yet
	if !force && !isStale(ec.LastUpdated, now) {
		return
	}

	// Call out to the unit and update object
	if e, err := libcomfo.GetErrors(comfoConn); err == nil {
		ec.Errors = e
		ec.LastUpdated = now
	} else {
		log.Printf("Errors.Update() - Error updating errors cache: %s", err)
	}
}

// isStale applies fuzzy logic to determine if a timestamp is within a certain
// margin of error (10%) from expiring.
func isStale(last time.Time, now time.Time) (stale bool) {

	errorMargin := fastCacheStale / 10

	if now.Sub(last) > fastCacheStale || now.Sub(last) > fastCacheStale-errorMargin {
		return true
	}

	return
}
