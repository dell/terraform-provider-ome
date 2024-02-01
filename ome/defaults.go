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

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Int64DefaultValue set default value for int64 type
func Int64DefaultValue(v types.Int64) planmodifier.Int64 {
	return &int64DefaultValuePlanModifier{v}
}

type int64DefaultValuePlanModifier struct {
	DefaultValue types.Int64
}

var _ planmodifier.Int64 = (*int64DefaultValuePlanModifier)(nil)

func (apm *int64DefaultValuePlanModifier) Description(ctx context.Context) string {
	return fmt.Sprintf("If value is not configured, defaults to %s", apm.DefaultValue)
}

func (apm *int64DefaultValuePlanModifier) MarkdownDescription(ctx context.Context) string {
	return fmt.Sprintf("If value is not configured, defaults to %s", apm.DefaultValue)
}

func (apm *int64DefaultValuePlanModifier) PlanModifyInt64(ctx context.Context, req planmodifier.Int64Request, res *planmodifier.Int64Response) {
	// If the attribute configuration is not null, we are done here
	if !req.ConfigValue.IsNull() {
		return
	}

	// If the attribute plan is "known" and "not null", then a previous plan modifier in the sequence
	// has already been applied, and we don't want to interfere.
	if !req.PlanValue.IsUnknown() && !req.PlanValue.IsNull() {
		return
	}

	res.PlanValue = apm.DefaultValue
}

// StringDefaultValue set default value for string type
func StringDefaultValue(v types.String) planmodifier.String {
	return &stringDefaultValuePlanModifier{v}
}

type stringDefaultValuePlanModifier struct {
	DefaultValue types.String
}

// Description implements planmodifier.String
func (apm *stringDefaultValuePlanModifier) Description(context.Context) string {
	return fmt.Sprintf("If value is not configured, defaults to %s", apm.DefaultValue)
}

// MarkdownDescription implements planmodifier.String
func (apm *stringDefaultValuePlanModifier) MarkdownDescription(context.Context) string {
	return fmt.Sprintf("If value is not configured, defaults to %s", apm.DefaultValue)
}

// PlanModifyString implements planmodifier.String
func (apm *stringDefaultValuePlanModifier) PlanModifyString(ctx context.Context, req planmodifier.StringRequest, res *planmodifier.StringResponse) {
	// If the attribute configuration is not null, we are done here
	if !req.ConfigValue.IsNull() {
		return
	}

	// If the attribute plan is "known" and "not null", then a previous plan modifier in the sequence
	// has already been applied, and we don't want to interfere.
	if !req.PlanValue.IsUnknown() && !req.PlanValue.IsNull() {
		return
	}

	res.PlanValue = apm.DefaultValue
}

// BoolDefaultValue sets default value for bool type
func BoolDefaultValue(v types.Bool) planmodifier.Bool {
	return &BoolDefaultValuePlanModifier{v}
}

// BoolDefaultValuePlanModifier for bool type
type BoolDefaultValuePlanModifier struct {
	DefaultValue types.Bool
}

// Description implements planmodifier.Bool
func (apm *BoolDefaultValuePlanModifier) Description(context.Context) string {
	return fmt.Sprintf("If value is not configured, defaults to %s", apm.DefaultValue)
}

// MarkdownDescription implements planmodifier.Bool
func (apm *BoolDefaultValuePlanModifier) MarkdownDescription(context.Context) string {
	return fmt.Sprintf("If value is not configured, defaults to %s", apm.DefaultValue)
}

// PlanModifyBool implements planmodifier.Bool
func (apm *BoolDefaultValuePlanModifier) PlanModifyBool(ctx context.Context, req planmodifier.BoolRequest, res *planmodifier.BoolResponse) {
	// If the attribute configuration is not null, we are done here
	if !req.ConfigValue.IsNull() {
		return
	}

	// If the attribute plan is "known" and "not null", then a previous plan modifier in the sequence
	// has already been applied, and we don't want to interfere.
	if !req.PlanValue.IsUnknown() && !req.PlanValue.IsNull() {
		return
	}

	res.PlanValue = apm.DefaultValue
}
