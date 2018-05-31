// Code generated by protoc-gen-go. DO NOT EDIT.
// source: github.com/ti-mo/comfo/rpc/comfo/service.proto

/*
Package comfo is a generated protocol buffer package.

It is generated from these files:
	github.com/ti-mo/comfo/rpc/comfo/service.proto

It has these top-level messages:
	Noop
	Bypass
	BootInfo
	ComfortTarget
	ComfortModified
	Errors
	Fans
	FanProfiles
	FanProfileTarget
	FanProfileModified
	FanSpeedTarget
	FanSpeedModified
	Hours
	Temps
	FlushCacheRequest
	FlushCacheResponse
*/
package comfo

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type Noop struct {
}

func (m *Noop) Reset()                    { *m = Noop{} }
func (m *Noop) String() string            { return proto.CompactTextString(m) }
func (*Noop) ProtoMessage()               {}
func (*Noop) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

type Bypass struct {
	Factor     uint32 `protobuf:"varint,1,opt,name=Factor" json:"Factor,omitempty"`
	Level      uint32 `protobuf:"varint,2,opt,name=Level" json:"Level,omitempty"`
	Correction uint32 `protobuf:"varint,3,opt,name=Correction" json:"Correction,omitempty"`
	SummerMode bool   `protobuf:"varint,4,opt,name=SummerMode" json:"SummerMode,omitempty"`
}

func (m *Bypass) Reset()                    { *m = Bypass{} }
func (m *Bypass) String() string            { return proto.CompactTextString(m) }
func (*Bypass) ProtoMessage()               {}
func (*Bypass) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *Bypass) GetFactor() uint32 {
	if m != nil {
		return m.Factor
	}
	return 0
}

func (m *Bypass) GetLevel() uint32 {
	if m != nil {
		return m.Level
	}
	return 0
}

func (m *Bypass) GetCorrection() uint32 {
	if m != nil {
		return m.Correction
	}
	return 0
}

func (m *Bypass) GetSummerMode() bool {
	if m != nil {
		return m.SummerMode
	}
	return false
}

type BootInfo struct {
	MajorVersion uint32 `protobuf:"varint,1,opt,name=MajorVersion" json:"MajorVersion,omitempty"`
	MinorVersion uint32 `protobuf:"varint,2,opt,name=MinorVersion" json:"MinorVersion,omitempty"`
	BetaVersion  uint32 `protobuf:"varint,3,opt,name=BetaVersion" json:"BetaVersion,omitempty"`
	DeviceName   string `protobuf:"bytes,4,opt,name=DeviceName" json:"DeviceName,omitempty"`
}

func (m *BootInfo) Reset()                    { *m = BootInfo{} }
func (m *BootInfo) String() string            { return proto.CompactTextString(m) }
func (*BootInfo) ProtoMessage()               {}
func (*BootInfo) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func (m *BootInfo) GetMajorVersion() uint32 {
	if m != nil {
		return m.MajorVersion
	}
	return 0
}

func (m *BootInfo) GetMinorVersion() uint32 {
	if m != nil {
		return m.MinorVersion
	}
	return 0
}

func (m *BootInfo) GetBetaVersion() uint32 {
	if m != nil {
		return m.BetaVersion
	}
	return 0
}

func (m *BootInfo) GetDeviceName() string {
	if m != nil {
		return m.DeviceName
	}
	return ""
}

type ComfortTarget struct {
	ComfortTemp uint32 `protobuf:"varint,1,opt,name=ComfortTemp" json:"ComfortTemp,omitempty"`
}

func (m *ComfortTarget) Reset()                    { *m = ComfortTarget{} }
func (m *ComfortTarget) String() string            { return proto.CompactTextString(m) }
func (*ComfortTarget) ProtoMessage()               {}
func (*ComfortTarget) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

func (m *ComfortTarget) GetComfortTemp() uint32 {
	if m != nil {
		return m.ComfortTemp
	}
	return 0
}

type ComfortModified struct {
	Modified     bool   `protobuf:"varint,1,opt,name=Modified" json:"Modified,omitempty"`
	OriginalTemp uint32 `protobuf:"varint,2,opt,name=OriginalTemp" json:"OriginalTemp,omitempty"`
	TargetTemp   uint32 `protobuf:"varint,3,opt,name=TargetTemp" json:"TargetTemp,omitempty"`
	ReqTime      string `protobuf:"bytes,4,opt,name=ReqTime" json:"ReqTime,omitempty"`
}

