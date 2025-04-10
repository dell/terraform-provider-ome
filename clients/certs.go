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
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"terraform-provider-ome/models"
)

// GetCSR is used to get certificate signing request from OME
func (c *Client) GetCSR(input models.CSRConfig) (string, error) {
	b, _ := json.Marshal(input)
	response, err := c.Post(CSRGenAPI, nil, b)
	if err != nil {
		return "", err
	}

	resp := make(map[string]string)
	bodyData, getBodyError := c.GetBodyData(response.Body)
	if getBodyError != nil {
		return "", getBodyError
	}
	err = c.JSONUnMarshal(bodyData, &resp)
	if err != nil {
		return "", err
	}
	return resp["CertificateData"], nil
}

// PostCert is used to upload an application certificate to OME
func (c *Client) PostCert(base64Encoded string) (string, error) {
	decodedData, err := base64.StdEncoding.DecodeString(base64Encoded)
	if err != nil {
		return "", fmt.Errorf("failed to decode base64: %v", err)
	}

	headers := map[string]string{
		"Content-Type": "application/octet-stream",
		"Accept":       "application/octet-stream",
	}

	b := bytes.NewBuffer(decodedData)

	response, errp := c.PostFile(CertUploadAPI, headers, b)
	if errp != nil {
		return "", errp
	}
	respStr, getBodyError := c.GetBodyData(response.Body)
	if getBodyError != nil {
		return "", getBodyError
	}
	return string(respStr), nil
}

// GetCert is used to get application certificate info from OME
func (c *Client) GetCert() (models.CertInfo, error) {
	var ret models.CertInfo
	response, err := c.Get(CertGetAPI, nil, nil)
	if err != nil {
		return ret, err
	}
	resp, getBodyError := c.GetBodyData(response.Body)
	if getBodyError != nil {
		return ret, getBodyError
	}
	err = c.JSONUnMarshalSingleValue(resp, &ret)
	return ret, err
}
