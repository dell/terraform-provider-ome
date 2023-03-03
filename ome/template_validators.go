package ome

import (
	"context"
	"fmt"
	"strings"
	"terraform-provider-ome/clients"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

var _ validator.String = &validFqddsValidator{}
var _ validator.String = &validTemplateViewTypeValidator{}
var _ validator.String = &validTemplateDeviceTypeValidator{}

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

type validTemplateViewTypeValidator struct {
}

// Description returns a plain text description of the validator's behavior, suitable for a practitioner to understand its impact.
func (v validTemplateViewTypeValidator) Description(ctx context.Context) string {
	return fmt.Sprintf("Allowed values are  :  %s ", clients.ValidTemplateViewTypes)
}

// MarkdownDescription returns a markdown formatted description of the validator's behavior, suitable for a practitioner to understand its impact.
func (v validTemplateViewTypeValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

// Validate runs the main validation logic of the validator, reading configuration data out of `req` and updating `resp` with diagnostics.
func (v validTemplateViewTypeValidator) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	templateViewType := req.ConfigValue
	if templateViewType.IsUnknown() || templateViewType.IsNull() {
		return
	}

	validTemplateViewTypes := strings.Split(clients.ValidTemplateViewTypes, ",")
	for _, validTemplateViewType := range validTemplateViewTypes {
		if strings.EqualFold(strings.TrimSpace(templateViewType.ValueString()), validTemplateViewType) {
			return
		}
	}
	resp.Diagnostics.AddAttributeError(
		req.Path,
		clients.ErrInvalidTemplateViewType,
		v.Description(ctx),
	)
}

type validTemplateDeviceTypeValidator struct {
}

// Description returns a plain text description of the validator's behavior, suitable for a practitioner to understand its impact.
func (v validTemplateDeviceTypeValidator) Description(ctx context.Context) string {
	return fmt.Sprintf("Allowed values are  :  %s ", clients.ValidTemplateDeviceTypes)
}

// MarkdownDescription returns a markdown formatted description of the validator's behavior, suitable for a practitioner to understand its impact.
func (v validTemplateDeviceTypeValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

// Validate runs the main validation logic of the validator, reading configuration data out of `req` and updating `resp` with diagnostics.
func (v validTemplateDeviceTypeValidator) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	templateDeviceType := req.ConfigValue
	if templateDeviceType.IsUnknown() || templateDeviceType.IsNull() {
		return
	}

	validTemplateDeviceTypes := strings.Split(clients.ValidTemplateDeviceTypes, ",")
	for _, validTemplateDeviceType := range validTemplateDeviceTypes {
		if strings.EqualFold(strings.TrimSpace(templateDeviceType.ValueString()), validTemplateDeviceType) {
			return
		}
	}
	resp.Diagnostics.AddAttributeError(
		req.Path,
		clients.ErrInvalidTemplateViewType,
		v.Description(ctx),
	)
}
