/*
Copyright (c) 2023 Dell Inc., or its subsidiaries. All Rights Reserved.
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
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNetworkParserMultiple(t *testing.T) {
	inputs := []string{
		"192.35.0.1",
		"10.36.0.0-192.36.0.255",
		"fe80::ffff:ffff:ffff:ffff",
		"fe80::ffff:192.0.2.0/125",
		"fe80::ffff:ffff:ffff:1111-fe80::ffff:ffff:ffff:ffff",
		"192.37.0.0/24",
	}

	outputs := map[string]bool{
		"192.35.0.1":                true,
		"10.36.0.20":                true,
		"fe80::ffff:ffff:ffff:ffff": true,
		"fe80::ffff:c000:202":       true,
		"fe80::ffff:192.0.2.11":     false,
		"fe80::ffff:ffff:ffff:1111": true,
		"fe80::ffff:ffff:ffff:f000": true,
		"192.37.0.5":                true,
		"192.39.0.5":                false,
	}
	pool, err := ParseNetworks(inputs)
	assert.Nil(t, err)

	for ipString, result := range outputs {
		ip := net.ParseIP(ipString)
		if ip == nil {
			t.Errorf("%s is not valid", ipString)
		} else {
			ok := pool.Contains(ip)
			if ok {
				assert.Truef(t, result, "%s found", ipString)
			} else {
				assert.Falsef(t, result, "%s not found", ipString)
			}
		}
	}
}

func TestNetworkParser(t *testing.T) {
	TCs := map[string]bool{
		"192.35.0.1":                true,
		"10.36.0.20":                true,
		"10.36.0.20  ":              true,
		"  10.36.0.20":              true,
		" 10.36.0.20 ":              true,
		"fe80::ffff:ffff:ffff:ffff": true,
		"  fe80::ffff:c000:202 ":    true,
		"fe80::ffff:192.0.2.11":     true,
		"fe80::ffff:ffff:ffff:1111": true,
		"10.36.0.0-192.36.0.255":    true,
		"10.36.0.0 - 192.36.0.255":  true,
		"fe80::ffff:ffff:ffff:1111-fe80::ffff:ffff:ffff:ffff":   true,
		"fe80::ffff:ffff:ffff:1111 - fe80::ffff:ffff:ffff:ffff": true,
		"192.37.0.0/24":        true,
		" 192.37.0.0/24 ":      true,
		"hallo":                false,
		"hallow - 192.37.0.0":  false,
		"10.36.0.20 - hallo00": false,
		"hallo/24":             false,
	}
	for v, ok := range TCs {
		_, err := ParseNetwork(v)
		if err != nil {
			assert.Falsef(t, ok, "No error expected for %s, but found %s", v, err.Error())
		} else {
			assert.Truef(t, ok, "Expected error, but none found for %s", v)
		}
	}
}

func TestIPRangePos(t *testing.T) {
	TCs := map[string]string{
		"192.35.0.1":               "192.35.0.1",
		"10.36.0.0-10.36.0.255":    "10.36.0.20",
		"fe80::ffff:192.0.2.0/125": "fe80::ffff:192.0.2.2",
		"fe80::ffff:ffff:ffff:1111 - fe80::ffff:ffff:ffff:ffff": "fe80::ffff:ffff:ffff:112b",
		" 192.37.0.0/24 ": "192.37.0.0",
	}
	for v, ips := range TCs {
		ipr, err := ParseNetwork(v)
		assert.Nil(t, err)
		ip := net.ParseIP(ips)
		assert.NotNilf(t, ip, "%s could not be parsed as an IP", ips)
		assert.Truef(t, ipr.Contains(ip), "%s not found in %s", ips, v)
	}
}

func TestIPRangeNeg(t *testing.T) {
	TCs := map[string]string{
		"192.35.0.1":               "192.35.0.3",
		"10.36.0.0-10.36.0.255":    "10.36.1.20",
		"fe80::ffff:192.0.2.0/125": "fe80::ffff:193.0.2.2",
		"fe80::ffff:ffff:ffff:1111 - fe80::ffff:ffff:ffff:ffff": "fe80::ffff:ffff:ffff:1100",
		"fe80::ffff:ffff:ffff:1111 - fe80::ffff:ffff:ffff:fff0": "fe80::ffff:ffff:ffff:ffff",
		" 192.37.0.0/24 ": "192.37.1.0",
	}
	for v, ips := range TCs {
		ipr, err := ParseNetwork(v)
		assert.Nil(t, err)
		ip := net.ParseIP(ips)
		assert.NotNilf(t, ip, "%s could not be parsed as an IP", ips)
		assert.Falsef(t, ipr.Contains(ip), "%s found in %s", ips, v)
	}
}
