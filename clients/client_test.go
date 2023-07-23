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
	"net/http/httptest"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func testPreReq(c *Client, r *http.Request) {
	r.Header.Add("x-test-header", "test-value")
	r.Header.Add("Content-Type", "application/json")
	r.Header.Add(AuthTokenHeader, c.GetSessionToken())
}

// initOptions internal impl to get init data
func initOptions(ts *httptest.Server) ClientOptions {
	opts := ClientOptions{
		URL:        ts.URL,
		SkipSSL:    true,
		RootCaPath: "",
		Timeout:    time.Second * 30,
		Retry:      1,
		Username:   "admin",
		Password:   "Password123!",
	}
	return opts
}

// getTestData to get test data Directory
func getTestData(fileName string) string {
	wd, _ := os.Getwd()
	parent := filepath.Dir(wd)
	rootCAs := filepath.Join(parent, "testdata", fileName)
	return rootCAs
}

// TestBodyData tests basic GET request.
func TestBodyData(t *testing.T) {
	ts := createNewTLSServer(t)
	defer ts.Close()

	opts := initOptions(ts)

	c, _ := NewClient(opts)

	body, err := c.GetBodyData(nil)

	assert.NotNil(t, err)
	assert.Nil(t, body)
	assert.ErrorContains(t, err, ErrEmptyBodyMsg)

	response, _ := c.Get("/emptyBody", nil, nil)
	body, _ = c.GetBodyData(response.Body)
	assert.Equal(t, []byte{}, body)

	response, _ = c.Post("/data", nil, nil)
	body, _ = c.GetBodyData(response.Body)
	//assert response body
	assert.Equal(t, []byte(`Hello from TLS server post body`), body)

	response, _ = c.Patch("/data", nil, nil)
	body, _ = c.GetBodyData(response.Body)
	//assert response body
	assert.Equal(t, []byte(`Hello from TLS server`), body)

	response, _ = c.Put("/data", nil, nil)
	body, _ = c.GetBodyData(response.Body)
	//assert response body
	assert.Equal(t, []byte(`Hello from TLS server`), body)

	response, _ = c.Delete("/data", nil, nil)
	body, _ = c.GetBodyData(response.Body)
	//assert response body
	assert.Equal(t, []byte{}, body)
}

// TestGetHttpClient tests the GetHttpClient method
func TestGetHttpClient(t *testing.T) {
	ts := createNewTLSServer(t)
	defer ts.Close()

	opts := initOptions(ts)
	c, _ := NewClient(opts)
	client := c.GetHTTPClient()

	assert.Equal(t, c.httpclient, client)
}

// TestClientVerifyTimeout verified the timeout value set
func TestClientVerifyTimeout(t *testing.T) {
	ts := createNewTLSServer(t)
	defer ts.Close()

	opts := initOptions(ts)

	c, _ := NewClient(opts)
	assert.Equal(t, c.httpclient.Timeout.Seconds(), 30.0)
}

// TestClientVerifyUrl verifies the url set scheme and host
func TestClientVerifyUrl(t *testing.T) {
	ts := createNewTLSServer(t)
	defer ts.Close()

	opts := initOptions(ts)

	c, _ := NewClient(opts)
	assert.Equal(t, ts.URL, c.GetURL())
}

// TestClientVerifyRetry verifies the retry set
func TestClientVerifyRetry(t *testing.T) {
	ts := createNewTLSServer(t)
	defer ts.Close()

	opts := initOptions(ts)

	c, _ := NewClient(opts)
	assert.Equal(t, opts.Retry, c.retry)
}

// TestClientVerifyUserNameAndPassword verifies the retry set
func TestClientVerifyUserNameAndPassword(t *testing.T) {
	ts := createNewTLSServer(t)
	defer ts.Close()

	opts := initOptions(ts)

	c, _ := NewClient(opts)
	assert.Equal(t, opts.Username, c.username)
	assert.Equal(t, opts.Password, c.password)
}

// TestClientVerifyPreRequestHook verifies the preRequestHook set
func TestClientVerifyPreRequestHook(t *testing.T) {
	ts := createNewTLSServer(t)
	defer ts.Close()

	opts := initOptions(ts)
	opts.PreRequestHook = testPreReq
	c, _ := NewClient(opts)

	assert.NotNil(t, c.preRequestHook)
}

