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
	"terraform-provider-ome/models"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCertUpload(t *testing.T) {
	ts := createNewTLSServer(t)
	defer ts.Close()

	opts := initOptions(ts)

	c, _ := NewClient(opts)

	_, err1 := c.PostCert("aGVsbG8gdGhlcmUgdmFsaWQuCg==")
	assert.Nil(t, err1)

	_, err2 := c.PostCert("aGVsbG8gdGhlcmUK")
	assert.NotNil(t, err2)
}

func TestCSR(t *testing.T) {
	ts := createNewTLSServer(t)
	defer ts.Close()

	opts := initOptions(ts)

	c, _ := NewClient(opts)

	_, err1 := c.GetCSR(models.CSRConfig{
		DistinguishedName: "valid",
	})
	assert.Nil(t, err1)

	_, err2 := c.GetCSR(models.CSRConfig{
		DistinguishedName: "invalid",
	})
	assert.NotNil(t, err2)
}

func TestCertGet(t *testing.T) {
	ts := createNewTLSServer(t)
	defer ts.Close()

	opts := initOptions(ts)

	c, _ := NewClient(opts)

	cert, err1 := c.GetCert()
	assert.Nil(t, err1)
	assert.NotEmpty(t, cert)
}
