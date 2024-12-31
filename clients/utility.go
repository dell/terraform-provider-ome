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

package clients

import (
	"fmt"
)

// DeviceMutuallyExclusive checks if the service tag , device ids  are mutually exclusive
func DeviceMutuallyExclusive(serviceTags []string, devIDs []int64) (string, error) {
	var usedDeviceInput string
	if len(serviceTags) == 0 && len(devIDs) == 0 {
		return "", fmt.Errorf("%s", ErrDeviceRequired)
	}

	if len(serviceTags) > 0 && len(devIDs) > 0 {
		return "", fmt.Errorf("%s", ErrDeviceMutuallyExclusive)
	}
	if len(serviceTags) > 0 {
		usedDeviceInput = ServiceTags
	} else if len(devIDs) > 0 {
		usedDeviceInput = DeviceIDs
	}
	return usedDeviceInput, nil
}

// CompareInt64 compares the two array and reurns the diff
func CompareInt64(comparing, comparedTo []int64) []int64 {
	compareToMap := make(map[int64]int64)
	for _, val := range comparedTo {
		compareToMap[val]++
	}

	var diff []int64
	for _, val := range comparing {
		if compareToMap[val] > 0 {
			compareToMap[val]--
			continue
		}
		diff = append(diff, val)
	}
	return diff
}

// CompareString compares the two array and returns the diff
func CompareString(comparing, comparedTo []string) []string {
	compareToMap := make(map[string]int64)
	for _, val := range comparedTo {
		compareToMap[val]++
	}

	var diff []string
	for _, val := range comparing {
		if compareToMap[val] > 0 {
			compareToMap[val]--
			continue
		}
		diff = append(diff, val)
	}
	return diff
}

// FindElementInIntArray finds the element in an array
func FindElementInIntArray(arr []int64, find int64) int {
	index := -1
	for i := 0; i < len(arr); i++ {
		if arr[i] == find {
			index = i
		}
	}
	return index
}
