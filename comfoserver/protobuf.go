package comfoserver

import (
	rpc "github.com/ti-mo/comfo/rpc/comfo"
)

type CacheType = rpc.FlushCacheRequest_CacheType

const (
	BootInfo = rpc.FlushCacheRequest_BootInfo
	Bypass   = rpc.FlushCacheRequest_Bypass
	Fans     = rpc.FlushCacheRequest_Fans
	Temps    = rpc.FlushCacheRequest_Temps
	Profiles = rpc.FlushCacheRequest_Profiles
	All      = rpc.FlushCacheRequest_All
)

// Protobuf takes out a lock on the BootInfoCache
// and returns a protobuf representation of BootInfo.
func (b *BootInfoCache) Protobuf() *rpc.BootInfo {

	b.CacheLock.Lock()
	defer b.CacheLock.Unlock()

	return &rpc.BootInfo{
		BetaVersion:  uint32(b.BetaVersion),
		MajorVersion: uint32(b.MajorVersion),
		MinorVersion: uint32(b.MinorVersion),
		DeviceName:   b.DeviceName,
	}
}

// Protobuf takes out a lock on the FanCache and
// returns a protobuf representation of Bypass.
func (b *BypassCache) Protobuf() (pb *rpc.Bypass) {

	b.CacheLock.Lock()
	defer b.CacheLock.Unlock()

	return &rpc.Bypass{
		Correction: uint32(b.Correction),
		Factor:     uint32(b.Factor),
		Level:      uint32(b.Level),
		SummerMode: bool(b.SummerMode),
	}
}

// Protobuf takes out a lock on the TempCache and
// returns a protobuf representation of Temps.
func (t *TempCache) Protobuf() (pb *rpc.Temps) {

	t.CacheLock.Lock()
	defer t.CacheLock.Unlock()

	return &rpc.Temps{
		Comfort:     float32(t.Comfort),
		OutsideAir:  float32(t.OutsideAir),
		SupplyAir:   float32(t.SupplyAir),
		InsideAir:   float32(t.InsideAir),
		ExhaustAir:  float32(t.ExhaustAir),
		GeoHeat:     float32(t.GeoHeat),
		Reheating:   float32(t.Reheating),
		KitchenHood: float32(t.KitchenHood),
	}
}

// Protobuf takes out a lock on the FanCache and
// returns a protobuf representation of Fans.
func (f *FanCache) Protobuf() (pb *rpc.Fans) {

	f.CacheLock.Lock()
	defer f.CacheLock.Unlock()

	return &rpc.Fans{
		InPercent:  uint32(f.InPercent),
		OutPercent: uint32(f.OutPercent),
		InSpeed:    uint32(f.InSpeed),
		OutSpeed:   uint32(f.OutSpeed),
	}
}

// Protobuf takes out a lock on the FanProfilesCache and
// returns a protobuf representation of FanProfiles.
func (fp *FanProfilesCache) Protobuf() (pb *rpc.FanProfiles) {

	fp.CacheLock.Lock()
	defer fp.CacheLock.Unlock()

	return &rpc.FanProfiles{
		OutAway: uint32(fp.OutAway),
		OutLow:  uint32(fp.OutLow),
		OutMid:  uint32(fp.OutMid),
		OutHigh: uint32(fp.OutHigh),

		InFanActive: fp.InFanActive,
		InAway:      uint32(fp.InAway),
		InLow:       uint32(fp.InLow),
		InMid:       uint32(fp.InMid),
		InHigh:      uint32(fp.InHigh),

		CurrentOut:  uint32(fp.CurrentOut),
		CurrentIn:   uint32(fp.CurrentIn),
		CurrentMode: uint32(fp.CurrentMode),
	}
}

// Protobuf takes out a lock on the ErrorsCache and
// returns a protobuf representation of Errors.
func (e *ErrorsCache) Protobuf() (pb *rpc.Errors) {

	e.CacheLock.Lock()
	defer e.CacheLock.Unlock()

	return &rpc.Errors{
		Filter: e.Filter,
	}
}
