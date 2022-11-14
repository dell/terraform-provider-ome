package models

// Devices - list of device response from on OME
type Devices struct {
	Value    []Device `json:"value"`
	NextLink string   `json:"@odata.nextLink"`
}

// Device - embedded device response from the Devices
type Device struct {
	ID               int64  `json:"Id"`
	DeviceServiceTag string `json:"DeviceServiceTag"`
}
