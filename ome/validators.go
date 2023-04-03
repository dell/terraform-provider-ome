package ome

import (
	"context"
	"fmt"
	"terraform-provider-ome/clients"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

// To be removed in the upcoming release and use framework

var _ validator.Set = sizeAtLeastValidator{}

// sizeAtLeastValidator validates that list contains at least min elements.
type sizeAtLeastValidator struct {
	min int
}

// Description describes the validation in plain text formatting.
func (v sizeAtLeastValidator) Description(_ context.Context) string {
	return fmt.Sprintf(clients.ErrBaseLineTargetsSize, v.min)
}

// MarkdownDescription describes the validation in Markdown formatting.
func (v sizeAtLeastValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

// Validate performs the validation.
func (v sizeAtLeastValidator) ValidateSet(ctx context.Context, req validator.SetRequest, resp *validator.SetResponse) {
	elems, ok := validateList(ctx, req, resp)
	if !ok {
		return
	}

	if len(elems) < v.min {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			v.Description(ctx),
			fmt.Sprintf("current size : %d", len(elems)),
		)
	}

}

func validateList(ctx context.Context, request validator.SetRequest, response *validator.SetResponse) ([]attr.Value, bool) {
	l := request.ConfigValue
	if l.IsUnknown() || l.IsNull() {
		return nil, false
	}
	return l.Elements(), true
}

// SizeAtLeast returns an AttributeValidator which ensures that any configured
// attribute value:
//
//   - Is a List.
//   - Contains at least min elements.
//
// Null (unconfigured) and unknown (known after apply) values are skipped.
func SizeAtLeast(min int) validator.Set {
	return sizeAtLeastValidator{
		min: min,
	}
}

var _ validator.String = complianceStateValidator{}

// sizeAtLeastValidator validates that list contains at least min elements.
type complianceStateValidator struct {
}

// Description describes the validation in plain text formatting.
func (v complianceStateValidator) Description(_ context.Context) string {
	return fmt.Sprintf(clients.ErrBaseLineComplianceStatus, clients.ValidComplainceStatus)
}

// MarkdownDescription describes the validation in Markdown formatting.
func (v complianceStateValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

// Validate performs the validation.
func (v complianceStateValidator) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	input := req.ConfigValue
	if input.IsUnknown() || input.IsNull() {
		return
	}
	if !(input.ValueString() == clients.ValidComplainceStatus) {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			v.Description(ctx),
			fmt.Sprintf("current value : %s", input.ValueString()),
		)
	}
}