// TestClientSslCertificate Verifies SSL certificate scenarios
func TestClientSslCertificate(t *testing.T) {
	ts := createNewTLSServer(t)
	defer ts.Close()

	var tests = []ClientOptions{
		{"https://127.0.0.1:8234", true, "", time.Second * 30, 1, "", "", nil},
		{"https://127.0.0.1:8234", true, getTestData("sample_ca.pem"), time.Second * 30, 1, "", "", nil},
		{"https://127.0.0.1:8234", true, getTestData("sample_ca_invalid.pem"), time.Second * 30, 1, "", "", nil},

		{"https://127.0.0.1:8234", false, "", time.Second * 30, 1, "", "", nil},
		{"https://127.0.0.1:8234", false, getTestData("sample_ca.pem"), time.Second * 30, 1, "", "", nil},
		{"https://127.0.0.1:8234", false, getTestData("sample_ca_invalid.pem"), time.Second * 30, 1, "", "", nil},
	}
	for _, v := range tests {
		c, err := NewClient(v)

		if v.SkipSSL {
			assert.Equal(t, v.SkipSSL, c.httpclient.Transport.(*http.Transport).TLSClientConfig.InsecureSkipVerify)
			assert.Nil(t, c.httpclient.Transport.(*http.Transport).TLSClientConfig.RootCAs)
		} else {
			if strings.Contains(v.RootCaPath, "sample_ca_invalid.pem") {
				assert.NotNil(t, err) // message can be different based on OS
			} else {
				assert.NotNil(t, c.httpclient.Transport.(*http.Transport).TLSClientConfig.RootCAs)
			}
		}
	}
}

