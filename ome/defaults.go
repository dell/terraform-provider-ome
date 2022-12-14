package ome

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
)

// DefaultAttributePlanModifier is to set default value for an attribute
// https://github.com/hashicorp/terraform-plugin-framework/issues/285
type DefaultAttributePlanModifier struct {
	value attr.Value
}

// DefaultAttribute modifer for the attributes
func DefaultAttribute(value attr.Value) DefaultAttributePlanModifier {
	return DefaultAttributePlanModifier{value: value}
}

// Modify method for DefaultAttribute
func (m DefaultAttributePlanModifier) Modify(
	ctx context.Context,
	req tfsdk.ModifyAttributePlanRequest,
	resp *tfsdk.ModifyAttributePlanResponse,
) {
	if req.AttributeConfig == nil || resp.AttributePlan == nil {
		return
	}

	// if configuration was provided, then don't use the default
	if !req.AttributeConfig.IsNull() {
		return
	}

	// If the plan is known and not null (for example due to another plan modifier),
	// don't set the default value
	if !resp.AttributePlan.IsUnknown() && !resp.AttributePlan.IsNull() {
		return
	}

	resp.AttributePlan = m.value
}

// Description impl for DefaultAttribute
func (m DefaultAttributePlanModifier) Description(ctx context.Context) string {
	return "Use a static default value for an attribute"
}

// MarkdownDescription for DefaultAttribute
func (m DefaultAttributePlanModifier) MarkdownDescription(ctx context.Context) string {
	return m.Description(ctx)
}