func (m *ComfortModified) Reset()                    { *m = ComfortModified{} }
func (m *ComfortModified) String() string            { return proto.CompactTextString(m) }
func (*ComfortModified) ProtoMessage()               {}
func (*ComfortModified) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{4} }

func (m *ComfortModified) GetModified() bool {
	if m != nil {
		return m.Modified
	}
	return false
}

func (m *ComfortModified) GetOriginalTemp() uint32 {
	if m != nil {
		return m.OriginalTemp
	}
	return 0
}

func (m *ComfortModified) GetTargetTemp() uint32 {
	if m != nil {
		return m.TargetTemp
	}
	return 0
}

func (m *ComfortModified) GetReqTime() string {
	if m != nil {
		return m.ReqTime
	}
	return ""
}

type Errors struct {
	Filter bool `protobuf:"varint,1,opt,name=Filter" json:"Filter,omitempty"`
}

func (m *Errors) Reset()                    { *m = Errors{} }
func (m *Errors) String() string            { return proto.CompactTextString(m) }
func (*Errors) ProtoMessage()               {}
func (*Errors) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{5} }

func (m *Errors) GetFilter() bool {
	if m != nil {
		return m.Filter
	}
	return false
}

type Fans struct {
	InPercent  uint32 `protobuf:"varint,1,opt,name=InPercent" json:"InPercent,omitempty"`
	OutPercent uint32 `protobuf:"varint,2,opt,name=OutPercent" json:"OutPercent,omitempty"`
	InSpeed    uint32 `protobuf:"varint,3,opt,name=InSpeed" json:"InSpeed,omitempty"`
	OutSpeed   uint32 `protobuf:"varint,4,opt,name=OutSpeed" json:"OutSpeed,omitempty"`
}

func (m *Fans) Reset()                    { *m = Fans{} }
func (m *Fans) String() string            { return proto.CompactTextString(m) }
func (*Fans) ProtoMessage()               {}
func (*Fans) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{6} }

func (m *Fans) GetInPercent() uint32 {
	if m != nil {
		return m.InPercent
	}
	return 0
}

func (m *Fans) GetOutPercent() uint32 {
	if m != nil {
		return m.OutPercent
	}
	return 0
}

func (m *Fans) GetInSpeed() uint32 {
	if m != nil {
		return m.InSpeed
	}
	return 0
}

func (m *Fans) GetOutSpeed() uint32 {
	if m != nil {
		return m.OutSpeed
	}
	return 0
}

type FanProfiles struct {
	OutAway      uint32 `protobuf:"varint,1,opt,name=OutAway" json:"OutAway,omitempty"`
	OutLow       uint32 `protobuf:"varint,2,opt,name=OutLow" json:"OutLow,omitempty"`
	OutMid       uint32 `protobuf:"varint,3,opt,name=OutMid" json:"OutMid,omitempty"`
	OutHigh      uint32 `protobuf:"varint,4,opt,name=OutHigh" json:"OutHigh,omitempty"`
	InFanActive  bool   `protobuf:"varint,5,opt,name=InFanActive" json:"InFanActive,omitempty"`
	InAway       uint32 `protobuf:"varint,6,opt,name=InAway" json:"InAway,omitempty"`
	InLow        uint32 `protobuf:"varint,7,opt,name=InLow" json:"InLow,omitempty"`
	InMid        uint32 `protobuf:"varint,8,opt,name=InMid" json:"InMid,omitempty"`
	InHigh       uint32 `protobuf:"varint,9,opt,name=InHigh" json:"InHigh,omitempty"`
	CurrentOut   uint32 `protobuf:"varint,10,opt,name=CurrentOut" json:"CurrentOut,omitempty"`
	CurrentIn    uint32 `protobuf:"varint,11,opt,name=CurrentIn" json:"CurrentIn,omitempty"`
	CurrentLevel uint32 `protobuf:"varint,12,opt,name=CurrentLevel" json:"CurrentLevel,omitempty"`
}

func (m *FanProfiles) Reset()                    { *m = FanProfiles{} }
func (m *FanProfiles) String() string            { return proto.CompactTextString(m) }
func (*FanProfiles) ProtoMessage()               {}
func (*FanProfiles) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{7} }

func (m *FanProfiles) GetOutAway() uint32 {
	if m != nil {
		return m.OutAway
	}
	return 0
}

func (m *FanProfiles) GetOutLow() uint32 {
	if m != nil {
		return m.OutLow
	}
	return 0
}

func (m *FanProfiles) GetOutMid() uint32 {
	if m != nil {
		return m.OutMid
	}
	return 0
}

