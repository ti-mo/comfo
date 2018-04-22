package comfoserver

import (
	"fmt"
	"net"
	"strings"

	rpc "github.com/ti-mo/comfo/rpc/comfo"
	"github.com/twitchtv/twirp"
)

// modifySpeed takes the unit's original speed (baseSpeed) and a protobuf FanSpeedTarget.
// Based on the target a new fan speed is returned between 1 and 4.
func modifySpeed(baseSpeed uint8, target *rpc.FanSpeedTarget) (tgtSpeed uint8, err error) {

	// Unit has 4 speed settings
	// 0 means auto, we don't use it
	var lowerBound uint8 = 1
	var upperBound uint8 = 4

	// Make sure only one of Abs and Rel is set
	if target.Abs != 0 && target.Rel != "" {
		return 0, twirp.InvalidArgumentError("Abs/Rel", errBothAbsRel.Error())
	} else if target.Abs == 0 && target.Rel == "" {
		return 0, twirp.InvalidArgumentError("Abs/Rel", errNoneAbsRel.Error())
	}

	// Determine Abs/Rel speed and target speed
	if target.Abs != 0 {
		tgtSpeed = uint8(target.Abs)
	} else if target.Rel != "" {
		if target.Rel == "+" {
			tgtSpeed = baseSpeed + 1
		} else if target.Rel == "-" {
			tgtSpeed = baseSpeed - 1
		} else {
			return 0, twirp.InvalidArgumentError("Rel", fmt.Sprintf("unknown value '%v'", target.Rel))
		}
	}

	// Bounds check
	if tgtSpeed < lowerBound || tgtSpeed > upperBound {
		return baseSpeed, twirp.InvalidArgumentError("FanSpeed", fmt.Sprintf("value '%v' out of range", tgtSpeed))
	}

	return
}

// ifaceAddrs returns a list of ipv4 and ipv6 addresses of the host.
func ifaceAddrs() (v4 []net.IP, v6 []net.IP) {

	// Get system interface addresses
	ifaces, _ := net.InterfaceAddrs()

	for _, i := range ifaces {
		ip, _, _ := net.ParseCIDR(i.String())

		if ip != nil && strings.Contains(ip.String(), ":") {
			v6 = append(v6, ip)
		} else {
			v4 = append(v4, ip)
		}
	}

	return
}

// printEndpoints prints a list of addresses the API is reachable on.
func printEndpoints(port string) {

	// Get host interface addresses
	v4, v6 := ifaceAddrs()

	fmt.Println("\nAPI listening on following endpoints:")

	fmt.Println("   IPv4:")
	for _, a := range v4 {
		fmt.Printf("    - http://%s:%s\n", a, port)
	}
	fmt.Println()

	fmt.Println("   IPv6:")
	for _, a := range v6 {
		fmt.Printf("    - http://[%s]:%s\n", a, port)
	}
	fmt.Println()
}
