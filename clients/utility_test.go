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
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDeviceMutuallyExclusive(t *testing.T) {
	type args struct {
		serviceTags []string
		devIDs      []int64
	}
	tests := []struct {
		name string
		args args
		want string
		err  string
	}{
		{"DeviceMutuallyExclusive Service Tags", args{[]string{"SVT1"}, nil}, ServiceTags, ""},
		{"DeviceMutuallyExclusive Device ids", args{nil, []int64{12}}, DeviceIDs, ""},
		{"DeviceMutuallyExclusive Error servicetags and deviceids", args{[]string{"SVT1"}, []int64{12}}, "", ErrDeviceMutuallyExclusive},
		{"DeviceMutuallyExclusive Error both inputs not specified", args{nil, nil}, "", ErrDeviceRequired},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := DeviceMutuallyExclusive(tt.args.serviceTags, tt.args.devIDs)
			if tt.err != "" {
				assert.NotNil(t, err)
				assert.Equal(t, got, tt.want)
				assert.ErrorContains(t, err, tt.err)
			} else {
				assert.NotNil(t, got)
				assert.Nil(t, err)
				assert.Equal(t, got, tt.want)
			}
		})
	}
}