func (m *FanProfiles) GetOutHigh() uint32 {
	if m != nil {
		return m.OutHigh
	}
	return 0
}

func (m *FanProfiles) GetInFanActive() bool {
	if m != nil {
		return m.InFanActive
	}
	return false
}

func (m *FanProfiles) GetInAway() uint32 {
	if m != nil {
		return m.InAway
	}
	return 0
}

func (m *FanProfiles) GetInLow() uint32 {
	if m != nil {
		return m.InLow
	}
	return 0
}

func (m *FanProfiles) GetInMid() uint32 {
	if m != nil {
		return m.InMid
	}
	return 0
}

func (m *FanProfiles) GetInHigh() uint32 {
	if m != nil {
		return m.InHigh
	}
	return 0
}

func (m *FanProfiles) GetCurrentOut() uint32 {
	if m != nil {
		return m.CurrentOut
	}
	return 0
}

func (m *FanProfiles) GetCurrentIn() uint32 {
	if m != nil {
		return m.CurrentIn
	}
	return 0
}

func (m *FanProfiles) GetCurrentLevel() uint32 {
	if m != nil {
		return m.CurrentLevel
	}
	return 0
}

type FanProfileTarget struct {
	Level       uint32 `protobuf:"varint,1,opt,name=Level" json:"Level,omitempty"`
	TargetSpeed uint32 `protobuf:"varint,2,opt,name=TargetSpeed" json:"TargetSpeed,omitempty"`
}

func (m *FanProfileTarget) Reset()                    { *m = FanProfileTarget{} }
func (m *FanProfileTarget) String() string            { return proto.CompactTextString(m) }
func (*FanProfileTarget) ProtoMessage()               {}
func (*FanProfileTarget) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{8} }

func (m *FanProfileTarget) GetLevel() uint32 {
	if m != nil {
		return m.Level
	}
	return 0
}

func (m *FanProfileTarget) GetTargetSpeed() uint32 {
	if m != nil {
		return m.TargetSpeed
	}
	return 0
}

type FanProfileModified struct {
	Modified      bool   `protobuf:"varint,1,opt,name=Modified" json:"Modified,omitempty"`
	OriginalSpeed uint32 `protobuf:"varint,2,opt,name=OriginalSpeed" json:"OriginalSpeed,omitempty"`
	TargetSpeed   uint32 `protobuf:"varint,3,opt,name=TargetSpeed" json:"TargetSpeed,omitempty"`
	ReqTime       string `protobuf:"bytes,4,opt,name=ReqTime" json:"ReqTime,omitempty"`
}

func (m *FanProfileModified) Reset()                    { *m = FanProfileModified{} }
func (m *FanProfileModified) String() string            { return proto.CompactTextString(m) }
func (*FanProfileModified) ProtoMessage()               {}
func (*FanProfileModified) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{9} }

func (m *FanProfileModified) GetModified() bool {
	if m != nil {
		return m.Modified
	}
	return false
}

func (m *FanProfileModified) GetOriginalSpeed() uint32 {
	if m != nil {
		return m.OriginalSpeed
	}
	return 0
}

func (m *FanProfileModified) GetTargetSpeed() uint32 {
	if m != nil {
		return m.TargetSpeed
	}
	return 0
}

func (m *FanProfileModified) GetReqTime() string {
	if m != nil {
		return m.ReqTime
	}
	return ""
}

// Abs and Rel are mutually exclusive
type FanSpeedTarget struct {
	Abs uint32 `protobuf:"varint,1,opt,name=Abs" json:"Abs,omitempty"`
	Rel bool   `protobuf:"varint,2,opt,name=Rel" json:"Rel,omitempty"`
}

func (m *FanSpeedTarget) Reset()                    { *m = FanSpeedTarget{} }
func (m *FanSpeedTarget) String() string            { return proto.CompactTextString(m) }
func (*FanSpeedTarget) ProtoMessage()               {}
func (*FanSpeedTarget) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{10} }

func (m *FanSpeedTarget) GetAbs() uint32 {
	if m != nil {
		return m.Abs
	}
	return 0
}

func (m *FanSpeedTarget) GetRel() bool {
	if m != nil {
		return m.Rel
	}
	return false
}

type FanSpeedModified struct {
	Modified      bool   `protobuf:"varint,1,opt,name=Modified" json:"Modified,omitempty"`
	OriginalSpeed uint32 `protobuf:"varint,2,opt,name=OriginalSpeed" json:"OriginalSpeed,omitempty"`
	TargetSpeed   uint32 `protobuf:"varint,3,opt,name=TargetSpeed" json:"TargetSpeed,omitempty"`
	ReqTime       string `protobuf:"bytes,4,opt,name=ReqTime" json:"ReqTime,omitempty"`
}

