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
