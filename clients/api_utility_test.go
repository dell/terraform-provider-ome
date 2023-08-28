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

package clients

import (
	"fmt"
	"net/http"
	"terraform-provider-ome/models"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestClienCreateSession - creates a session test
func TestClienCreateSession(t *testing.T) {

	ts := createNewTLSServer(t)
	defer ts.Close()

	tests := []struct {
		name     string
		id       int
		username string
		password string
	}{
		{"Create a session", 1, "admin", "Password123!"},
		{"Create a session with invalid username - unauthorized", 2, "myuser", "Password123!"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			opts := initOptions(ts)
			opts.Username = tt.username
			opts.Password = tt.password

			c, _ := NewClient(opts)

			resp, err := c.CreateSession()
			if tt.id == 1 {
				assert.Nil(t, err)
				assert.Equal(t, "13bc3f63-9376-44dc-a09f-3a94591a7c5d", c.GetSessionToken())
				assert.Equal(t, "e1817fe6-97e5-4ea0-88a9-d865c73021529", c.GetSessionID())
			} else {
				assert.NotNil(t, err)
				assert.Nil(t, resp)
			}
		})
	}
}

// TestClienCreateSession - creates a session test
func TestClienRemoveSession(t *testing.T) {

	ts := createNewTLSServer(t)
	defer ts.Close()

	opts := initOptions(ts)

	c, _ := NewClient(opts)

	tests := []struct {
		name      string
		id        int
		sessionID string
	}{
		{"Remove a session", 1, "e1817fe6-97e5-4ea0-88a9-d865c73021529"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c.SetSessionID(tt.sessionID) //Ideally done by the createSession
			resp, err := c.RemoveSession()
			assert.Nil(t, err)
			assert.NotNil(t, resp)
			assert.Equal(t, http.StatusNoContent, resp.StatusCode)
			assert.Equal(t, "", c.GetSessionToken())
			assert.Equal(t, "", c.GetSessionID())
		})
	}
}

func TestClient_TrackJob(t *testing.T) {
	ts := createNewTLSServer(t)
	defer ts.Close()

	opts := initOptions(ts)

	c, _ := NewClient(opts)

	type args struct {
		jobID         int64
		maxRetries    int64
		sleepInterval int64
	}
	tests := []struct {
		name string
		args args
	}{
		{"Track Job To Return Success On First Try", args{12345, 5, 1}},
		{"Track Job To Return Failure On First Try", args{23456, 5, 1}},
		{"Track Job To Return Failure with errors On First Try", args{34567, 5, 1}},
		{"Track Job To Return Running On First Try, Success on Second Try", args{45678, 5, 1}},
		{"Track Job To Return Running on all tries", args{56789, 5, 1}},
		{"Track Job with invalid job ID", args{13456, 5, 1}},
		{"Track Job with valid job ID with no last execution details", args{14567, 5, 1}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, message := c.TrackJob(tt.args.jobID, tt.args.maxRetries, tt.args.sleepInterval)
			if tt.args.jobID == 12345 || tt.args.jobID == 45678 {
				assert.Equal(t, true, got)
				assert.Equal(t, SuccessMsg, message)
			} else if tt.args.jobID == 23456 {
				assert.Equal(t, false, got)
				assert.Equal(t, "LastExecutionDetail Failure", message)
			} else if tt.args.jobID == 34567 {
				assert.Equal(t, false, got)
				assert.Equal(t, "LastExecutionDetail Warning", message)
			} else if tt.args.jobID == 56789 {
				assert.Equal(t, false, got)
				assert.Equal(t, fmt.Sprintf(JobIncompleteMsg, tt.args.jobID, tt.args.maxRetries), message)
			} else if tt.args.jobID == 13456 {
				assert.Equal(t, false, got)
				assert.Contains(t, message, "status: 400")
			} else if tt.args.jobID == 14567 {
				assert.Equal(t, false, got)
				assert.Contains(t, message, "No recent execution details were found for the provided job id.")
			}

		})
	}
}

func TestGetURL(t *testing.T) {
	host := "localhost"
	port := int64(443)

	actualURL := GetURL(host, port)

	assert.Equal(t, fmt.Sprintf("https://%s:%d", host, port), actualURL)
}

func TestClientPreReqHook(t *testing.T) {
	ts := createNewTLSServer(t)
	defer ts.Close()

	opts := initOptions(ts)

	c, _ := NewClient(opts)
	request, _ := http.NewRequest(http.MethodGet, c.url, nil)

	tests := []struct {
		name         string
		sessionToken string
	}{
		{"Empty Session Token", ""},
		{"With Session Token", "session_token"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c.SetSessionToken(tt.sessionToken)
			ClientPreReqHook(c, request)

			assert.Equal(t, "application/json", request.Header.Get("Content-Type"))
			assert.Equal(t, tt.sessionToken, request.Header.Get(AuthTokenHeader))
		})
	}
	ClientPreReqHook(c, request)
}

func TestClient_GetPaginatedData(t *testing.T) {

	ts := createNewTLSServer(t)
	defer ts.Close()

	opts := initOptions(ts)

	c, _ := NewClient(opts)

	type args struct {
		url string
		in  []models.Device
	}
	tests := []struct {
		name     string
		args     args
		wantData []models.Device
		wantErr  bool
	}{
		{"Test", args{fmt.Sprintf(GroupServiceDevicesAPI, 1013), []models.Device{}}, []models.Device{
			{ID: 10337},
			{ID: 10338},
		}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := c.GetPaginatedData(tt.args.url, &tt.args.in)
			if tt.wantErr {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
				for i, d := range tt.wantData {
					assert.Equal(t, d.ID, tt.args.in[i].ID)
				}

			}
		})
	}
}
func TestGetJob(t *testing.T) {
	ts := createNewTLSServer(t)
	defer ts.Close()

	opts := initOptions(ts)

	c, _ := NewClient(opts)

	tests := []struct {
		name  string
		id    int64
		jr    JobResp
		isErr bool
	}{
		{"Get Job respone success", 1, JobResp{ID: 1, LastRunStatus: JobStatus{ID: 2060}}, false},
		{"Get Job respone failure", 2, JobResp{ID: 1, LastRunStatus: JobStatus{ID: 2070}}, false},
		{"Get Job respone warning", 3, JobResp{ID: 1, LastRunStatus: JobStatus{ID: 2090}}, false},
		{"Get Job respone error", 4, JobResp{ID: 1, LastRunStatus: JobStatus{ID: 2016}}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := c.GetJob(tt.id)
			if tt.isErr {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
				assert.NotNil(t, resp)
				assert.Equal(t, tt.jr.LastRunStatus.ID, resp.LastRunStatus.ID)
			}
		})
	}
}

func TestJobOptsGetSchedule(t *testing.T) {
	tests := []struct {
		name   string
		input  JobOpts
		output string
	}{
		{"OnlyRunNow", JobOpts{"", "", true, ""}, RunNowSchedule},
		{"Both", JobOpts{"", "", true, "something"}, RunNowSchedule},
		{"None", JobOpts{"", "", false, ""}, RunNowSchedule},
		{"OnlySchedule", JobOpts{"", "", false, "something"}, "something"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.output, tt.input.getSchedule())
		})
	}
}
