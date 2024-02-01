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

package ome

import (
	"context"
	"fmt"
	"strings"
	"terraform-provider-ome/clients"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

var _ validator.String = &validFqddsValidator{}

type validFqddsValidator struct {
}

// Description returns a plain text description of the validator's behavior, suitable for a practitioner to understand its impact.
func (v validFqddsValidator) Description(ctx context.Context) string {
	return fmt.Sprintf("Allowed values are : %s", clients.ValidFQDDS)
}

// MarkdownDescription returns a markdown formatted description of the validator's behavior, suitable for a practitioner to understand its impact.
func (v validFqddsValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

// Validate runs the main validation logic of the validator, reading configuration data out of `req` and updating `resp` with diagnostics.
func (v validFqddsValidator) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	fqdds := req.ConfigValue
	if fqdds.IsUnknown() || fqdds.IsNull() {
		return
	}
	inputFqdds := fqdds.ValueString()
	multipleInputFqdds := strings.Split(inputFqdds, ",")
	multipleValidFqdds := strings.Split(clients.ValidFQDDS, ",")
	isValid := false

	for _, inpFqdds := range multipleInputFqdds {
		inputFqddsVal := strings.TrimSpace(inpFqdds)
		isValid = false
		for _, validFqdds := range multipleValidFqdds {
			if strings.EqualFold(validFqdds, inputFqddsVal) {
				isValid = true
				break
			}
		}
		if !isValid {
			resp.Diagnostics.AddAttributeError(
				req.Path,
				clients.ErrInvalidFqdds,
				v.Description(ctx),
			)
			break
		}
	}

}