func (m *FanSpeedModified) Reset()                    { *m = FanSpeedModified{} }
func (m *FanSpeedModified) String() string            { return proto.CompactTextString(m) }
func (*FanSpeedModified) ProtoMessage()               {}
func (*FanSpeedModified) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{11} }

func (m *FanSpeedModified) GetModified() bool {
	if m != nil {
		return m.Modified
	}
	return false
}

func (m *FanSpeedModified) GetOriginalSpeed() uint32 {
	if m != nil {
		return m.OriginalSpeed
	}
	return 0
}

func (m *FanSpeedModified) GetTargetSpeed() uint32 {
	if m != nil {
		return m.TargetSpeed
	}
	return 0
}

func (m *FanSpeedModified) GetReqTime() string {
	if m != nil {
		return m.ReqTime
	}
	return ""
}

type Hours struct {
	FanAway      uint32 `protobuf:"varint,1,opt,name=FanAway" json:"FanAway,omitempty"`
	FanLow       uint32 `protobuf:"varint,2,opt,name=FanLow" json:"FanLow,omitempty"`
	FanMid       uint32 `protobuf:"varint,3,opt,name=FanMid" json:"FanMid,omitempty"`
	FanHigh      uint32 `protobuf:"varint,4,opt,name=FanHigh" json:"FanHigh,omitempty"`
	FrostProtect uint32 `protobuf:"varint,5,opt,name=FrostProtect" json:"FrostProtect,omitempty"`
	Reheating    uint32 `protobuf:"varint,6,opt,name=Reheating" json:"Reheating,omitempty"`
	BypassOpen   uint32 `protobuf:"varint,7,opt,name=BypassOpen" json:"BypassOpen,omitempty"`
	Filter       uint32 `protobuf:"varint,8,opt,name=Filter" json:"Filter,omitempty"`
}

func (m *Hours) Reset()                    { *m = Hours{} }
func (m *Hours) String() string            { return proto.CompactTextString(m) }
func (*Hours) ProtoMessage()               {}
func (*Hours) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{12} }

func (m *Hours) GetFanAway() uint32 {
	if m != nil {
		return m.FanAway
	}
	return 0
}

func (m *Hours) GetFanLow() uint32 {
	if m != nil {
		return m.FanLow
	}
	return 0
}

func (m *Hours) GetFanMid() uint32 {
	if m != nil {
		return m.FanMid
	}
	return 0
}

func (m *Hours) GetFanHigh() uint32 {
	if m != nil {
		return m.FanHigh
	}
	return 0
}

func (m *Hours) GetFrostProtect() uint32 {
	if m != nil {
		return m.FrostProtect
	}
	return 0
}

func (m *Hours) GetReheating() uint32 {
	if m != nil {
		return m.Reheating
	}
	return 0
}

func (m *Hours) GetBypassOpen() uint32 {
	if m != nil {
		return m.BypassOpen
	}
	return 0
}

func (m *Hours) GetFilter() uint32 {
	if m != nil {
		return m.Filter
	}
	return 0
}

type Temps struct {
	Comfort     float32 `protobuf:"fixed32,1,opt,name=Comfort" json:"Comfort,omitempty"`
	OutsideAir  float32 `protobuf:"fixed32,2,opt,name=OutsideAir" json:"OutsideAir,omitempty"`
	SupplyAir   float32 `protobuf:"fixed32,3,opt,name=SupplyAir" json:"SupplyAir,omitempty"`
	OutAir      float32 `protobuf:"fixed32,4,opt,name=OutAir" json:"OutAir,omitempty"`
	ExhaustAir  float32 `protobuf:"fixed32,5,opt,name=ExhaustAir" json:"ExhaustAir,omitempty"`
	GeoHeat     float32 `protobuf:"fixed32,6,opt,name=GeoHeat" json:"GeoHeat,omitempty"`
	Reheating   float32 `protobuf:"fixed32,7,opt,name=Reheating" json:"Reheating,omitempty"`
	KitchenHood float32 `protobuf:"fixed32,8,opt,name=KitchenHood" json:"KitchenHood,omitempty"`
}

func (m *Temps) Reset()                    { *m = Temps{} }
func (m *Temps) String() string            { return proto.CompactTextString(m) }
func (*Temps) ProtoMessage()               {}
func (*Temps) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{13} }