func TestDo(t *testing.T) {
	ts := createNewTLSServer(t)
	defer ts.Close()

	opts := initOptions(ts)

	c, _ := NewClient(opts)

	headers := make(map[string]string)
	headers["Content-Type"] = "application/json"
	headers["X-Auth-Token"] = "wwrwetwyhbdsvsdhkoqndfsjhgrutcasfbdfhgch"

	queryParams := make(map[string]string)
	queryParams["query1"] = "value1"
	queryParams["query2"] = "value2"

	clientBody := []byte(`Hello from Client post body`)

	type args struct {
		method      string
		path        string
		headers     map[string]string
		queryParams map[string]string
		body        []byte
	}
	tests := []struct {
		name string
		c    *Client
		args args
	}{
		{"Test GET method No Error and Response is Not Nil", c, args{"GET", "/data", nil, nil, nil}},
		{"Test GET method with Header", c, args{"GET", "/data", map[string]string{"Content-Type": "application/json"}, nil, nil}},
		{"Test GET method with multiple Headers", c, args{"GET", "/data", headers, nil, nil}},
		{"Test GET method with Queryparams", c, args{"GET", "/data", nil, map[string]string{"query1": "value1"}, nil}},
		{"Test GET method with multiple Queryparams", c, args{"GET", "/data", nil, queryParams, nil}},
		{"Test GET method with Multiple Header and Queryparams", c, args{"GET", "/data", headers, queryParams, nil}},

		{"Test POST method No Error and Response is Not Nil", c, args{"POST", "/data", nil, nil, nil}},
		{"Test POST method No Error and Response is Not Nil", c, args{"POST", "/data", nil, nil, clientBody}},
		{"Test POST method with Header", c, args{"POST", "/data", map[string]string{"Content-Type": "application/json"}, nil, nil}},
		{"Test POST method with multiple Headers", c, args{"POST", "/data", headers, nil, nil}},
		{"Test POST method with Queryparams", c, args{"POST", "/data", nil, map[string]string{"query1": "value1"}, nil}},
		{"Test POST method with multiple Queryparams", c, args{"POST", "/data", nil, queryParams, nil}},
		{"Test POST method with Multiple Header and Queryparams", c, args{"POST", "/data", headers, queryParams, nil}},

		{"Test PATCH method No Error and Response is Not Nil", c, args{"PATCH", "/data", nil, nil, nil}},
		{"Test PATCH method No Error and Response is Not Nil", c, args{"PATCH", "/data", nil, nil, clientBody}},
		{"Test PATCH method with Header", c, args{"PATCH", "/data", map[string]string{"Content-Type": "application/json"}, nil, nil}},
		{"Test PATCH method with multiple Headers", c, args{"PATCH", "/data", headers, nil, nil}},
		{"Test PATCH method with Queryparams", c, args{"PATCH", "/data", nil, map[string]string{"query1": "value1"}, nil}},
		{"Test PATCH method with multiple Queryparams", c, args{"PATCH", "/data", nil, queryParams, nil}},
		{"Test PATCH method with Multiple Header and Queryparams", c, args{"PATCH", "/data", headers, queryParams, nil}},

		{"Test PUT method No Error and Response is Not Nil", c, args{"PUT", "/data", nil, nil, nil}},
		{"Test PUT method No Error and Response is Not Nil", c, args{"PUT", "/data", nil, nil, clientBody}},
		{"Test PUT method with Header", c, args{"PUT", "/data", map[string]string{"Content-Type": "application/json"}, nil, nil}},
		{"Test PUT method with multiple Headers", c, args{"PUT", "/data", headers, nil, nil}},
		{"Test PUT method with Queryparams", c, args{"PUT", "/data", nil, map[string]string{"query1": "value1"}, nil}},
		{"Test PUT method with multiple Queryparams", c, args{"PUT", "/data", nil, queryParams, nil}},
		{"Test PUT method with Multiple Header and Queryparams", c, args{"PUT", "/data", headers, queryParams, nil}},

		{"Test DELETE method No Error and Response is Not Nil", c, args{"DELETE", "/data", nil, nil, nil}},
		{"Test DELETE method with Header", c, args{"DELETE", "/data", map[string]string{"Content-Type": "application/json"}, nil, nil}},
		{"Test DELETE method with multiple Headers", c, args{"DELETE", "/data", headers, nil, nil}},
		{"Test DELETE method with Queryparams", c, args{"DELETE", "/data", nil, map[string]string{"query1": "value1"}, nil}},
		{"Test DELETE method with multiple Queryparams", c, args{"DELETE", "/data", nil, queryParams, nil}},
		{"Test DELETE method with Multiple Header and Queryparams", c, args{"DELETE", "/data", headers, queryParams, nil}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			response, err := tt.c.Do(tt.args.method, tt.args.path, tt.args.headers, tt.args.queryParams, tt.args.body)

			assert.Nil(t, err)
			assert.NotNil(t, response)

			//assert response body and status
			body, _ := c.GetBodyData(response.Body)

			if tt.args.method == "GET" || tt.args.method == "PATCH" || tt.args.method == "PUT" {
				assert.Equal(t, http.StatusOK, response.StatusCode)
				// GET response body assertion
				assert.Equal(t, []byte(`Hello from TLS server`), body)
			}
			if tt.args.method == "POST" {
				assert.Equal(t, http.StatusCreated, response.StatusCode)
				// POST response body assertion
				assert.Equal(t, []byte(`Hello from TLS server post body`), body)
			}
			if tt.args.method == "DELETE" {
				assert.Equal(t, http.StatusNoContent, response.StatusCode)
				// GET response body assertion
				assert.Equal(t, []byte{}, body)
			}
			// Response assertion for content type
			assert.Equal(t, "application/json", response.Header.Get("Content-Type"))
			//assert Headers
			if tt.args.headers != nil {
				for k, v := range tt.args.headers {
					//assert header map for request
					assert.Equal(t, v, response.Request.Header.Get(k))

				}
			}
			//assert Query params
			if tt.args.queryParams != nil {
				//assert query params map for request
				for k, v := range tt.args.queryParams {
					assert.Equal(t, v, response.Request.URL.Query().Get(k))
				}
			}
		})
	}
}

