package comfoserver

import (
	rpc "github.com/ti-mo/comfo/rpc/comfo"
)

// Protobuf returns a protobuf representation of TempCache
func (t *TempCache) Protobuf() (pb *rpc.Temps) {

	return &rpc.Temps{
		Comfort:     float32(t.Comfort),
		OutsideAir:  float32(t.OutsideAir),
		SupplyAir:   float32(t.SupplyAir),
		OutAir:      float32(t.OutAir),
		ExhaustAir:  float32(t.ExhaustAir),
		GeoHeat:     float32(t.GeoHeat),
		Reheating:   float32(t.Reheating),
		KitchenHood: float32(t.KitchenHood),
	}
}

// Protobuf returns a protobuf representation of FanCache.
func (f *FanCache) Protobuf() (pb *rpc.Fans) {

	return &rpc.Fans{
		InPercent:  uint32(f.InPercent),
		OutPercent: uint32(f.OutPercent),
		InSpeed:    uint32(f.InSpeed),
		OutSpeed:   uint32(f.OutSpeed),
	}
}

// Protobuf returns a protobuf representation of FanProfilesCache.
func (fp *FanProfilesCache) Protobuf() (pb *rpc.FanProfiles) {

	return &rpc.FanProfiles{
		OutAway: uint32(fanProfilesCache.OutAway),
		OutLow:  uint32(fanProfilesCache.OutLow),
		OutMid:  uint32(fanProfilesCache.OutMid),
		OutHigh: uint32(fanProfilesCache.OutHigh),

		InFanActive: fanProfilesCache.InFanActive,
		InAway:      uint32(fanProfilesCache.InAway),
		InLow:       uint32(fanProfilesCache.InLow),
		InMid:       uint32(fanProfilesCache.InMid),
		InHigh:      uint32(fanProfilesCache.InHigh),

		CurrentOut:   uint32(fanProfilesCache.CurrentOut),
		CurrentIn:    uint32(fanProfilesCache.CurrentIn),
		CurrentLevel: uint32(fanProfilesCache.CurrentLevel),
	}
}