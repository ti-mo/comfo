package comfoserver

import (
	"fmt"

	rpc "github.com/ti-mo/comfo/rpc/comfo"
	"github.com/twitchtv/twirp"
)

// modifySpeed takes the unit's original speed (baseSpeed) and a protobuf FanSpeedTarget.
// Based on the target a new fan speed is returned between 1 and 4. An empty FanSpeedTarget
// will default to a relative speed decrease.
func modifySpeed(baseSpeed uint8, target *rpc.FanSpeedTarget) (tgtSpeed uint8, err error) {

	// Unit has 4 speed settings
	// 0 means auto, we don't use it
	var lowerBound uint8 = 1
	var upperBound uint8 = 4

	// Make sure only one of Abs and Rel is set
	if target.Abs != 0 && target.Rel {
		return 0, twirp.InvalidArgumentError("Abs/Rel", errBothAbsRel.Error())
	}

	// Determine Abs/Rel speed and target speed
	if target.Abs != 0 {
		tgtSpeed = uint8(target.Abs)
	} else if target.Rel {
		tgtSpeed = baseSpeed + 1
	} else {
		// Default decrease
		tgtSpeed = baseSpeed - 1
	}

	// Bounds check
	if tgtSpeed < lowerBound || tgtSpeed > upperBound {
		return baseSpeed, twirp.InvalidArgumentError("FanSpeed", fmt.Sprintf("value '%v' out of range", tgtSpeed))
	}

	return
}
