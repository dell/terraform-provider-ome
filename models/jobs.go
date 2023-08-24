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

package models

import (
	"encoding/json"
)

type JobPayload struct {
	ID             int64     `json:"Id"`
	Enabled        StateType `json:"State"`
	JobName        string
	JobDescription string `json:"JobDescription,omitempty"`
	Schedule       string
	JobType        JobType
	Params         JobParams `json:"Params,omitempty"`
	Targets        []JobTargetType
}

type StateType bool

// MarshalJSON - implements marshaller interface
func (s StateType) MarshalJSON() ([]byte, error) {
	if s {
		return json.Marshal("Enabled")
	}
	return json.Marshal("Disabled")
}

type JobParams map[string]string

// MarshalJSON - implements marshaller interface
func (j JobParams) MarshalJSON() ([]byte, error) {
	type item struct {
		Key   string
		Value string
	}
	k := make([]item, 0)
	for key, value := range j {
		k = append(k, item{key, value})
	}
	return json.Marshal(&k)
}

type JobTargetType struct {
	ID         int64 `json:"Id"`
	Data       string
	TargetType TargetType
}

type TargetType uint8

const (
	DeviceTargetType TargetType = iota
)

// MarshalJSON - implements marshaller interface
func (TargetType) MarshalJSON() ([]byte, error) {
	return json.Marshal(&map[string]any{
		"Id":   8,
		"Name": "Inventory_Task",
	})
}

type JobType uint8

const (
	InventoryRefreshJobType JobType = iota
	ResetIDRACJobType
	ClearJobQueueJobType
)

// MarshalJSON - implements marshaller interface
func (j JobType) MarshalJSON() ([]byte, error) {
	job_type_map := map[JobType]uint8{InventoryRefreshJobType: 8, ResetIDRACJobType: 3, ClearJobQueueJobType: 3}
	jtype_map := map[uint8]string{3: "DeviceAction_Task", 8: "Inventory_Task"}

	return json.Marshal(&struct {
		ID   uint8  `json:"Id"`
		Name string `json:"Name"`
	}{
		ID:   job_type_map[j],
		Name: jtype_map[job_type_map[j]],
	})
}
