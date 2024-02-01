/*
Copyright (c) 2024 Dell Inc., or its subsidiaries. All Rights Reserved.
Licensed under the Mozilla Public License Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://mozilla.org/MPL/2.0/
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package utils

import (
	"errors"
	"fmt"
	"net"
	"net/netip"
	"strings"
)

// IPRange represents a continuous range of IPs
type IPRange struct {
	start netip.Addr
	stop  netip.Addr
	cidr  *net.IPNet
}

// Contains check if a net.IP is contained in the IP range
func (ir IPRange) Contains(ip net.IP) bool {
	if ir.cidr != nil {
		return ir.cidr.Contains(ip)
	}
	addr, _ := netip.ParseAddr(ip.String())
	return ir.start.Compare(addr)*ir.stop.Compare(addr) <= 0
}

// ParseNetwork converts a string to IPRange if valid
func ParseNetwork(network string) (IPRange, error) {
	var netRange IPRange
	network = strings.TrimSpace(network)
	if strings.Contains(network, "/") {
		// a candidate for a CIDR
		_, cidr, err := net.ParseCIDR(network)
		if err != nil {
			return netRange, err
		}
		netRange.cidr = cidr
		return netRange, nil
	}
	start, stop, ok := strings.Cut(network, "-")
	if !ok {
		// a candidate for basic IP
		stop = start
	}
	start, stop = strings.TrimSpace(start), strings.TrimSpace(stop)
	startIP, startErr := netip.ParseAddr(start)
	stopIP, stopErr := netip.ParseAddr(stop)
	err := errors.Join(startErr, stopErr)
	if err != nil {
		return netRange, err
	}
	if startIP.BitLen() != stopIP.BitLen() {
		return netRange, fmt.Errorf("%s has different bit lengths for first and last addresses", network)
	}
	netRange.start, netRange.stop = startIP, stopIP
	return netRange, nil
}

// IPSet resents a set of IPs (need not be continuous range)
type IPSet struct {
	set []IPRange
}

// ParseNetworks converts a slice of strings to IPSet if valid
func ParseNetworks(networks []string) (IPSet, error) {
	set := IPSet{
		set: make([]IPRange, 0),
	}
	var err error
	for _, network := range networks {
		ipr, cerr := ParseNetwork(network)
		if cerr != nil {
			err = errors.Join(cerr, err)
			continue
		}
		set.set = append(set.set, ipr)
	}
	return set, err
}

// Contains check if a net.IP is contained in the IP set
func (is IPSet) Contains(ip net.IP) bool {
	for _, ipr := range is.set {
		if ipr.Contains(ip) {
			return true
		}
	}
	return false
}
