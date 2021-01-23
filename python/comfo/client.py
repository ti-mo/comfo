"""Comfo Twirp API client."""
import asyncio
import functools

from twirp.context import Context
from twirp.errors import Errors as TwirpErrors
from twirp.exceptions import TwirpServerException

from . import comfo_pb2, exceptions
from .comfo_twirp import ComfoClient
from .types import BootInfo, Bypass, Errors, FanProfiles, Fans, Temps


class Comfo:
    """
    Comfo implements the comfo Python API.

    All methods that perform any form of network i/o have async variants
    denoted by the 'async_' prefix.
    """

    def __init__(self, host, loop=None):
        """Construct a new Comfo."""
        self.host = host
        self.twirp = ComfoClient(f"http://{host}:3094")

        if loop is None:
            try:
                self.loop = asyncio.get_running_loop()
            except RuntimeError:
                pass

    def ping(self) -> bool:
        """Run a simple get query against the unit to verify connectivity."""
        self.get_bootinfo()
        return True

    async def async_ping(self):
        """Async version of ping()."""
        return await self.loop.run_in_executor(None, self.ping)

    def get_bootinfo(self) -> BootInfo:
        """
        Get the unit's BootInfo.

        This contains the firmware's version and the name of the unit.
        """
        return BootInfo(
            self.twirp.GetBootInfo(
                ctx=Context(),
                request=comfo_pb2.Noop(),
            )
        )

    async def async_get_bootinfo(self) -> BootInfo:
        """Async version of get_bootinfo()."""
        return await self.loop.run_in_executor(None, self.get_bootinfo)

    def get_bypass(self) -> Bypass:
        """Get the unit's Bypass (heat exchanger) status."""
        return Bypass(
            self.twirp.GetBypass(
                ctx=Context(),
                request=comfo_pb2.Noop(),
            )
        )

    async def async_get_bypass(self) -> Bypass:
        """Async version of get_bypass()."""
        return await self.loop.run_in_executor(None, self.get_bypass)

    def get_temps(self) -> Temps:
        """Get the unit's temperature sensor data."""
        return Temps(
            self.twirp.GetTemps(
                ctx=Context(),
                request=comfo_pb2.Noop(),
            )
        )

    async def async_get_temps(self) -> Temps:
        """Async version of get_temps()."""
        return await self.loop.run_in_executor(None, self.get_temps)

    def get_fans(self) -> Fans:
        """Get the unit's fan speeds and duties."""
        return Fans(
            self.twirp.GetFans(
                ctx=Context(),
                request=comfo_pb2.Noop(),
            )
        )

    async def async_get_fans(self) -> Fans:
        """Async version of get_fans()."""
        return await self.loop.run_in_executor(None, self.get_fans)

    def get_fan_profiles(self) -> FanProfiles:
        """Get the unit's fan speeds and duties."""
        return FanProfiles(
            self.twirp.GetFanProfiles(
                ctx=Context(),
                request=comfo_pb2.Noop(),
            )
        )

    async def async_get_fan_profiles(self) -> FanProfiles:
        """Async version of get_fan_profiles()."""
        return await self.loop.run_in_executor(None, self.get_fan_profiles)

    def get_errors(self) -> Errors:
        """Get the unit's error statuses."""
        return Errors(
            self.twirp.GetErrors(
                ctx=Context(),
                request=comfo_pb2.Noop(),
            )
        )

    async def async_get_errors(self) -> Errors:
        """Async version of get_errors()."""
        return await self.loop.run_in_executor(None, self.get_errors)

    def set_fan_speed(
        self,
        speed: int,
    ) -> bool:
        """Set the unit's fan speed to a specific value.

        The given speed must be a value between 1 and 4. The FanProfiles class
        provides SPEED_ constants to use as values.

        Returns True if the call caused the unit's speed to be modified.
        """
        if speed not in FanProfiles.VALID_SPEEDS:
            raise exceptions.InvalidFanSpeed

        return self.twirp.SetFanSpeed(
            ctx=Context(),
            request=comfo_pb2.FanSpeedTarget(Abs=speed),
        ).Modified

    async def async_set_fan_speed(self, speed: int) -> bool:
        """Async version of set_fan_speed()."""
        return await self.loop.run_in_executor(
            None,
            functools.partial(self.set_fan_speed, speed=speed),
        )

    def increase_fan_speed(self) -> bool:
        """Activate the unit's next duty level.

        Returns True if the call caused the unit's speed to be modified.
        """
        try:
            return self.twirp.SetFanSpeed(
                ctx=Context(),
                request=comfo_pb2.FanSpeedTarget(Rel=True),
            ).Modified
        except TwirpServerException as e:
            # InvalidArgument means the fan duty went out of bounds
            # and wasn't modified.
            if e.code is TwirpErrors.InvalidArgument:
                return False
            raise

    async def async_increase_fan_speed(self) -> bool:
        """Async version of increase_fan_speed()."""
        return await self.loop.run_in_executor(None, self.increase_fan_speed)

    def decrease_fan_speed(self) -> bool:
        """Activate the unit's lower duty level.

        Returns True if the call caused the unit's speed to be modified.
        """
        try:
            return self.twirp.SetFanSpeed(
                ctx=Context(),
                request=comfo_pb2.FanSpeedTarget(Rel=False),
            ).Modified
        except TwirpServerException as e:
            # InvalidArgument means the fan duty went out of bounds
            # and wasn't modified.
            if e.code is TwirpErrors.InvalidArgument:
                return False
            raise

        return True

    async def async_decrease_fan_speed(self) -> bool:
        """Async version of decrease_fan_speed()."""
        return await self.loop.run_in_executor(None, self.decrease_fan_speed)

    def set_comfort_temperature(self, temperature: int = None) -> bool:
        """
        Set the comfort temperature on the unit's heat exchanger.

        The given temperature must range between 0-255.
        """
        if temperature is None or temperature < 0 or temperature > 255:
            raise exceptions.InvalidTemperature

        return self.twirp.SetComfortTemp(
            ctx=Context(),
            request=comfo_pb2.ComfortTarget(
                ComfortTemp=temperature,
            ),
        ).Modified

    async def async_set_comfort_temperature(self, temperature: int) -> bool:
        """Async version of set_comfort_temperature()."""
        return await self.loop.run_in_executor(
            None,
            functools.partial(
                self.set_comfort_temperature,
                temperature=temperature,
            ),
        )

    def configure_fan_profile(self, profile: int, percent: int) -> bool:
        """Configure a fan profile's duty level, in percent.

        A unit has 4 profiles, addressed from 1-4. Each profile is associated
        with a duty level (percentage) at which to run the unit's fan.
        The percentage must range between 15 and 100.
        """
        if profile not in FanProfiles.VALID_SPEEDS:
            raise exceptions.InvalidFanSpeed

        if percent < 15 or percent > 100:
            raise exceptions.InvalidDutyPercent

        return self.twirp.SetFanProfile(
            ctx=Context(),
            request=comfo_pb2.FanProfileTarget(
                Mode=profile,
                TargetSpeed=percent,
            ),
        ).Modified

    async def async_configure_fan_profile(
        self,
        profile: int,
        percent: int,
    ) -> bool:
        """Async version of configure_fan_profile()."""
        return await self.loop.run_in_executor(
            None,
            functools.partial(
                self.configure_fan_profile,
                profile=profile,
                percent=percent,
            ),
        )
