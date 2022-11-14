package ome

import (
	"context"
	"fmt"
	"strings"
	"terraform-provider-ome/clients"

	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ tfsdk.AttributeValidator = &validFqddsValidator{}
var _ tfsdk.AttributeValidator = &validTemplateViewTypeValidator{}

type validFqddsValidator struct {
}

// Description returns a plain text description of the validator's behavior, suitable for a practitioner to understand its impact.
func (v validFqddsValidator) Description(ctx context.Context) string {
	return fmt.Sprintf("Allowed values are  either all or one of these FQDDS: %s", clients.ValidFQDDS)
}

// MarkdownDescription returns a markdown formatted description of the validator's behavior, suitable for a practitioner to understand its impact.
func (v validFqddsValidator) MarkdownDescription(ctx context.Context) string {
	return fmt.Sprintf("Allowed values are  either all or one of these FQDDS: %s", clients.ValidFQDDS)
}

// Validate runs the main validation logic of the validator, reading configuration data out of `req` and updating `resp` with diagnostics.
func (v validFqddsValidator) Validate(ctx context.Context, req tfsdk.ValidateAttributeRequest, resp *tfsdk.ValidateAttributeResponse) {
	var fqdds types.String
	diags := tfsdk.ValueAs(ctx, req.AttributeConfig, &fqdds)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	if fqdds.Unknown || fqdds.Null {
		return
	}
	inputFqdds := fqdds.Value
	multipleFqdds := []string{}
	if strings.Contains(inputFqdds, ",") {
		multipleFqdds = strings.Split(inputFqdds, ",")
	}
	if len(multipleFqdds) == 0 {
		if strings.Contains(clients.ValidFQDDS, inputFqdds) {
			return
		}
		resp.Diagnostics.AddAttributeError(
			req.AttributePath,
			clients.ErrInvalidFqdds,
			fmt.Sprintf("FQDDS provided is not in this list of allowed values  %s", clients.ValidFQDDS),
		)
		return
	}
	for _, inpFqdds := range multipleFqdds {
		if strings.Contains(clients.ValidFQDDS, strings.TrimSpace(inpFqdds)) {
			continue
		} else {
			resp.Diagnostics.AddAttributeError(
				req.AttributePath,
				clients.ErrInvalidFqdds,
				fmt.Sprintf("FQDDS provided is not in this list of allowed values  :  %s", clients.ValidFQDDS),
			)
			break
		}
	}

}

type validTemplateViewTypeValidator struct {
}

// Description returns a plain text description of the validator's behavior, suitable for a practitioner to understand its impact.
func (v validTemplateViewTypeValidator) Description(ctx context.Context) string {
	return fmt.Sprintf("Allowed values are  either all or one of these FQDDS: %s", clients.ValidFQDDS)
}

// MarkdownDescription returns a markdown formatted description of the validator's behavior, suitable for a practitioner to understand its impact.
func (v validTemplateViewTypeValidator) MarkdownDescription(ctx context.Context) string {
	return fmt.Sprintf("Allowed values are  either all or one of these FQDDS: %s", clients.ValidFQDDS)
}

// Validate runs the main validation logic of the validator, reading configuration data out of `req` and updating `resp` with diagnostics.
func (v validTemplateViewTypeValidator) Validate(ctx context.Context, req tfsdk.ValidateAttributeRequest, resp *tfsdk.ValidateAttributeResponse) {
	var templateViewType types.String
	diags := tfsdk.ValueAs(ctx, req.AttributeConfig, &templateViewType)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	if templateViewType.Unknown || templateViewType.Null {
		return
	}

	validTemplateViewTypes := strings.Split(clients.ValidTemplateViewTypes, ",")
	for _, validTemplateViewType := range validTemplateViewTypes {
		if templateViewType.Value == validTemplateViewType {
			return
		}
	}
	resp.Diagnostics.AddAttributeError(
		req.AttributePath,
		clients.ErrInvalidTemplateViewType,
		fmt.Sprintf("Allowed values are either all or one of  :  %s", clients.ValidTemplateViewTypes),
	)
}
