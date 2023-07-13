package ome

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func OmeDiscoveryJobSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{

		"discovery_config_group_id": schema.Int64Attribute{
			MarkdownDescription: "Discovery Config Group ID",
			Description:         "Discovery Config Group ID",
			Optional:            true,
			Computed:            true,
		},

		"discovery_config_group_name": schema.StringAttribute{
			MarkdownDescription: "Discovery Config Group Name",
			Description:         "Discovery Config Group Name",
			Optional:            true,
			Computed:            true,
		},

		"discovery_config_group_description": schema.StringAttribute{
			MarkdownDescription: "Discovery Config Group Description",
			Description:         "Discovery Config Group Description",
			Optional:            true,
			Computed:            true,
		},

		"discovery_status_email_recipient": schema.StringAttribute{
			MarkdownDescription: "Discovery Status Email Recipient",
			Description:         "Discovery Status Email Recipient",
			Optional:            true,
			Computed:            true,
		},

		"discovery_config_parent_group_id": schema.Int64Attribute{
			MarkdownDescription: "Discovery Config Parent Group ID",
			Description:         "Discovery Config Parent Group ID",
			Optional:            true,
			Computed:            true,
		},

		"create_group": schema.BoolAttribute{
			MarkdownDescription: "Create Group",
			Description:         "Create Group",
			Optional:            true,
			Computed:            true,
		},

		"discovery_config_models": schema.SetNestedAttribute{
			MarkdownDescription: "Discovery Config Models",
			Description:         "Discovery Config Models",
			Optional:            true,
			Computed:            true,
			NestedObject:        schema.NestedAttributeObject{Attributes: OmeDiscoveryConfigModelsSchema()},
		},

		"discovery_config_task_param": schema.SetNestedAttribute{
			MarkdownDescription: "Discovery Config Task Param",
			Description:         "Discovery Config Task Param",
			Optional:            true,
			Computed:            true,
			NestedObject:        schema.NestedAttributeObject{Attributes: OmeDiscoveryConfigTaskParamSchema()},
		},

		"discovery_config_tasks": schema.SetNestedAttribute{
			MarkdownDescription: "Discovery Config Tasks",
			Description:         "Discovery Config Tasks",
			Optional:            true,
			Computed:            true,
			NestedObject:        schema.NestedAttributeObject{Attributes: OmeDiscoveryConfigTasksSchema()},
		},

		"schedule": schema.SingleNestedAttribute{
			MarkdownDescription: "Schedule",
			Description:         "Schedule",
			Optional:            true,
			Computed:            true,
			Attributes:          OmeScheduleJobSchema(),
		},

		"trap_destination": schema.BoolAttribute{
			MarkdownDescription: "Trap Destination",
			Description:         "Trap Destination",
			Optional:            true,
			Computed:            true,
		},

		"community_string": schema.BoolAttribute{
			MarkdownDescription: "Community String",
			Description:         "Community String",
			Optional:            true,
			Computed:            true,
		},

		"chassis_identifier": schema.StringAttribute{
			MarkdownDescription: "Chassis Identifier",
			Description:         "Chassis Identifier",
			Optional:            true,
			Computed:            true,
		},

		"use_all_profiles": schema.BoolAttribute{
			MarkdownDescription: "Use All Profiles",
			Description:         "Use All Profiles",
			Optional:            true,
			Computed:            true,
		},
	}
}

func OmeDiscoveryConfigTargetsSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{

		"discovery_config_target_id": schema.Int64Attribute{
			MarkdownDescription: "Discovery Config Target ID",
			Description:         "Discovery Config Target ID",
			Optional:            true,
			Computed:            true,
		},

		"network_address_detail": schema.StringAttribute{
			MarkdownDescription: "Network Address Detail",
			Description:         "Network Address Detail",
			Optional:            true,
			Computed:            true,
		},

		"subnet_mask": schema.StringAttribute{
			MarkdownDescription: "Subnet Mask",
			Description:         "Subnet Mask",
			Optional:            true,
			Computed:            true,
		},

		"address_type": schema.Int64Attribute{
			MarkdownDescription: "Address Type",
			Description:         "Address Type",
			Optional:            true,
			Computed:            true,
		},

		"disabled": schema.BoolAttribute{
			MarkdownDescription: "Disabled",
			Description:         "Disabled",
			Optional:            true,
			Computed:            true,
		},

		"exclude": schema.BoolAttribute{
			MarkdownDescription: "Exclude",
			Description:         "Exclude",
			Optional:            true,
			Computed:            true,
		},
	}
}

func OmeDiscoveryConfigModelsSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{

		"discovery_config_id": schema.Int64Attribute{
			MarkdownDescription: "Discovery Config ID",
			Description:         "Discovery Config ID",
			Optional:            true,
			Computed:            true,
		},

		"discovery_config_description": schema.StringAttribute{
			MarkdownDescription: "Discovery Config Description",
			Description:         "Discovery Config Description",
			Optional:            true,
			Computed:            true,
		},

		"discovery_config_status": schema.StringAttribute{
			MarkdownDescription: "Discovery Config Status",
			Description:         "Discovery Config Status",
			Optional:            true,
			Computed:            true,
		},

		"discovery_config_targets": schema.SetNestedAttribute{
			MarkdownDescription: "Discovery Config Targets",
			Description:         "Discovery Config Targets",
			Optional:            true,
			Computed:            true,
			NestedObject:        schema.NestedAttributeObject{Attributes: OmeDiscoveryConfigTargetsSchema()},
		},

		"connection_profile_id": schema.Int64Attribute{
			MarkdownDescription: "Connection Profile ID",
			Description:         "Connection Profile ID",
			Optional:            true,
			Computed:            true,
		},

		"connection_profile": schema.StringAttribute{
			MarkdownDescription: "Connection Profile",
			Description:         "Connection Profile",
			Optional:            true,
			Computed:            true,
		},

		"device_type": schema.ListAttribute{
			MarkdownDescription: "Device Type",
			Description:         "Device Type",
			Optional:            true,
			Computed:            true,
			ElementType:         types.Int64Type,
		},

		"discovery_config_vendor_platforms": schema.SetNestedAttribute{
			MarkdownDescription: "Discovery Config Vendor Platforms",
			Description:         "Discovery Config Vendor Platforms",
			Optional:            true,
			Computed:            true,
			NestedObject:        schema.NestedAttributeObject{Attributes: OmeDiscoveryConfigVendorPlatformsSchema()},
		},
	}
}

func OmeDiscoveryConfigTaskParamSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{

		"task_id": schema.Int64Attribute{
			MarkdownDescription: "Task ID",
			Description:         "Task ID",
			Optional:            true,
			Computed:            true,
		},

		"task_type_id": schema.Int64Attribute{
			MarkdownDescription: "Task Type ID",
			Description:         "Task Type ID",
			Optional:            true,
			Computed:            true,
		},

		"execution_sequence": schema.Int64Attribute{
			MarkdownDescription: "Execution Sequence",
			Description:         "Execution Sequence",
			Optional:            true,
			Computed:            true,
		},
	}
}

func OmeScheduleJobSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{

		"run_now": schema.BoolAttribute{
			MarkdownDescription: "Run Now",
			Description:         "Run Now",
			Optional:            true,
			Computed:            true,
		},

		"run_later": schema.BoolAttribute{
			MarkdownDescription: "Run Later",
			Description:         "Run Later",
			Optional:            true,
			Computed:            true,
		},

		"recurring": schema.SingleNestedAttribute{
			MarkdownDescription: "Recurring",
			Description:         "Recurring",
			Optional:            true,
			Computed:            true,
			Attributes:          OmeRecurringSchema(),
		},

		"cron": schema.StringAttribute{
			MarkdownDescription: "Cron",
			Description:         "Cron",
			Optional:            true,
			Computed:            true,
		},

		"start_time": schema.StringAttribute{
			MarkdownDescription: "Start Time",
			Description:         "Start Time",
			Optional:            true,
			Computed:            true,
		},

		"end_time": schema.StringAttribute{
			MarkdownDescription: "End Time",
			Description:         "End Time",
			Optional:            true,
			Computed:            true,
		},
	}
}

func OmeDiscoveryConfigTasksSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{

		"discovery_config_description": schema.StringAttribute{
			MarkdownDescription: "Discovery Config Description",
			Description:         "Discovery Config Description",
			Optional:            true,
			Computed:            true,
		},

		"discovery_config_email_recipient": schema.StringAttribute{
			MarkdownDescription: "Discovery Config Email Recipient",
			Description:         "Discovery Config Email Recipient",
			Optional:            true,
			Computed:            true,
		},

		"discovery_config_discovered_device_count": schema.StringAttribute{
			MarkdownDescription: "Discovery Config Discovered Device Count",
			Description:         "Discovery Config Discovered Device Count",
			Optional:            true,
			Computed:            true,
		},

		"discovery_config_request_id": schema.Int64Attribute{
			MarkdownDescription: "Discovery Config Request Id",
			Description:         "Discovery Config Request Id",
			Optional:            true,
			Computed:            true,
		},

		"discovery_config_expected_device_count": schema.StringAttribute{
			MarkdownDescription: "Discovery Config Expected Device Count",
			Description:         "Discovery Config Expected Device Count",
			Optional:            true,
			Computed:            true,
		},

		"discovery_config_name": schema.StringAttribute{
			MarkdownDescription: "Discovery Config Name",
			Description:         "Discovery Config Name",
			Optional:            true,
			Computed:            true,
		},
	}
}

func OmeDiscoveryConfigVendorPlatformsSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{

		"vendor_platform_id": schema.Int64Attribute{
			MarkdownDescription: "Vendor Platform Id",
			Description:         "Vendor Platform Id",
			Optional:            true,
			Computed:            true,
		},

		"discovery_config_vendor_platform_id": schema.Int64Attribute{
			MarkdownDescription: "Discovery Config Vendor Platform Id",
			Description:         "Discovery Config Vendor Platform Id",
			Optional:            true,
			Computed:            true,
		},
	}
}

func OmeRecurringSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{

		"hourly": schema.SingleNestedAttribute{
			MarkdownDescription: "Hourly",
			Description:         "Hourly",
			Optional:            true,
			Computed:            true,
			Attributes:          OmeHourlySchema(),
		},

		"daily": schema.SingleNestedAttribute{
			MarkdownDescription: "Daily",
			Description:         "Daily",
			Optional:            true,
			Computed:            true,
			Attributes:          OmeDailySchema(),
		},

		"weekley": schema.SingleNestedAttribute{
			MarkdownDescription: "Weekley",
			Description:         "Weekley",
			Optional:            true,
			Computed:            true,
			Attributes:          OmeWeekleySchema(),
		},
	}
}

func OmeHourlySchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{

		"frequency": schema.Int64Attribute{
			MarkdownDescription: "Frequency",
			Description:         "Frequency",
			Optional:            true,
			Computed:            true,
		},
	}
}

func OmeDailySchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{

		"frequency": schema.Int64Attribute{
			MarkdownDescription: "Frequency",
			Description:         "Frequency",
			Optional:            true,
			Computed:            true,
		},

		"time": schema.SingleNestedAttribute{
			MarkdownDescription: "Time",
			Description:         "Time",
			Optional:            true,
			Computed:            true,
			Attributes:          OmeTimeSchema(),
		},
	}
}

func OmeWeekleySchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{

		"day": schema.StringAttribute{
			MarkdownDescription: "Day",
			Description:         "Day",
			Optional:            true,
			Computed:            true,
		},

		"time": schema.SingleNestedAttribute{
			MarkdownDescription: "Time",
			Description:         "Time",
			Optional:            true,
			Computed:            true,
			Attributes:          OmeTimeSchema(),
		},
	}
}

func OmeTimeSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{

		"minutes": schema.Int64Attribute{
			MarkdownDescription: "Minutes",
			Description:         "Minutes",
			Optional:            true,
			Computed:            true,
		},

		"hour": schema.Int64Attribute{
			MarkdownDescription: "Hour",
			Description:         "Hour",
			Optional:            true,
			Computed:            true,
		},
	}
}
