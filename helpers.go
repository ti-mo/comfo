package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"strconv"
	"strings"
)

// jsonError sends an error in a json response body.
func jsonError(w http.ResponseWriter, err error) {

	resp := map[string]interface{}{
		"error":  err.Error(),
		"failed": true,
	}

	w.WriteHeader(http.StatusInternalServerError)
	jsonWrite(w, resp)
}

// jsonWrite writes data as a json response body.
func jsonWrite(w io.Writer, body interface{}) {

	// Encode and write the body to the Writer
	if err := json.NewEncoder(w).Encode(body); err != nil {
		panic(err)
	}
}

// modifySpeed takes the unit's original ('base') speed and applies an action
// 'actionStr' to it. (one of up/down/int) An error is returned when actionStr
// fails to be converted or when the result is out of bounds.
func modifySpeed(baseSpeed int, actionStr string) (tgtSpeed int, err error) {

	// Unit has 4 speed settings
	// 0 means auto, we don't use it
	lowerBound := 1
	upperBound := 4

	switch actionStr {
	case "up":
		tgtSpeed = baseSpeed + 1
	case "down":
		tgtSpeed = baseSpeed - 1
	default:
		// Extract target speed integer from string
		tgtSpeed, err = strconv.Atoi(actionStr)
		if err != nil {
			return baseSpeed, err
		}
	}

	// Bounds check
	if tgtSpeed < lowerBound || tgtSpeed > upperBound {
		return baseSpeed, errOutOfRange
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
func printEndpoints(port int) {

	// Get host interface addresses
	v4, v6 := ifaceAddrs()

	fmt.Println("\nAPI Endpoints:")

	fmt.Println("   IPv4:")
	for _, a := range v4 {
		fmt.Printf("    - http://%v:%d\n", a.String(), port)
	}
	fmt.Println()

	fmt.Println("   IPv6:")
	for _, a := range v6 {
		fmt.Printf("    - http://[%v]:%d\n", a.String(), port)
	}
	fmt.Println()
}
