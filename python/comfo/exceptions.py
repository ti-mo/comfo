"""Exceptions thrown by the Comfo API client."""


class InvalidFanSpeed(Exception):
    """Raised when invalid fan speed profile is given."""


class InvalidDutyPercent(Exception):
    """Raised when invalid duty level percentage is given."""


class InvalidTemperature(Exception):
    """Raised when invalid temperature input is given."""
