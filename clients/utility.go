package clients

import (
	"fmt"
)

// DeviceMutuallyExclusive checks if the service tag , device ids  are mutually exclusive
func DeviceMutuallyExclusive(serviceTags []string, devIDs []int64) (string, error) {
	var usedDeviceInput string
	if len(serviceTags) == 0 && len(devIDs) == 0 {
		return "", fmt.Errorf(ErrDeviceRequired)
	}

	if len(serviceTags) > 0 && len(devIDs) > 0 {
		return "", fmt.Errorf(ErrDeviceMutuallyExclusive)
	}
	if len(serviceTags) > 0 {
		usedDeviceInput = ServiceTags
	} else if len(devIDs) > 0 {
		usedDeviceInput = DeviceIDs
	}
	return usedDeviceInput, nil
}
