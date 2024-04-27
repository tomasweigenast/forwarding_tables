package main

import (
	"encoding/json"
	"fmt"
	"net"
	"strconv"
	"strings"
)

// ftable represents a forwarding table
type ftable []ftableEntry

// ftableEntry contains information about an entry inside a ftable
type ftableEntry struct {
	dst  net.IPNet // the destination net
	nhop net.IP    // the next hop
	intf string    // the interface name where the data must be send
}

// NewFtable creates a new ftable
func newftable() ftable {
	return make(ftable, 0)
}

func ftableFromJson(j []byte) ftable {
	data := []map[string]any{}
	if err := json.Unmarshal(j, &data); err != nil {
		panic(fmt.Errorf("wrong ftable json format: %s", err))
	}

	ftable := newftable()
	for _, entry := range data {
		dst := net.ParseIP(entry["destination"].(string))
		mask := parseSubnetMask(entry["mask"].(string))
		nextHop := net.ParseIP(entry["next_hop"].(string))
		interf := entry["interface"].(string)

		ftable = append(ftable, ftableEntry{
			dst: net.IPNet{
				IP:   dst,
				Mask: mask,
			},
			nhop: nextHop,
			intf: interf,
		})
	}
	return ftable
}

func parseSubnetMask(maskStr string) net.IPMask {
	// Split the mask string into octets
	octets := strings.Split(maskStr, ".")

	// Convert each octet to uint8
	var octetBytes []byte
	for _, octetStr := range octets {
		octet, err := strconv.Atoi(octetStr)
		if err != nil || octet < 0 || octet > 255 {
			return nil // Return nil if an octet is invalid
		}
		octetBytes = append(octetBytes, byte(octet))
	}

	// Create the subnet mask
	mask := net.IPv4Mask(octetBytes[0], octetBytes[1], octetBytes[2], octetBytes[3])

	return mask
}

func (f *ftable) lookup(ip net.IP) (next_hop net.IP, interface_name string) {
	for _, entry := range *f {
		if entry.dst.Contains(ip) {
			return entry.nhop, entry.intf
		}
	}

	return nil, ""
}

func (ftable *ftable) add(destination, nextHop, interfaceName string) {
	_, destinationNet, err := net.ParseCIDR(destination)
	if err != nil {
		panic(fmt.Errorf("unable to add entry to ftable: %s", err))
	}

	*ftable = append(*ftable, ftableEntry{
		dst:  *destinationNet,
		nhop: net.ParseIP(nextHop),
		intf: interfaceName,
	})
}
