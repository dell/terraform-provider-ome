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
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"strings"
	"time"
)

// Client type is to hold http client information
type Client struct {
	// httpclient from net/http
	httpclient *http.Client
	//url - base url of the form https://ipaddr:port, with no trailing slash
	url string
	//retry count for http request retries on timeout
	retry int
	//username - used to set the username for authentication
	username string
	//password - used to set the password for authentication
	password string
	//token - OAuth token received after authentication
	token string
	//sessionID - received after authentication
	sessionID string
	//PreRequestHook is the function to be invoked before making the http requests
	preRequestHook PreRequestHook
}

// PreRequestHook is the function to be invoked before making the http requests
type (
	PreRequestHook func(*Client, *http.Request)
)

// ClientOptions - different arguments to the New client implementation
type ClientOptions struct {
	// url - base url of the form https://ipaddr:port, with no trailing slash
	URL string
	// skipSsl - used to set InsecureSkipVerify for SSL validation
	SkipSSL bool
	// RootCaPath - path of the root ca
	RootCaPath string
	// Timeout - used to set timeout for http request
	Timeout time.Duration
	// Retry - used to set the number of retries to be done on timeout
	Retry int
	// Username - used to set the username for client
	Username string
	//password - used to set the password for client
	Password string
	// PreRequestHook - used to set the pre-request function.
	PreRequestHook PreRequestHook
}

// NewClient creates a https client by accepting ClientOptions as an argument
func NewClient(opts ClientOptions) (*Client, error) {
	omeClient := &Client{
		httpclient:     &http.Client{Timeout: opts.Timeout},
		url:            opts.URL,
		retry:          opts.Retry,
		username:       opts.Username,
		password:       opts.Password,
		preRequestHook: opts.PreRequestHook,
	}

	if opts.SkipSSL { // #nosec G402
		omeClient.httpclient.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: opts.SkipSSL,
			},
		}

	} else {
		pool, _ := x509.SystemCertPool() //return the system certificate pool
		if opts.RootCaPath != "" {
			rootCAsData, readErr := os.ReadFile(opts.RootCaPath)
			if readErr != nil {
				return nil, readErr
			}
			pool.AppendCertsFromPEM(rootCAsData)
		}
		// #nosec G402
		omeClient.httpclient.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: opts.SkipSSL,
				RootCAs:            pool,
			},
		}
	}
	return omeClient, nil
}

// GetURL return the ome url
func (c *Client) GetURL() string {
	return c.url
}

// GetHTTPClient returns the https client
func (c *Client) GetHTTPClient() *http.Client {
	return c.httpclient
}

// GetSessionID returns the sessionID
func (c *Client) GetSessionID() string {
	return c.sessionID
}

// GetSessionToken returns the auth token
func (c *Client) GetSessionToken() string {
	return c.token
}

// SetSessionID sets the sessionID
func (c *Client) SetSessionID(in string) {
	c.sessionID = in
}

// SetSessionToken sets the auth token
func (c *Client) SetSessionToken(in string) {
	c.token = in
}

// SetSessionParams sets the Session Params
func (c *Client) SetSessionParams(token, sessionID string) {
	c.SetSessionID(sessionID)
	c.SetSessionToken(token)
}

// Get sends an HTTP request using the GET method to the API.
func (c *Client) Get(
	path string,
	headers map[string]string,
	queryParams map[string]string) (*http.Response, error) {

	return c.Do(http.MethodGet, path, headers, queryParams, nil)
}

// Post sends an HTTP request using the POST method to the API.
func (c *Client) Post(
	path string,
	headers map[string]string,
	body []byte) (*http.Response, error) {

	return c.Do(http.MethodPost, path, headers, nil, body)
}

// Patch sends an HTTP request using the PATCH method to the API.
func (c *Client) Patch(
	path string,
	headers map[string]string,
	body []byte) (*http.Response, error) {

	return c.Do(http.MethodPatch, path, headers, nil, body)
}

// Put sends an HTTP request using the Put method to the API.
func (c *Client) Put(
	path string,
	headers map[string]string,
	body []byte) (*http.Response, error) {

	return c.Do(http.MethodPut, path, headers, nil, body)
}

// Delete sends an HTTP request using the Delete method to the API.
func (c *Client) Delete(
	path string,
	headers map[string]string,
	queryParams map[string]string) (*http.Response, error) {

	return c.Do(http.MethodDelete, path, headers, queryParams, nil)
}

// Do sends an HTTP request using the given method to the API.
func (c *Client) Do(
	method string,
	path string,
	headers map[string]string,
	queryParams map[string]string,
	body []byte) (*http.Response, error) {

	pathURL := c.url + path

	request, _ := http.NewRequest(method, pathURL, strings.NewReader(string(body)))

	//PrereqHook
	if c.preRequestHook != nil {
		c.preRequestHook(c, request)
	}
	//Add Request Header if any
	c.addHeaders(request, headers)

	//Add Request query params if any
	c.addQueryParams(request, queryParams)

	var response *http.Response
	var err error

	for attempt := 1; attempt <= c.retry; attempt++ {
		response, err = c.GetHTTPClient().Do(request)
		if err != nil {
			if e, ok := err.(net.Error); ok && e.Timeout() {
				time.Sleep(waitTime)
				err = fmt.Errorf(ErrRetryTimeoutMsg, attempt)
				response = nil
			} else {
				response = nil
				break
			}

		} else {
			break
		}
	}

	if response != nil && response.StatusCode != http.StatusOK && response.StatusCode != http.StatusAccepted &&
		response.StatusCode != http.StatusCreated && response.StatusCode != http.StatusNoContent {
		data, _ := c.GetBodyData(response.Body)
		return nil, fmt.Errorf(ErrResponseMsg, response.StatusCode, string(data))
	}

	return response, err
}

// addHeaders to add header to the request
func (c *Client) addHeaders(request *http.Request, headers map[string]string) {
	for k, value := range headers {
		request.Header.Add(k, value)
	}
}

// addQueryParams Adds query params to the request
func (c *Client) addQueryParams(request *http.Request, queryParams map[string]string) {
	q := request.URL.Query()
	for k, value := range queryParams {
		q.Add(k, value)
	}
	// request.URL.RawQuery = q.Encode()
	request.URL.RawQuery = strings.ReplaceAll(q.Encode(), "+", "%20")
}

// GetBodyData returns the body data
func (c *Client) GetBodyData(body io.ReadCloser) ([]byte, error) {

	if body == nil {
		return nil, errors.New(ErrEmptyBodyMsg)
	}

	data, _ := io.ReadAll(body)

	err := body.Close()
	if err != nil {
		return nil, err
	}
	return data, nil
}

// JSONMarshal - marshals the interface to bytes
func (c *Client) JSONMarshal(in interface{}) ([]byte, error) {
	return json.Marshal(in)
}

// JSONUnMarshal - unmarshals the byte to a interface
func (c *Client) JSONUnMarshal(data []byte, in interface{}) error {
	return json.Unmarshal(data, in)
}
