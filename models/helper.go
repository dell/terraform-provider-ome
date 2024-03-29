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

package models

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

func stringPointerValue(k *string) types.String {
	ret := types.StringNull()
	if k != nil {
		ret = types.StringValue(*k)
	}
	return ret
}

func int64ListValue(inputs []int64) types.List {
	retVals := []attr.Value{}
	for _, input := range inputs {
		retVals = append(retVals, types.Int64Value(input))
	}
	ret, _ := types.ListValue(
		types.Int64Type,
		retVals,
	)
	return ret
}

func stringListValue(inputs []string) types.List {
	ret, _ := types.ListValueFrom(
		context.TODO(),
		types.StringType,
		inputs,
	)
	return ret
}

func objListValue(typeObj map[string]attr.Type, inputs any) (types.List, diag.Diagnostics) {
	return types.ListValueFrom(
		context.TODO(),
		basetypes.ObjectType{
			AttrTypes: typeObj,
		},
		inputs,
	)
}
