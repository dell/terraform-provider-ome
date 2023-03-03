package ome

import (
	"context"
	"fmt"
	"strings"
	"terraform-provider-ome/clients"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

type outputFormatValidator struct {
}

// Description returns a plain text description of the validator's behavior, suitable for a practitioner to understand its impact.
func (o outputFormatValidator) Description(ctx context.Context) string {
	return fmt.Sprintf("Allowed values are one of these output formats: %s", clients.ValidOutputFormat)
}

// MarkdownDescription returns a markdown formatted description of the validator's behavior, suitable for a practitioner to understand its impact.
func (o outputFormatValidator) MarkdownDescription(ctx context.Context) string {
	return fmt.Sprintf("Allowed values are one of these output formats: %s", clients.ValidOutputFormat)
}

// Validate runs the main validation logic of the validator, reading configuration data out of `req` and updating `resp` with diagnostics.
func (o outputFormatValidator) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	outputFormat := req.ConfigValue

	if outputFormat.IsUnknown() || outputFormat.IsNull() {
		return
	}

	validOutputFormatTypes := strings.Split(clients.ValidOutputFormat, ",")
	for _, validOutputFormat := range validOutputFormatTypes {
		if outputFormat.ValueString() == validOutputFormat {
			return
		}
	}
	resp.Diagnostics.AddAttributeError(
		req.Path,
		clients.ErrInvalidTemplateViewType,
		fmt.Sprintf("Allowed values are one of  :  %s", clients.ValidOutputFormat),
	)

}
