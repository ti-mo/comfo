package comfoserver

import (
	"io"
	"log/slog"
	"time"
)

var (
	// 'Fast' cache to be refreshed every 5 seconds.
	fastCacheTimer time.Ticker
	fastCacheStale = time.Second * 5

	// Slower caching tier for data to be polled every 5 minutes.
	slowCacheTimer time.Ticker
	slowCacheStale = time.Minute * 5

	// Data that rarely changes, eg. bootinfo and errors.
	glacialCacheTimer time.Ticker
	glacialCacheStale = time.Hour
)

// StartCaches initializes the data caches,
// starts worker timers and starts the cache worker.
func StartCaches(conn io.ReadWriteCloser) {

	comfoConn = conn

	slog.Info("Populating caches")

	// Pull in first data into caches
	UpdateCaches(true)

	// Initialize both cache timers to their respective staleness thresholds
	fastCacheTimer = *time.NewTicker(fastCacheStale)
	slowCacheTimer = *time.NewTicker(slowCacheStale)
	glacialCacheTimer = *time.NewTicker(glacialCacheStale)

	// Start cache worker
	go CacheWorker()

	slog.Info("Cache worker started")
}

// CacheWorker is responsible for keeping the application
// caches up-to-date and making sure only one operation
// is being executed at a time.
func CacheWorker() {
	for {
		select {
		case <-fastCacheTimer.C:
			tempCache.Update(false)
			fanCache.Update(false)
		case <-slowCacheTimer.C:
			fanProfilesCache.Update(false)
			bypassCache.Update(false)
		case <-glacialCacheTimer.C:
			bootInfoCache.Update(false)
			errorsCache.Update(false)
		case c := <-flushCache:
			FlushCaches(c)
		}
	}
}