func (m *Temps) GetComfort() float32 {
	if m != nil {
		return m.Comfort
	}
	return 0
}

func (m *Temps) GetOutsideAir() float32 {
	if m != nil {
		return m.OutsideAir
	}
	return 0
}

func (m *Temps) GetSupplyAir() float32 {
	if m != nil {
		return m.SupplyAir
	}
	return 0
}

func (m *Temps) GetOutAir() float32 {
	if m != nil {
		return m.OutAir
	}
	return 0
}

func (m *Temps) GetExhaustAir() float32 {
	if m != nil {
		return m.ExhaustAir
	}
	return 0
}

func (m *Temps) GetGeoHeat() float32 {
	if m != nil {
		return m.GeoHeat
	}
	return 0
}

func (m *Temps) GetReheating() float32 {
	if m != nil {
		return m.Reheating
	}
	return 0
}

func (m *Temps) GetKitchenHood() float32 {
	if m != nil {
		return m.KitchenHood
	}
	return 0
}

type FlushCacheRequest struct {
	// The type of cache to flush. One of all/fans/temps/profiles.
	Type string `protobuf:"bytes,1,opt,name=Type" json:"Type,omitempty"`
}

func (m *FlushCacheRequest) Reset()                    { *m = FlushCacheRequest{} }
func (m *FlushCacheRequest) String() string            { return proto.CompactTextString(m) }
func (*FlushCacheRequest) ProtoMessage()               {}
func (*FlushCacheRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{14} }

func (m *FlushCacheRequest) GetType() string {
	if m != nil {
		return m.Type
	}
	return ""
}

type FlushCacheResponse struct {
	Success bool   `protobuf:"varint,1,opt,name=Success" json:"Success,omitempty"`
	ReqTime string `protobuf:"bytes,2,opt,name=ReqTime" json:"ReqTime,omitempty"`
}

func (m *FlushCacheResponse) Reset()                    { *m = FlushCacheResponse{} }
func (m *FlushCacheResponse) String() string            { return proto.CompactTextString(m) }
func (*FlushCacheResponse) ProtoMessage()               {}
func (*FlushCacheResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{15} }

func (m *FlushCacheResponse) GetSuccess() bool {
	if m != nil {
		return m.Success
	}
	return false
}

func (m *FlushCacheResponse) GetReqTime() string {
	if m != nil {
		return m.ReqTime
	}
	return ""
}

func init() {
	proto.RegisterType((*Noop)(nil), "comfo.comfoserver.Noop")
	proto.RegisterType((*Bypass)(nil), "comfo.comfoserver.Bypass")
	proto.RegisterType((*BootInfo)(nil), "comfo.comfoserver.BootInfo")
	proto.RegisterType((*ComfortTarget)(nil), "comfo.comfoserver.ComfortTarget")
	proto.RegisterType((*ComfortModified)(nil), "comfo.comfoserver.ComfortModified")
	proto.RegisterType((*Errors)(nil), "comfo.comfoserver.Errors")
	proto.RegisterType((*Fans)(nil), "comfo.comfoserver.Fans")
	proto.RegisterType((*FanProfiles)(nil), "comfo.comfoserver.FanProfiles")
	proto.RegisterType((*FanProfileTarget)(nil), "comfo.comfoserver.FanProfileTarget")
	proto.RegisterType((*FanProfileModified)(nil), "comfo.comfoserver.FanProfileModified")
	proto.RegisterType((*FanSpeedTarget)(nil), "comfo.comfoserver.FanSpeedTarget")
	proto.RegisterType((*FanSpeedModified)(nil), "comfo.comfoserver.FanSpeedModified")
	proto.RegisterType((*Hours)(nil), "comfo.comfoserver.Hours")
	proto.RegisterType((*Temps)(nil), "comfo.comfoserver.Temps")
	proto.RegisterType((*FlushCacheRequest)(nil), "comfo.comfoserver.FlushCacheRequest")
	proto.RegisterType((*FlushCacheResponse)(nil), "comfo.comfoserver.FlushCacheResponse")
}

