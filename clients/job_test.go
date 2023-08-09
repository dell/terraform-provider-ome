package clients

import (
	_ "embed"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNetwork_GetJobByID(t *testing.T) {
	ts := createNewTLSServer(t)
	defer ts.Close()

	opts := initOptions(ts)
	c, _ := NewClient(opts)

	tests := []struct {
		name string
		args int
	}{
		{"Get Job Successfully", 13881},
		{"Get Job Failed", -1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			getJob, err := c.GetJobByID(tt.args)
			t.Log(getJob, err)
			if err == nil {
				assert.Equal(t, tt.args, getJob.ID)
			}
		})
	}
}
