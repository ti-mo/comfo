package main

import (
	"log"
	"time"
)

var (
	fastCacheTimer time.Ticker
	fastCacheStale = time.Second * 2

	slowCacheTimer time.Ticker
	slowCacheStale = time.Minute * 5
)

// StartCaches initializes the data caches,
// starts worker timers and starts the cache worker.
func StartCaches() {

	log.Println("Updating initial caches.")

	// Pull in first data into caches
	UpdateCaches(true)

	// Initialize both cache timers to their respective staleness thresholds
	fastCacheTimer = *time.NewTicker(fastCacheStale)
	slowCacheTimer = *time.NewTicker(slowCacheStale)

	// Start cache worker
	go CacheWorker()

	log.Println("Started cache worker.")
}

// CacheWorker is responsible for keeping the application
// caches up-to-date and making sure only one operation
// is being executed at a time.
func CacheWorker() {
	for {
		select {
		case <-fastCacheTimer.C:
			// Update short-lived caches
			tempCache.Update(false)
			fanCache.Update(false)
		case <-slowCacheTimer.C:
			// Update long-lived caches
			fanProfilesCache.Update(false)
		case c := <-flushCache:
			FlushCaches(c)
		}
	}
}