func init() { proto.RegisterFile("github.com/ti-mo/comfo/rpc/comfo/service.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 996 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xcc, 0x56, 0xdd, 0x6e, 0x1b, 0x45,
	0x14, 0xd6, 0x3a, 0xb6, 0xe3, 0x1c, 0xd7, 0xa5, 0x5d, 0x21, 0x30, 0x11, 0xaa, 0xcc, 0xb6, 0x88,
	0xde, 0xe0, 0x88, 0x9f, 0x3b, 0x84, 0x50, 0x12, 0xea, 0xd8, 0xd0, 0xd4, 0xd1, 0x3a, 0x54, 0x02,
	0xae, 0x36, 0xeb, 0x13, 0x7b, 0x90, 0x3d, 0xb3, 0x9d, 0x9d, 0x4d, 0x09, 0xaf, 0xc0, 0x0d, 0x42,
	0xbc, 0x18, 0x0f, 0xc1, 0x03, 0x70, 0xc3, 0x35, 0x3a, 0x33, 0x67, 0x77, 0xc7, 0x8d, 0x13, 0x71,
	0xc9, 0x8d, 0xb5, 0xdf, 0x77, 0x66, 0xce, 0xdf, 0x7c, 0x73, 0xc6, 0x30, 0x5c, 0x08, 0xb3, 0x2c,
	0x2e, 0x86, 0xa9, 0x5a, 0x1f, 0x18, 0xf1, 0xf1, 0x5a, 0x1d, 0xa4, 0x6a, 0x7d, 0xa9, 0x0e, 0x74,
	0x96, 0xf2, 0x57, 0x8e, 0xfa, 0x4a, 0xa4, 0x38, 0xcc, 0xb4, 0x32, 0x2a, 0x7c, 0x68, 0xc9, 0xa1,
	0xfd, 0x25, 0x0b, 0xea, 0xa8, 0x0d, 0xcd, 0x17, 0x4a, 0x65, 0xd1, 0x15, 0xb4, 0x8f, 0xae, 0xb3,
	0x24, 0xcf, 0xc3, 0x77, 0xa0, 0x3d, 0x4a, 0x52, 0xa3, 0x74, 0x3f, 0x18, 0x04, 0x4f, 0x7b, 0x31,
	0xa3, 0xf0, 0x6d, 0x68, 0x3d, 0xc7, 0x2b, 0x5c, 0xf5, 0x1b, 0x96, 0x76, 0x20, 0x7c, 0x04, 0x70,
	0xac, 0xb4, 0xc6, 0xd4, 0x08, 0x25, 0xfb, 0x3b, 0xd6, 0xe4, 0x31, 0x64, 0x9f, 0x15, 0xeb, 0x35,
	0xea, 0x53, 0x35, 0xc7, 0x7e, 0x73, 0x10, 0x3c, 0xed, 0xc4, 0x1e, 0x13, 0xfd, 0x11, 0x40, 0xe7,
	0x48, 0x29, 0x33, 0x91, 0x97, 0x2a, 0x8c, 0xe0, 0xde, 0x69, 0xf2, 0x93, 0xd2, 0x2f, 0x51, 0xe7,
	0xe4, 0xce, 0x25, 0xb0, 0xc1, 0xd9, 0x35, 0x42, 0xd6, 0x6b, 0x1a, 0xbc, 0xc6, 0xe3, 0xc2, 0x01,
	0x74, 0x8f, 0xd0, 0x24, 0xe5, 0x12, 0x97, 0x95, 0x4f, 0x51, 0x5a, 0x5f, 0x23, 0x75, 0xe6, 0x45,
	0xb2, 0x76, 0x69, 0xed, 0xc5, 0x1e, 0x13, 0x7d, 0x02, 0xbd, 0x63, 0xea, 0x92, 0x36, 0xe7, 0x89,
	0x5e, 0xa0, 0x21, 0x97, 0x25, 0x81, 0xeb, 0x8c, 0x33, 0xf3, 0xa9, 0xe8, 0xd7, 0x00, 0xde, 0x62,
	0x7c, 0xaa, 0xe6, 0xe2, 0x52, 0xe0, 0x3c, 0xdc, 0x87, 0x4e, 0xf9, 0x6d, 0xb7, 0x74, 0xe2, 0x0a,
	0x53, 0x21, 0x53, 0x2d, 0x16, 0x42, 0x26, 0x2b, 0xeb, 0x92, 0x0b, 0xf1, 0x39, 0x4a, 0xd3, 0xc5,
	0xb7, 0x2b, 0xb8, 0xbb, 0x35, 0x13, 0xf6, 0x61, 0x37, 0xc6, 0x57, 0xe7, 0xa2, 0xaa, 0xa1, 0x84,
	0xd1, 0x00, 0xda, 0xcf, 0xb4, 0x56, 0xda, 0x9d, 0xa7, 0x58, 0x19, 0xd4, 0x9c, 0x01, 0xa3, 0xe8,
	0x17, 0x68, 0x8e, 0x12, 0x99, 0x87, 0xef, 0xc3, 0xde, 0x44, 0x9e, 0xa1, 0x4e, 0x51, 0x1a, 0xae,
	0xab, 0x26, 0x28, 0x83, 0x69, 0x61, 0x4a, 0xb3, 0xcb, 0xd1, 0x63, 0x28, 0x83, 0x89, 0x9c, 0x65,
	0x88, 0x73, 0x4e, 0xaf, 0x84, 0x54, 0xfb, 0xb4, 0x30, 0xce, 0xd4, 0xb4, 0xa6, 0x0a, 0x47, 0x7f,
	0x36, 0xa0, 0x3b, 0x4a, 0xe4, 0x99, 0x56, 0x97, 0x62, 0x85, 0x39, 0x79, 0x99, 0x16, 0xe6, 0xf0,
	0x75, 0x72, 0xcd, 0x19, 0x94, 0x90, 0xb2, 0x9f, 0x16, 0xe6, 0xb9, 0x7a, 0xcd, 0xb1, 0x19, 0x31,
	0x7f, 0x2a, 0xca, 0xb0, 0x8c, 0xd8, 0xd3, 0x58, 0x2c, 0x96, 0x1c, 0xb4, 0x84, 0x74, 0x82, 0x13,
	0x39, 0x4a, 0xe4, 0x61, 0x6a, 0xc4, 0x15, 0xf6, 0x5b, 0xb6, 0x19, 0x3e, 0x45, 0x3e, 0x27, 0xd2,
	0x26, 0xd1, 0x76, 0x3e, 0x1d, 0x22, 0xe5, 0x4f, 0x24, 0xa5, 0xb0, 0xeb, 0x94, 0x6f, 0x81, 0x63,
	0x29, 0x81, 0x4e, 0xc9, 0x52, 0x7c, 0xeb, 0xc3, 0x86, 0xdf, 0x2b, 0x7d, 0xd8, 0xe8, 0x74, 0x4f,
	0x0a, 0xad, 0x51, 0x9a, 0x69, 0x61, 0xfa, 0xc0, 0xf7, 0xa4, 0x62, 0xe8, 0x14, 0x18, 0x4d, 0x64,
	0xbf, 0xeb, 0x4e, 0xa1, 0x22, 0x48, 0x2b, 0x0c, 0xdc, 0x15, 0xbc, 0xe7, 0xb4, 0xe2, 0x73, 0xd1,
	0x37, 0xf0, 0xa0, 0x6e, 0x29, 0xab, 0xb6, 0xba, 0xb3, 0x81, 0x7f, 0x67, 0x07, 0xd0, 0x75, 0x76,
	0x77, 0x38, 0xae, 0xb1, 0x3e, 0x15, 0xfd, 0x1e, 0x40, 0x58, 0x3b, 0xfb, 0x4f, 0x72, 0x7e, 0x02,
	0xbd, 0x52, 0xba, 0xbe, 0xdb, 0x4d, 0xf2, 0xcd, 0xd0, 0x3b, 0x37, 0x42, 0xdf, 0x21, 0xe9, 0xcf,
	0xe1, 0xfe, 0x28, 0x71, 0xe2, 0xe2, 0xf2, 0x1e, 0xc0, 0xce, 0xe1, 0x45, 0xce, 0xc5, 0xd1, 0x27,
	0x31, 0x31, 0x8f, 0xa8, 0x4e, 0x4c, 0x9f, 0xd1, 0x6f, 0x81, 0xed, 0x8b, 0xdd, 0xf6, 0x3f, 0x29,
	0xe4, 0xaf, 0x00, 0x5a, 0x63, 0x55, 0x68, 0xab, 0x7b, 0x92, 0x9f, 0xa7, 0x7b, 0x86, 0x6e, 0x0a,
	0x4b, 0x4f, 0xf7, 0x0e, 0x31, 0xef, 0xe9, 0xde, 0x21, 0xf6, 0xe4, 0xeb, 0x9e, 0x21, 0x69, 0x67,
	0xa4, 0x55, 0x6e, 0xce, 0xb4, 0x32, 0x98, 0x1a, 0x2b, 0xfc, 0x5e, 0xbc, 0xc1, 0x91, 0xfa, 0x62,
	0x5c, 0x62, 0x62, 0x84, 0x5c, 0xb0, 0xf8, 0x6b, 0x82, 0xb4, 0xeb, 0xde, 0x86, 0x69, 0x86, 0x92,
	0x2f, 0x81, 0xc7, 0x78, 0x13, 0xa6, 0xc3, 0x39, 0xb9, 0x09, 0xf3, 0x77, 0x00, 0x2d, 0x1a, 0x53,
	0xb6, 0x4e, 0x1e, 0x8d, 0xb6, 0xce, 0x46, 0x5c, 0x42, 0x9e, 0x2f, 0xb9, 0x98, 0xe3, 0xa1, 0xd0,
	0xb6, 0xd6, 0x46, 0xec, 0x31, 0x94, 0xd9, 0xac, 0xc8, 0xb2, 0xd5, 0x35, 0x99, 0x77, 0xac, 0xb9,
	0x26, 0x78, 0x0a, 0x90, 0xa9, 0x69, 0x4d, 0x8c, 0xc8, 0xeb, 0xb3, 0x9f, 0x97, 0x49, 0x91, 0x5b,
	0x5b, 0xcb, 0x79, 0xad, 0x19, 0xca, 0xe7, 0x04, 0xd5, 0x18, 0x13, 0x63, 0xab, 0x6d, 0xc4, 0x25,
	0xdc, 0xec, 0xc4, 0xae, 0x8b, 0x57, 0x77, 0x62, 0x00, 0xdd, 0x6f, 0x85, 0x49, 0x97, 0x28, 0xc7,
	0x4a, 0xb9, 0x9b, 0xdf, 0x88, 0x7d, 0x2a, 0xfa, 0x08, 0x1e, 0x8e, 0x56, 0x45, 0xbe, 0x3c, 0x4e,
	0xd2, 0x25, 0xc6, 0xf8, 0xaa, 0xc0, 0xdc, 0x84, 0x21, 0x34, 0xcf, 0xaf, 0x33, 0xb4, 0xb5, 0xef,
	0xc5, 0xf6, 0x3b, 0x1a, 0x43, 0xe8, 0x2f, 0xcc, 0x33, 0x25, 0x73, 0xa4, 0xc4, 0x66, 0x45, 0x9a,
	0x62, 0x9e, 0xb3, 0x2e, 0x4b, 0xe8, 0xcb, 0xa9, 0xb1, 0x21, 0xa7, 0x4f, 0xff, 0x69, 0x42, 0xcb,
	0xb6, 0x33, 0xfc, 0x12, 0x3a, 0x27, 0xee, 0x65, 0xc8, 0xc3, 0x77, 0x87, 0x37, 0x1e, 0xfb, 0x21,
	0xbd, 0xf4, 0xfb, 0xfd, 0x2d, 0x06, 0xb7, 0xe5, 0x0b, 0xea, 0x8a, 0xb1, 0x8f, 0xc2, 0xad, 0xbb,
	0xb7, 0x19, 0xec, 0x8e, 0x09, 0xdc, 0x77, 0x9b, 0xab, 0xa1, 0x7e, 0xab, 0x8f, 0x47, 0xdb, 0x7d,
	0x54, 0x1b, 0xbf, 0x82, 0xbd, 0x13, 0x34, 0xfc, 0x7c, 0xdd, 0xea, 0xe5, 0xbd, 0x2d, 0x06, 0xde,
	0xf3, 0x1d, 0x74, 0x67, 0x36, 0x17, 0x77, 0x13, 0x3f, 0xd8, 0x1e, 0xcf, 0x9b, 0x24, 0xfb, 0x8f,
	0xef, 0x58, 0x52, 0x4d, 0x86, 0x97, 0x70, 0x7f, 0x86, 0xc6, 0x7b, 0xf3, 0xc3, 0xc1, 0x96, 0x6d,
	0x1b, 0xff, 0x1b, 0xf6, 0xa3, 0xdb, 0x57, 0x54, 0x7e, 0x7f, 0x84, 0xde, 0xcc, 0x6f, 0x5d, 0xf8,
	0xf8, 0xce, 0x06, 0xb1, 0xe7, 0x0f, 0xef, 0x5c, 0x54, 0x39, 0xff, 0x1e, 0xa0, 0xd6, 0x59, 0xf8,
	0x64, 0xdb, 0xa6, 0x37, 0xf5, 0xba, 0xdd, 0xf5, 0x0d, 0xb1, 0x1e, 0xed, 0xfe, 0xd0, 0xb2, 0x2b,
	0x2e, 0xda, 0xf6, 0xef, 0xe5, 0x67, 0xff, 0x06, 0x00, 0x00, 0xff, 0xff, 0xdc, 0x8e, 0xae, 0x26,
	0x90, 0x0a, 0x00, 0x00,
}
