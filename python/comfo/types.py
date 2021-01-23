"""
Concrete types for converting from protobuf objects.

Since the Python protobuf compiler outputs code that generates
classes on the fly from a symbol database, we convert them as
concrete types here to ensure the expected fields are present.
"""
from pprint import pformat


class ComfoMessage:
    """Mixin for protobuf wrappers in this package."""

    def __repr__(self):
        """Print the object as a dict."""
        return pformat(self.__dict__, sort_dicts=False)


class BootInfo(ComfoMessage):
    """BootInfo is a concrete type embedding protobuf values."""

    def __init__(self, pb):
        """Copy the protobuf's values to instance variables."""
        self.MajorVersion = pb.MajorVersion
        self.MinorVersion = pb.MinorVersion
        self.BetaVersion = pb.BetaVersion
        self.DeviceName = pb.DeviceName


class Bypass(ComfoMessage):
    """Bypass is a concrete type embedding protobuf values."""

    def __init__(self, pb):
        """Copy the protobuf's values to instance variables."""
        self.Factor = pb.Factor
        self.Level = pb.Level
        self.Correction = pb.Correction
        self.SummerMode = pb.SummerMode


class Temps(ComfoMessage):
    """Temps is a concrete type embedding protobuf values."""

    def __init__(self, pb):
        """Copy the protobuf's values to instance variables."""
        self.Comfort = pb.Comfort
        self.OutsideAir = pb.OutsideAir
        self.SupplyAir = pb.SupplyAir
        self.InsideAir = pb.InsideAir
        self.ExhaustAir = pb.ExhaustAir
        self.GeoHeat = pb.GeoHeat
        self.Reheating = pb.Reheating
        self.KitchenHood = pb.KitchenHood


class Fans(ComfoMessage):
    """Fans is a concrete type embedding protobuf values."""

    def __init__(self, pb):
        """Copy the protobuf's values to instance variables."""
        self.InPercent = pb.InPercent
        self.OutPercent = pb.OutPercent
        self.InSpeed = pb.InSpeed
        self.OutSpeed = pb.OutSpeed


class FanProfiles(ComfoMessage):
    """FanProfiles is a concrete type embedding protobuf values."""

    SPEED_OFF = 1
    SPEED_LOW = 2
    SPEED_MEDIUM = 3
    SPEED_HIGH = 4

    VALID_SPEEDS = [SPEED_OFF, SPEED_LOW, SPEED_MEDIUM, SPEED_HIGH]

    def __init__(self, pb):
        """Copy the protobuf's values to instance variables."""
        self.OutAway = pb.OutAway
        self.OutLow = pb.OutLow
        self.OutMid = pb.OutMid
        self.OutHigh = pb.OutHigh

        self.InFanActive = pb.InFanActive
        self.InAway = pb.InAway
        self.InLow = pb.InLow
        self.InMid = pb.InMid
        self.InHigh = pb.InHigh

        self.CurrentOut = pb.CurrentOut
        self.CurrentIn = pb.CurrentIn
        self.CurrentMode = pb.CurrentMode


class Errors(ComfoMessage):
    """Errors is a concrete type embedding protobuf values."""

    def __init__(self, pb):
        """Copy the protobuf's values to instance variables."""
        self.Filter = pb.Filter
