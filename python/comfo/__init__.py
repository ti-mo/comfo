"""
Comfo client to be used against the Go-based comfoserver API.

Wraps the generated Twirp (protobuf) methods to make them more Pythonic.
All public methods have asyncio variants with the 'async_' prefix.
"""

from .client import Comfo  # noqa: F401
