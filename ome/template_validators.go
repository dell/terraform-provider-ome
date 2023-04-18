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