func TestDoError(t *testing.T) {
	ts := createNewTLSServer(t)
	defer ts.Close()

	opts := initOptions(ts)

	opts.URL = "https://invalid.domain:8088"

	c, _ := NewClient(opts)

	type args struct {
		method      string
		path        string
		headers     map[string]string
		queryParams map[string]string
		body        []byte
	}
	tests := []struct {
		name string
		c    *Client
		args args
	}{
		{"Test GET method for Error and Response is Nil", c, args{"GET", "/data", nil, nil, nil}},
		{"Test Post method for Error and Response is Nil", c, args{"POST", "/data", nil, nil, nil}},
		{"Test Patch method for Error and Response is Nil", c, args{"PATCH", "/data", nil, nil, nil}},
		{"Test Put method for Error and Response is Nil", c, args{"PUT", "/data", nil, nil, nil}},
		{"Test Delete method for Error and Response is Nil", c, args{"DELETE", "/data", nil, nil, nil}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			response, err := tt.c.Do(tt.args.method, tt.args.path, tt.args.headers, tt.args.queryParams, tt.args.body)
			//Assert that err is not nill
			assert.NotNil(t, err)
			//Assert that response is not nill
			assert.Nil(t, response)

		})
	}
}

// TestDoRetry tests the timeout returned by the client
func TestDoRetry(t *testing.T) {

	ts := createNewTLSServer(t)
	defer ts.Close()

	opts := initOptions(ts)

	type args struct {
		method      string
		path        string
		headers     map[string]string
		queryParams map[string]string
		body        []byte
	}
	tests := []struct {
		name    string
		args    args
		timeout time.Duration
		retry   int
	}{
		{"Test timeout value with 1", args{"GET", "/timeout", nil, nil, nil}, 1 * time.Second, 1},
		{"Test timeout value with 2", args{"GET", "/timeout", nil, nil, nil}, 1 * time.Second, 2},
		{"Test timeout value with 3", args{"GET", "/timeout", nil, nil, nil}, 1 * time.Second, 3},
		{"Test timeout value Fails on 1st and 2nd attempt and successfull on 3rd attempt", args{"GET", "/timeout-success", nil, nil, nil}, 3 * time.Second, 4},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			opts.Timeout = tt.timeout
			opts.Retry = tt.retry
			c, _ := NewClient(opts)
			response, err := c.Do(tt.args.method, tt.args.path, tt.args.headers, tt.args.queryParams, tt.args.body)
			if opts.Retry <= 3 { // just a condition to access
				assert.NotNil(t, err)
				assert.Nil(t, response)
				assert.ErrorContains(t, err, fmt.Sprintf(ErrRetryTimeoutMsg, tt.retry))
			} else {
				assert.NotNil(t, response)
				assert.Nil(t, err)
			}
		})
	}
}

// TestDoPreReqHook
func TestDoPreReqHook(t *testing.T) {
	ts := createNewTLSServer(t)
	defer ts.Close()

	opts := initOptions(ts)
	opts.PreRequestHook = testPreReq
	c, _ := NewClient(opts)

	response, _ := c.Get("/test", nil, nil)

	assert.Equal(t, "test-value", response.Request.Header.Get("x-test-header"))

}

// TestClientJsonMarshal
func TestClientJsonMarshal(t *testing.T) {
	ts := createNewTLSServer(t)
	defer ts.Close()

	opts := initOptions(ts)

	c, _ := NewClient(opts)
	type SampleTest struct {
		Name        string `json:"Name"`
		DisplayName string `json:"DisplayName"`
		ID          int    `json:"Id"`
	}

	st := SampleTest{
		Name:        "YourName",
		DisplayName: "Your Name",
		ID:          123,
	}

	type args struct {
		in interface{}
	}
	tests := []struct {
		name string
		args args
	}{
		{"", args{st}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := c.JSONMarshal(tt.args.in)
			assert.NotNil(t, got)
			assert.Nil(t, err)
			assert.Equal(t, []byte(`{"Name":"YourName","DisplayName":"Your Name","Id":123}`), got)

		})
	}
}

// TestClientJsonMarshal
func TestClientJsonUnMarshal(t *testing.T) {
	ts := createNewTLSServer(t)
	defer ts.Close()

	opts := initOptions(ts)

	c, _ := NewClient(opts)
	type SampleTest struct {
		Name        string `json:"Name"`
		DisplayName string `json:"DisplayName"`
		ID          int    `json:"Id"`
	}

	st := SampleTest{}

	type args struct {
		in   interface{}
		data []byte
	}
	tests := []struct {
		name string
		args args
	}{
		{"", args{&st, []byte(`{"Name":"YourName","DisplayName":"Your Name","Id":123}`)}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := c.JSONUnMarshal(tt.args.data, tt.args.in)
			assert.Nil(t, err)
			assert.Equal(t, 123, st.ID)
			assert.Equal(t, "YourName", st.Name)
			assert.Equal(t, "Your Name", st.DisplayName)
		})
	}
}

func TestUnmarshalValue(t *testing.T) {
	c := Client{}

	t.Run("not_a_json_neg", func(t *testing.T) {
		in, data := "&st", []byte(`not a json`)
		err := c.JSONUnMarshalValue(data, &in)
		assert.NotEqual(t, nil, err)
	})

	t.Run("no_value_field_neg", func(t *testing.T) {
		in, data := "&st", []byte(`{"name": "is a json", "description": "no value field"}`)
		err := c.JSONUnMarshalValue(data, &in)
		assert.NotEqual(t, nil, err)
	})

	t.Run("has_value_field", func(t *testing.T) {
		in, data, out := "&st", []byte(`{"name": "is a json", "value": "here is the value field"}`), "here is the value field"
		err := c.JSONUnMarshalValue(data, &in)
		assert.Nil(t, err)
		assert.Truef(t, reflect.DeepEqual(in, out), "Expected:%v Got:%v", out, in)
	})

	t.Run("value_type_mismatch_neg", func(t *testing.T) {
		in, data := make([]int, 0), []byte(`{"name": "is a json", "value": ["this","field", "is", "not","a", "list","of","ints"]`)
		err := c.JSONUnMarshalValue(data, &in)
		assert.NotEqual(t, nil, err)
	})

	t.Run("value_type_match", func(t *testing.T) {
		in, data, out := make([]string, 0), []byte(`{"name": "is a json", "value": ["a", "b", "c"]}`), []string{"a", "b", "c"}
		err := c.JSONUnMarshalValue(data, &in)
		assert.Nil(t, err)
		assert.Truef(t, reflect.DeepEqual(in, out), "Expected:%v Got:%v", out, in)
	})
}

func TestUnmarshalSingleValue(t *testing.T) {
	c := Client{}
	t.Run("no_value_field_neg", func(t *testing.T) {
		in, data := "dummy", []byte(`{"name": "is a json", "description": "no value field"}`)
		err := c.JSONUnMarshalSingleValue(data, &in)
		assert.NotEqual(t, nil, err)
	})

	t.Run("value_is_not_list_neg", func(t *testing.T) {
		in, data := "dummy", []byte(`{"name": "is a json", "value": "not a list"}`)
		err := c.JSONUnMarshalSingleValue(data, &in)
		assert.NotEqual(t, nil, err)
	})

	t.Run("multiple_items_neg", func(t *testing.T) {
		in, data := "dummy", []byte(`{"name": "multiple items", "value": [1, 2, 3]}`)
		err := c.JSONUnMarshalSingleValue(data, &in)
		assert.NotEqual(t, nil, err)
	})

	t.Run("zero_items_neg", func(t *testing.T) {
		in, data := "dummy", []byte(`{"name": "zero items", "value": []}`)
		err := c.JSONUnMarshalSingleValue(data, &in)
		assert.NotEqual(t, nil, err)
	})

	t.Run("item_type_mismatch_neg", func(t *testing.T) {
		var in int
		data := []byte(`{"name": "string", "value": ["this item is not an int"]}`)
		err := c.JSONUnMarshalSingleValue(data, &in)
		assert.NotEqual(t, nil, err)
	})

	t.Run("item_type_match", func(t *testing.T) {
		var in int
		data, out := []byte(`{"name": "int", "value": [100]}`), 100
		err := c.JSONUnMarshalSingleValue(data, &in)
		assert.Nil(t, err)
		assert.Equal(t, out, in)
	})
}
