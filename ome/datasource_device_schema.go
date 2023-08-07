/*
Copyright (c) 2023 Dell Inc., or its subsidiaries. All Rights Reserved.
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
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func omeDeviceDataSchema() map[string]schema.Attribute {
	acceptedInventoryTypes := []string{
		"serverDeviceCards",
		"serverProcessors",
		"serverDellVideos",
		"serverNetworkInterfaces",
		"serverFcCards",
		"serverOperatingSystems",
		"serverVirtualFlashes",
		"serverPowerSupplies",
		"serverArrayDisks",
		"serverRaidControllers",
		"serverMemoryDevices",
		"serverStorageEnclosures",
		"serverSupportedPowerStates",
		"deviceLicense",
		"deviceCapabilities",
		"deviceFru",
		"deviceManagement",
		"deviceSoftware",
		"subsystemRollupStatus",
		"deviceInventory",
	}
	return map[string]schema.Attribute{
		"filters": schema.SingleNestedAttribute{
			MarkdownDescription: "Filters to apply while fetching devices." +
				" Only one among `filter_expression`, `ids` and `device_service_tags` can be configured.",
			Description: "Filters to apply while fetching devices." +
				" Only one among 'filter_expression', 'ids' and 'device_service_tags' can be configured.",
			Optional: true,
			Attributes: map[string]schema.Attribute{
				"ids": schema.ListAttribute{
					MarkdownDescription: "IDs of the devices to fetch.",
					Description:         "IDs of the devices to fetch.",
					Optional:            true,
					ElementType:         types.Int64Type,
				},
				"device_service_tags": schema.ListAttribute{
					MarkdownDescription: "Service tags of the devices to fetch.",
					Description:         "Service tags of the devices to fetch.",
					Optional:            true,
					ElementType:         types.StringType,
					Validators: []validator.List{
						listvalidator.ConflictsWith(path.MatchRelative().AtName("ids")),
					},
				},
				"ip_expressions": schema.ListAttribute{
					MarkdownDescription: "IP expressions of the devices to fetch." +
						" Supported expressions are IPv4, IPv6, CIDRs and IP ranges.",
					Description: "IP expressions of the devices to fetch." +
						" Supported expressions are IPv4, IPv6, CIDRs and IP ranges.",
					Optional:    true,
					ElementType: types.StringType,
				},
				"filter_expression": schema.StringAttribute{
					MarkdownDescription: "OData `$filter` compatible expression to be used for querying devices.",
					Description:         "OData '$filter' compatible expression to be used for querying devices.",
					Optional:            true,
					Validators: []validator.String{
						stringvalidator.ConflictsWith(path.MatchRelative().AtName("ids")),
						stringvalidator.ConflictsWith(path.MatchRelative().AtName("device_service_tags")),
					},
				},
			},
		},
		"devices": schema.ListNestedAttribute{
			MarkdownDescription: "Devices fetched.",
			Description:         "Devices fetched.",
			Computed:            true,
			NestedObject:        schema.NestedAttributeObject{Attributes: omeSingleDeviceDataSchema()},
		},
		"id": schema.Int64Attribute{
			MarkdownDescription: "Dummy ID of the datasource.",
			Description:         "Dummy ID of the datasource.",
			Computed:            true,
		},
		"inventory_types": schema.ListAttribute{
			MarkdownDescription: "The types of inventory types to fetch." +
				makeSchemaAcceptedValues(acceptedInventoryTypes, "`") +
				" If not configured, all inventory types are fetched.",
			Description: "The types of inventory to fetch." +
				makeSchemaAcceptedValues(acceptedInventoryTypes, "'") +
				" If not configured, all inventory types are fetched.",
			Optional:    true,
			ElementType: types.StringType,
			Validators: []validator.List{
				listvalidator.ValueStringsAre(
					stringvalidator.OneOf(acceptedInventoryTypes...),
				),
				listvalidator.UniqueValues(),
			},
		},
	}
}

func omeSingleDeviceDataSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"id": schema.Int64Attribute{
			MarkdownDescription: "ID of the device.",
			Description:         "ID of the device.",
			Computed:            true,
		},
		"type": schema.Int64Attribute{
			MarkdownDescription: "Type of the device.",
			Description:         "Type of the device.",
			Computed:            true,
		},
		"identifier": schema.StringAttribute{
			MarkdownDescription: "Identifier of the device.",
			Description:         "Identifier of the device.",
			Computed:            true,
		},
		"device_service_tag": schema.StringAttribute{
			MarkdownDescription: "Device Service Tag of the device.",
			Description:         "Device Service Tag of the device.",
			Computed:            true,
		},
		"chassis_service_tag": schema.StringAttribute{
			MarkdownDescription: "Chassis Service Tag of the device.",
			Description:         "Chassis Service Tag of the device.",
			Computed:            true,
		},
		"model": schema.StringAttribute{
			MarkdownDescription: "Model of the device.",
			Description:         "Model of the device.",
			Computed:            true,
		},
		"power_state": schema.Int64Attribute{
			MarkdownDescription: "Power State of the device.",
			Description:         "Power State of the device.",
			Computed:            true,
		},
		"managed_state": schema.Int64Attribute{
			MarkdownDescription: "Managed State of the device.",
			Description:         "Managed State of the device.",
			Computed:            true,
		},
		"status": schema.Int64Attribute{
			MarkdownDescription: "Status of the device.",
			Description:         "Status of the device.",
			Computed:            true,
		},
		"connection_state": schema.BoolAttribute{
			MarkdownDescription: "Connection State of the device.",
			Description:         "Connection State of the device.",
			Computed:            true,
		},
		"asset_tag": schema.StringAttribute{
			MarkdownDescription: "Asset Tag of the device.",
			Description:         "Asset Tag of the device.",
			Computed:            true,
		},
		"system_id": schema.Int64Attribute{
			MarkdownDescription: "System ID of the device.",
			Description:         "System ID of the device.",
			Computed:            true,
		},
		"device_name": schema.StringAttribute{
			MarkdownDescription: "Device Name of the device.",
			Description:         "Device Name of the device.",
			Computed:            true,
		},
		"last_inventory_time": schema.StringAttribute{
			MarkdownDescription: "Last Inventory Time of the device.",
			Description:         "Last Inventory Time of the device.",
			Computed:            true,
		},
		"last_status_time": schema.StringAttribute{
			MarkdownDescription: "Last Status Time of the device.",
			Description:         "Last Status Time of the device.",
			Computed:            true,
		},
		"device_subscription": schema.StringAttribute{
			MarkdownDescription: "Device Subscription of the device.",
			Description:         "Device Subscription of the device.",
			Computed:            true,
		},
		"device_capabilities": schema.ListAttribute{
			MarkdownDescription: "Device Capabilities of the device.",
			Description:         "Device Capabilities of the device.",
			Computed:            true,
			ElementType:         types.Int64Type,
		},
		"slot_configuration": schema.SingleNestedAttribute{
			MarkdownDescription: "Slot Configuration of the device.",
			Description:         "Slot Configuration of the device.",
			Computed:            true,
			Attributes:          omeSlotConfigurationDataSchema(),
		},
		"device_management": schema.ListNestedAttribute{
			MarkdownDescription: "Device Management of the device.",
			Description:         "Device Management of the device.",
			Computed:            true,
			NestedObject:        schema.NestedAttributeObject{Attributes: omeDeviceManagementDataSchema()},
		},
		"enabled": schema.BoolAttribute{
			MarkdownDescription: "Whether the device is enabled or not.",
			Description:         "Whether the device is enabled or not.",
			Computed:            true,
		},
		"connection_state_reason": schema.Int64Attribute{
			MarkdownDescription: "Connection State Reason of the device.",
			Description:         "Connection State Reason of the device.",
			Computed:            true,
		},
		"chassis_ip": schema.StringAttribute{
			MarkdownDescription: "Chassis IP of the device.",
			Description:         "Chassis IP of the device.",
			Computed:            true,
		},
		"discovery_configuration_job_information": schema.ListNestedAttribute{
			MarkdownDescription: "Discovery Configuration Job Info of the device.",
			Description:         "Discovery Configuration Job Info of the device.",
			Computed:            true,
			NestedObject:        schema.NestedAttributeObject{Attributes: omeDiscoveryConfigurationJobDataSchema()},
		},
		"detailed_inventory": schema.SingleNestedAttribute{
			MarkdownDescription: "Detailed inventory of the device." +
				" Detailed inventory is only fetched if only a single device is fetched by this datasource.",
			Description: "Detailed inventory of the device." +
				" Detailed inventory is only fetched if only a single device is fetched by this datasource.",
			Computed:   true,
			Attributes: omeDeviceInventorySchema(),
		},
	}
}

func omeDiscoveryConfigurationJobDataSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"group_id": schema.StringAttribute{
			MarkdownDescription: "Group ID",
			Description:         "Group ID",
			Computed:            true,
		},
		"created_by": schema.StringAttribute{
			MarkdownDescription: "Created By",
			Description:         "Created By",
			Computed:            true,
		},
		"discovery_job_name": schema.StringAttribute{
			MarkdownDescription: "Discovery Job Name",
			Description:         "Discovery Job Name",
			Computed:            true,
		},
	}
}

func omeSlotConfigurationDataSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"chassis_name": schema.StringAttribute{
			MarkdownDescription: "Chassis Name",
			Description:         "Chassis Name",
			Computed:            true,
		},
	}
}

func omeDeviceManagementDataSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"management_id": schema.Int64Attribute{
			MarkdownDescription: "Management ID",
			Description:         "Management ID",
			Computed:            true,
		},
		"network_address": schema.StringAttribute{
			MarkdownDescription: "Network Address",
			Description:         "Network Address",
			Computed:            true,
		},
		"mac_address": schema.StringAttribute{
			MarkdownDescription: "Mac Address",
			Description:         "Mac Address",
			Computed:            true,
		},
		"management_type": schema.Int64Attribute{
			MarkdownDescription: "Management Type",
			Description:         "Management Type",
			Computed:            true,
		},
		"instrumentation_name": schema.StringAttribute{
			MarkdownDescription: "Instrumentation Name",
			Description:         "Instrumentation Name",
			Computed:            true,
		},
		"dns_name": schema.StringAttribute{
			MarkdownDescription: "DNSName",
			Description:         "DNSName",
			Computed:            true,
		},
		"management_profile": schema.SetNestedAttribute{
			MarkdownDescription: "Management Profile",
			Description:         "Management Profile",
			Computed:            true,
			NestedObject:        schema.NestedAttributeObject{Attributes: omeManagementProfileDataSchema()},
		},
	}
}

func omeManagementProfileDataSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"management_profile_id": schema.Int64Attribute{
			MarkdownDescription: "Management Profile ID",
			Description:         "Management Profile ID",
			Computed:            true,
		},
		"profile_id": schema.StringAttribute{
			MarkdownDescription: "Profile ID",
			Description:         "Profile ID",
			Computed:            true,
		},
		"management_id": schema.Int64Attribute{
			MarkdownDescription: "Management ID",
			Description:         "Management ID",
			Computed:            true,
		},
		"agent_name": schema.StringAttribute{
			MarkdownDescription: "Agent Name",
			Description:         "Agent Name",
			Computed:            true,
		},
		"version": schema.StringAttribute{
			MarkdownDescription: "Version",
			Description:         "Version",
			Computed:            true,
		},
		"management_url": schema.StringAttribute{
			MarkdownDescription: "Management URL",
			Description:         "Management URL",
			Computed:            true,
		},
		"has_creds": schema.Int64Attribute{
			MarkdownDescription: "Has Creds",
			Description:         "Has Creds",
			Computed:            true,
		},
		"status": schema.Int64Attribute{
			MarkdownDescription: "Status",
			Description:         "Status",
			Computed:            true,
		},
		"status_date_time": schema.StringAttribute{
			MarkdownDescription: "Status Date Time",
			Description:         "Status Date Time",
			Computed:            true,
		},
	}
}

// ################## Inventory Schema

func omeDeviceInventorySchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"server_device_cards": schema.ListNestedAttribute{
			MarkdownDescription: "Server Device Cards.",
			Description:         "Server Device Cards.",
			Computed:            true,
			NestedObject:        schema.NestedAttributeObject{Attributes: omeServerDeviceCardInfoSchema()},
		},
		"cpus": schema.ListNestedAttribute{
			MarkdownDescription: "CPU related Information.",
			Description:         "CPU related Information.",
			Computed:            true,
			NestedObject:        schema.NestedAttributeObject{Attributes: omeCPUInfoSchema()},
		},

		"nics": schema.ListNestedAttribute{
			MarkdownDescription: "NIC related Information.",
			Description:         "NIC related Information.",

			Computed:     true,
			NestedObject: schema.NestedAttributeObject{Attributes: omeNICInfoSchema()},
		},

		"fcis": schema.ListNestedAttribute{
			MarkdownDescription: "FCI related Information.",
			Description:         "FCI related Information.",

			Computed:     true,
			NestedObject: schema.NestedAttributeObject{Attributes: omeFCInfoSchema()},
		},

		"os": schema.ListNestedAttribute{
			MarkdownDescription: "OS related Information.",
			Description:         "OS related Information.",

			Computed:     true,
			NestedObject: schema.NestedAttributeObject{Attributes: omeOSInfoSchema()},
		},

		"power_supply": schema.ListNestedAttribute{
			MarkdownDescription: "Power Supply related Information.",
			Description:         "Power Supply related Information.",

			Computed:     true,
			NestedObject: schema.NestedAttributeObject{Attributes: omePowerSupplyInfoSchema()},
		},

		"disks": schema.ListNestedAttribute{
			MarkdownDescription: "Disk related Information.",
			Description:         "Disk related Information.",

			Computed:     true,
			NestedObject: schema.NestedAttributeObject{Attributes: omeDiskInfoSchema()},
		},

		"raid_controllers": schema.ListNestedAttribute{
			MarkdownDescription: "RAIDController Information.",
			Description:         "RAIDController Information.",
			Computed:            true,
			NestedObject:        schema.NestedAttributeObject{Attributes: omeRAIDControllerInfoSchema()},
		},

		"memory": schema.ListNestedAttribute{
			MarkdownDescription: "Memory Information.",
			Description:         "Memory Information.",
			Computed:            true,
			NestedObject:        schema.NestedAttributeObject{Attributes: omeMemoryInfoSchema()},
		},
		"storage_enclosures": schema.ListNestedAttribute{
			MarkdownDescription: "Storage Enclosure Information.",
			Description:         "Storage Enclosure Information.",
			Computed:            true,
			NestedObject:        schema.NestedAttributeObject{Attributes: omeStorageEnclosureInfoSchema()},
		},

		"power_state": schema.ListNestedAttribute{
			MarkdownDescription: "Server Power States",
			Description:         "Server Power States",

			Computed:     true,
			NestedObject: schema.NestedAttributeObject{Attributes: omeServerPowerStateSchema()},
		},

		"licenses": schema.ListNestedAttribute{
			MarkdownDescription: "Device Licenses",
			Description:         "Device Licenses",

			Computed:     true,
			NestedObject: schema.NestedAttributeObject{Attributes: omeDeviceLicenseSchema()},
		},

		"capabilities": schema.ListNestedAttribute{
			MarkdownDescription: "Device Capabilities",
			Description:         "Device Capabilities",

			Computed:     true,
			NestedObject: schema.NestedAttributeObject{Attributes: omeDeviceCapabilitySchema()},
		},

		"frus": schema.ListNestedAttribute{
			MarkdownDescription: "Device FRUs",
			Description:         "Device FRUs",

			Computed:     true,
			NestedObject: schema.NestedAttributeObject{Attributes: omeDeviceFruSchema()},
		},

		"locations": schema.ListNestedAttribute{
			MarkdownDescription: "Device Locations",
			Description:         "Device Locations",

			Computed:     true,
			NestedObject: schema.NestedAttributeObject{Attributes: omeDeviceLocationSchema()},
		},

		"management_info": schema.ListNestedAttribute{
			MarkdownDescription: "Device Management",
			Description:         "Device Management",

			Computed:     true,
			NestedObject: schema.NestedAttributeObject{Attributes: omeDeviceManagementInfoSchema()},
		},

		"softwares": schema.ListNestedAttribute{
			MarkdownDescription: "Device Softwares",
			Description:         "Device Softwares",

			Computed:     true,
			NestedObject: schema.NestedAttributeObject{Attributes: omeDeviceSoftwareSchema()},
		},

		"subsytem_rollup_status": schema.ListNestedAttribute{
			MarkdownDescription: "Sub System Rollup Status",
			Description:         "Sub System Rollup Status",

			Computed:     true,
			NestedObject: schema.NestedAttributeObject{Attributes: omeSubSystemRollupStatusSchema()},
		},
	}
}

func omeSubSystemRollupStatusSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{

		"id": schema.Int64Attribute{
			MarkdownDescription: "ID",
			Description:         "ID",

			Computed: true,
		},

		"status": schema.Int64Attribute{
			MarkdownDescription: "Status",
			Description:         "Status",

			Computed: true,
		},

		"subsystem_name": schema.StringAttribute{
			MarkdownDescription: "Subsystem Name",
			Description:         "Subsystem Name",

			Computed: true,
		},
	}
}

func omeDeviceSoftwareSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{

		"version": schema.StringAttribute{
			MarkdownDescription: "Version",
			Description:         "Version",

			Computed: true,
		},

		"installation_date": schema.StringAttribute{
			MarkdownDescription: "Installation Date",
			Description:         "Installation Date",

			Computed: true,
		},

		"status": schema.StringAttribute{
			MarkdownDescription: "Status",
			Description:         "Status",

			Computed: true,
		},

		"software_type": schema.StringAttribute{
			MarkdownDescription: "Software Type",
			Description:         "Software Type",

			Computed: true,
		},

		"vendor_id": schema.StringAttribute{
			MarkdownDescription: "Vendor ID",
			Description:         "Vendor ID",

			Computed: true,
		},

		"sub_device_id": schema.StringAttribute{
			MarkdownDescription: "Sub Device ID",
			Description:         "Sub Device ID",

			Computed: true,
		},

		"sub_vendor_id": schema.StringAttribute{
			MarkdownDescription: "Sub Vendor ID",
			Description:         "Sub Vendor ID",

			Computed: true,
		},

		"component_id": schema.StringAttribute{
			MarkdownDescription: "Component ID",
			Description:         "Component ID",

			Computed: true,
		},

		"pci_device_id": schema.StringAttribute{
			MarkdownDescription: "Pci Device ID",
			Description:         "Pci Device ID",

			Computed: true,
		},

		"device_description": schema.StringAttribute{
			MarkdownDescription: "Device Description",
			Description:         "Device Description",

			Computed: true,
		},

		"instance_id": schema.StringAttribute{
			MarkdownDescription: "Instance ID",
			Description:         "Instance ID",

			Computed: true,
		},
	}
}

func omeDeviceManagementInfoSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{

		"management_id": schema.Int64Attribute{
			MarkdownDescription: "Management ID",
			Description:         "Management ID",

			Computed: true,
		},

		"ip_address": schema.StringAttribute{
			MarkdownDescription: "IPAddress",
			Description:         "IPAddress",

			Computed: true,
		},

		"mac_address": schema.StringAttribute{
			MarkdownDescription: "MACAddress",
			Description:         "MACAddress",

			Computed: true,
		},

		"instrumentation_name": schema.StringAttribute{
			MarkdownDescription: "Instrumentation Name",
			Description:         "Instrumentation Name",

			Computed: true,
		},

		"dns_name": schema.StringAttribute{
			MarkdownDescription: "DNSName",
			Description:         "DNSName",

			Computed: true,
		},

		"management_type": schema.SingleNestedAttribute{
			MarkdownDescription: "Management Type",
			Description:         "Management Type",

			Computed:   true,
			Attributes: omeManagementTypeSchema(),
		},

		"end_point_agents": schema.ListNestedAttribute{
			MarkdownDescription: "End Point Agents",
			Description:         "End Point Agents",

			Computed:     true,
			NestedObject: schema.NestedAttributeObject{Attributes: omeEndPointAgentSchema()},
		},
	}
}

func omeManagementTypeSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{

		"management_type": schema.Int64Attribute{
			MarkdownDescription: "Management Type",
			Description:         "Management Type",

			Computed: true,
		},

		"name": schema.StringAttribute{
			MarkdownDescription: "Name",
			Description:         "Name",

			Computed: true,
		},

		"description": schema.StringAttribute{
			MarkdownDescription: "Description",
			Description:         "Description",

			Computed: true,
		},
	}
}

func omeEndPointAgentSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{

		"management_profile_id": schema.Int64Attribute{
			MarkdownDescription: "Management Profile ID",
			Description:         "Management Profile ID",

			Computed: true,
		},

		"profile_id": schema.StringAttribute{
			MarkdownDescription: "Profile ID",
			Description:         "Profile ID",

			Computed: true,
		},

		"agent_name": schema.StringAttribute{
			MarkdownDescription: "Agent Name",
			Description:         "Agent Name",

			Computed: true,
		},

		"version": schema.StringAttribute{
			MarkdownDescription: "Version",
			Description:         "Version",

			Computed: true,
		},

		"management_url": schema.StringAttribute{
			MarkdownDescription: "Management URL",
			Description:         "Management URL",

			Computed: true,
		},

		"has_creds": schema.Int64Attribute{
			MarkdownDescription: "Has Creds",
			Description:         "Has Creds",

			Computed: true,
		},

		"status": schema.Int64Attribute{
			MarkdownDescription: "Status",
			Description:         "Status",

			Computed: true,
		},

		"status_date_time": schema.StringAttribute{
			MarkdownDescription: "Status Date Time",
			Description:         "Status Date Time",

			Computed: true,
		},
	}
}

func omeDeviceLocationSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"id": schema.Int64Attribute{
			MarkdownDescription: "ID",
			Description:         "ID",
			Computed:            true,
		},
		"room": schema.StringAttribute{
			MarkdownDescription: "Room",
			Description:         "Room",
			Computed:            true,
		},
		"rack": schema.StringAttribute{
			MarkdownDescription: "Rack",
			Description:         "Rack",
			Computed:            true,
		},
		"aisle": schema.StringAttribute{
			MarkdownDescription: "Aisle",
			Description:         "Aisle",
			Computed:            true,
		},
		"datacenter": schema.StringAttribute{
			MarkdownDescription: "Datacenter",
			Description:         "Datacenter",
			Computed:            true,
		},
		"rackslot": schema.StringAttribute{
			MarkdownDescription: "Rackslot",
			Description:         "Rackslot",
			Computed:            true,
		},
		"management_system_unit": schema.Int64Attribute{
			MarkdownDescription: "Management System Unit",
			Description:         "Management System Unit",
			Computed:            true,
		},
	}
}

func omeDeviceFruSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{

		"revision": schema.StringAttribute{
			MarkdownDescription: "Revision",
			Description:         "Revision",

			Computed: true,
		},

		"id": schema.Int64Attribute{
			MarkdownDescription: "ID",
			Description:         "ID",

			Computed: true,
		},

		"manufacturer": schema.StringAttribute{
			MarkdownDescription: "Manufacturer",
			Description:         "Manufacturer",

			Computed: true,
		},

		"name": schema.StringAttribute{
			MarkdownDescription: "Name",
			Description:         "Name",

			Computed: true,
		},

		"part_number": schema.StringAttribute{
			MarkdownDescription: "Part Number",
			Description:         "Part Number",

			Computed: true,
		},

		"serial_number": schema.StringAttribute{
			MarkdownDescription: "Serial Number",
			Description:         "Serial Number",

			Computed: true,
		},
	}
}

func omeDeviceCapabilitySchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{

		"id": schema.Int64Attribute{
			MarkdownDescription: "ID",
			Description:         "ID",

			Computed: true,
		},

		"capability_type": schema.SingleNestedAttribute{
			MarkdownDescription: "Capability Type",
			Description:         "Capability Type",

			Computed:   true,
			Attributes: omeDeviceCapabilityTypeSchema(),
		},
	}
}

func omeDeviceCapabilityTypeSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{

		"capability_id": schema.Int64Attribute{
			MarkdownDescription: "Capability ID",
			Description:         "Capability ID",

			Computed: true,
		},

		"name": schema.StringAttribute{
			MarkdownDescription: "Name",
			Description:         "Name",

			Computed: true,
		},

		"description": schema.StringAttribute{
			MarkdownDescription: "Description",
			Description:         "Description",

			Computed: true,
		},

		"id_owner": schema.Int64Attribute{
			MarkdownDescription: "IDOwner",
			Description:         "IDOwner",

			Computed: true,
		},
	}
}

func omeServerDeviceCardInfoSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{

		"id": schema.Int64Attribute{
			MarkdownDescription: "ID",
			Description:         "ID",

			Computed: true,
		},

		"slot_number": schema.StringAttribute{
			MarkdownDescription: "Slot Number",
			Description:         "Slot Number",

			Computed: true,
		},

		"manufacturer": schema.StringAttribute{
			MarkdownDescription: "Manufacturer",
			Description:         "Manufacturer",

			Computed: true,
		},

		"description": schema.StringAttribute{
			MarkdownDescription: "Description",
			Description:         "Description",

			Computed: true,
		},

		"databus_width": schema.StringAttribute{
			MarkdownDescription: "Databus Width",
			Description:         "Databus Width",

			Computed: true,
		},

		"slot_length": schema.StringAttribute{
			MarkdownDescription: "Slot Length",
			Description:         "Slot Length",

			Computed: true,
		},

		"slot_type": schema.StringAttribute{
			MarkdownDescription: "Slot Type",
			Description:         "Slot Type",

			Computed: true,
		},
	}
}

func omeCPUInfoSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{

		"id": schema.Int64Attribute{
			MarkdownDescription: "ID",
			Description:         "ID",

			Computed: true,
		},

		"family": schema.StringAttribute{
			MarkdownDescription: "Family",
			Description:         "Family",

			Computed: true,
		},

		"max_speed": schema.Int64Attribute{
			MarkdownDescription: "Max Speed",
			Description:         "Max Speed",

			Computed: true,
		},

		"current_speed": schema.Int64Attribute{
			MarkdownDescription: "Current Speed",
			Description:         "Current Speed",

			Computed: true,
		},

		"slot_number": schema.StringAttribute{
			MarkdownDescription: "Slot Number",
			Description:         "Slot Number",

			Computed: true,
		},

		"status": schema.Int64Attribute{
			MarkdownDescription: "Status",
			Description:         "Status",

			Computed: true,
		},

		"number_of_cores": schema.Int64Attribute{
			MarkdownDescription: "Number Of Cores",
			Description:         "Number Of Cores",

			Computed: true,
		},

		"number_of_enabled_cores": schema.Int64Attribute{
			MarkdownDescription: "Number Of Enabled Cores",
			Description:         "Number Of Enabled Cores",

			Computed: true,
		},

		"brand_name": schema.StringAttribute{
			MarkdownDescription: "Brand Name",
			Description:         "Brand Name",

			Computed: true,
		},

		"model_name": schema.StringAttribute{
			MarkdownDescription: "Model Name",
			Description:         "Model Name",

			Computed: true,
		},

		"instance_id": schema.StringAttribute{
			MarkdownDescription: "Instance ID",
			Description:         "Instance ID",

			Computed: true,
		},

		"voltage": schema.StringAttribute{
			MarkdownDescription: "Voltage",
			Description:         "Voltage",

			Computed: true,
		},
	}
}

func omePartitionSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{

		"fqdd": schema.StringAttribute{
			MarkdownDescription: "Fqdd",
			Description:         "Fqdd",

			Computed: true,
		},

		"current_mac_address": schema.StringAttribute{
			MarkdownDescription: "Current Mac Address",
			Description:         "Current Mac Address",

			Computed: true,
		},

		"permanent_mac_address": schema.StringAttribute{
			MarkdownDescription: "Permanent Mac Address",
			Description:         "Permanent Mac Address",

			Computed: true,
		},

		"permanent_iscsi_mac_address": schema.StringAttribute{
			MarkdownDescription: "Permanent Iscsi Mac Address",
			Description:         "Permanent Iscsi Mac Address",

			Computed: true,
		},

		"permanent_fcoe_mac_address": schema.StringAttribute{
			MarkdownDescription: "Permanent Fcoe Mac Address",
			Description:         "Permanent Fcoe Mac Address",

			Computed: true,
		},

		"wwn": schema.StringAttribute{
			MarkdownDescription: "Wwn",
			Description:         "Wwn",

			Computed: true,
		},

		"wwpn": schema.StringAttribute{
			MarkdownDescription: "Wwpn",
			Description:         "Wwpn",

			Computed: true,
		},

		"virtual_wwn": schema.StringAttribute{
			MarkdownDescription: "Virtual Wwn",
			Description:         "Virtual Wwn",

			Computed: true,
		},

		"virtual_wwpn": schema.StringAttribute{
			MarkdownDescription: "Virtual Wwpn",
			Description:         "Virtual Wwpn",

			Computed: true,
		},

		"virtual_mac_address": schema.StringAttribute{
			MarkdownDescription: "Virtual Mac Address",
			Description:         "Virtual Mac Address",

			Computed: true,
		},

		"nic_mode": schema.StringAttribute{
			MarkdownDescription: "Nic Mode",
			Description:         "Nic Mode",

			Computed: true,
		},

		"fcoe_mode": schema.StringAttribute{
			MarkdownDescription: "Fcoe Mode",
			Description:         "Fcoe Mode",

			Computed: true,
		},

		"iscsi_mode": schema.StringAttribute{
			MarkdownDescription: "Iscsi Mode",
			Description:         "Iscsi Mode",

			Computed: true,
		},

		"min_bandwidth": schema.Int64Attribute{
			MarkdownDescription: "Min Bandwidth",
			Description:         "Min Bandwidth",

			Computed: true,
		},

		"max_bandwidth": schema.Int64Attribute{
			MarkdownDescription: "Max Bandwidth",
			Description:         "Max Bandwidth",

			Computed: true,
		},
	}
}

func omePortSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{

		"port_id": schema.StringAttribute{
			MarkdownDescription: "Port ID",
			Description:         "Port ID",

			Computed: true,
		},

		"product_name": schema.StringAttribute{
			MarkdownDescription: "Product Name",
			Description:         "Product Name",

			Computed: true,
		},

		"link_status": schema.StringAttribute{
			MarkdownDescription: "Link Status",
			Description:         "Link Status",

			Computed: true,
		},

		"link_speed": schema.Int64Attribute{
			MarkdownDescription: "Link Speed",
			Description:         "Link Speed",

			Computed: true,
		},

		"partitions": schema.ListNestedAttribute{
			MarkdownDescription: "Partitions",
			Description:         "Partitions",

			Computed:     true,
			NestedObject: schema.NestedAttributeObject{Attributes: omePartitionSchema()},
		},
	}
}

func omeNICInfoSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{

		"nic_id": schema.StringAttribute{
			MarkdownDescription: "Nic ID",
			Description:         "Nic ID",

			Computed: true,
		},

		"vendor_name": schema.StringAttribute{
			MarkdownDescription: "Vendor Name",
			Description:         "Vendor Name",

			Computed: true,
		},

		"ports": schema.ListNestedAttribute{
			MarkdownDescription: "Ports",
			Description:         "Ports",

			Computed:     true,
			NestedObject: schema.NestedAttributeObject{Attributes: omePortSchema()},
		},
	}
}

func omeFCInfoSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{

		"id": schema.Int64Attribute{
			MarkdownDescription: "ID",
			Description:         "ID",

			Computed: true,
		},

		"fqdd": schema.StringAttribute{
			MarkdownDescription: "Fqdd",
			Description:         "Fqdd",

			Computed: true,
		},

		"device_description": schema.StringAttribute{
			MarkdownDescription: "Device Description",
			Description:         "Device Description",

			Computed: true,
		},

		"device_name": schema.StringAttribute{
			MarkdownDescription: "Device Name",
			Description:         "Device Name",

			Computed: true,
		},

		"first_fctarget_lun": schema.StringAttribute{
			MarkdownDescription: "First Fctarget Lun",
			Description:         "First Fctarget Lun",

			Computed: true,
		},

		"first_fctarget_wwpn": schema.StringAttribute{
			MarkdownDescription: "First Fctarget Wwpn",
			Description:         "First Fctarget Wwpn",

			Computed: true,
		},

		"port_number": schema.Int64Attribute{
			MarkdownDescription: "Port Number",
			Description:         "Port Number",

			Computed: true,
		},

		"port_speed": schema.StringAttribute{
			MarkdownDescription: "Port Speed",
			Description:         "Port Speed",

			Computed: true,
		},

		"second_fctarget_lun": schema.StringAttribute{
			MarkdownDescription: "Second Fctarget Lun",
			Description:         "Second Fctarget Lun",

			Computed: true,
		},

		"second_fctarget_wwpn": schema.StringAttribute{
			MarkdownDescription: "Second Fctarget Wwpn",
			Description:         "Second Fctarget Wwpn",

			Computed: true,
		},

		"vendor_name": schema.StringAttribute{
			MarkdownDescription: "Vendor Name",
			Description:         "Vendor Name",

			Computed: true,
		},

		"wwn": schema.StringAttribute{
			MarkdownDescription: "Wwn",
			Description:         "Wwn",

			Computed: true,
		},

		"wwpn": schema.StringAttribute{
			MarkdownDescription: "Wwpn",
			Description:         "Wwpn",

			Computed: true,
		},

		"link_status": schema.StringAttribute{
			MarkdownDescription: "Link Status",
			Description:         "Link Status",

			Computed: true,
		},

		"virtual_wwn": schema.StringAttribute{
			MarkdownDescription: "Virtual Wwn",
			Description:         "Virtual Wwn",

			Computed: true,
		},

		"virtual_wwpn": schema.StringAttribute{
			MarkdownDescription: "Virtual Wwpn",
			Description:         "Virtual Wwpn",

			Computed: true,
		},
	}
}

func omeOSInfoSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{

		"id": schema.Int64Attribute{
			MarkdownDescription: "ID",
			Description:         "ID",

			Computed: true,
		},

		"os_name": schema.StringAttribute{
			MarkdownDescription: "Os Name",
			Description:         "Os Name",

			Computed: true,
		},

		"os_version": schema.StringAttribute{
			MarkdownDescription: "Os Version",
			Description:         "Os Version",

			Computed: true,
		},

		"hostname": schema.StringAttribute{
			MarkdownDescription: "Hostname",
			Description:         "Hostname",

			Computed: true,
		},
	}
}

func omePowerSupplyInfoSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{

		"id": schema.Int64Attribute{
			MarkdownDescription: "ID",
			Description:         "ID",

			Computed: true,
		},

		"name": schema.StringAttribute{
			MarkdownDescription: "Name",
			Description:         "Name",

			Computed: true,
		},

		"power_supply_type": schema.Int64Attribute{
			MarkdownDescription: "Power Supply Type",
			Description:         "Power Supply Type",

			Computed: true,
		},

		"output_watts": schema.Int64Attribute{
			MarkdownDescription: "Output Watts",
			Description:         "Output Watts",

			Computed: true,
		},

		"location": schema.StringAttribute{
			MarkdownDescription: "Location",
			Description:         "Location",

			Computed: true,
		},

		"redundancy_state": schema.StringAttribute{
			MarkdownDescription: "Redundancy State",
			Description:         "Redundancy State",

			Computed: true,
		},

		"status": schema.Int64Attribute{
			MarkdownDescription: "Status",
			Description:         "Status",

			Computed: true,
		},

		"state": schema.StringAttribute{
			MarkdownDescription: "State",
			Description:         "State",

			Computed: true,
		},

		"firmware_version": schema.StringAttribute{
			MarkdownDescription: "Firmware Version",
			Description:         "Firmware Version",

			Computed: true,
		},

		"input_voltage": schema.Int64Attribute{
			MarkdownDescription: "Input Voltage",
			Description:         "Input Voltage",

			Computed: true,
		},

		"model": schema.StringAttribute{
			MarkdownDescription: "Model",
			Description:         "Model",

			Computed: true,
		},

		"manufacturer": schema.StringAttribute{
			MarkdownDescription: "Manufacturer",
			Description:         "Manufacturer",

			Computed: true,
		},

		"range1_max_input_power_watts": schema.Int64Attribute{
			MarkdownDescription: "Range1Max Input Power Watts",
			Description:         "Range1Max Input Power Watts",

			Computed: true,
		},

		"serial_number": schema.StringAttribute{
			MarkdownDescription: "Serial Number",
			Description:         "Serial Number",

			Computed: true,
		},

		"active_input_voltage": schema.StringAttribute{
			MarkdownDescription: "Active Input Voltage",
			Description:         "Active Input Voltage",

			Computed: true,
		},

		"input_power_units": schema.StringAttribute{
			MarkdownDescription: "Input Power Units",
			Description:         "Input Power Units",

			Computed: true,
		},

		"operational_status": schema.StringAttribute{
			MarkdownDescription: "Operational Status",
			Description:         "Operational Status",

			Computed: true,
		},

		"range1_max_input_voltage_high_milli_volts": schema.Int64Attribute{
			MarkdownDescription: "Range1Max Input Voltage High Milli Volts",
			Description:         "Range1Max Input Voltage High Milli Volts",

			Computed: true,
		},

		"rated_max_output_power": schema.Int64Attribute{
			MarkdownDescription: "Rated Max Output Power",
			Description:         "Rated Max Output Power",

			Computed: true,
		},

		"requested_state": schema.Int64Attribute{
			MarkdownDescription: "Requested State",
			Description:         "Requested State",

			Computed: true,
		},

		"ac_input": schema.BoolAttribute{
			MarkdownDescription: "Ac Input",
			Description:         "Ac Input",

			Computed: true,
		},

		"ac_output": schema.BoolAttribute{
			MarkdownDescription: "Ac Output",
			Description:         "Ac Output",

			Computed: true,
		},

		"switching_supply": schema.BoolAttribute{
			MarkdownDescription: "Switching Supply",
			Description:         "Switching Supply",

			Computed: true,
		},
	}
}

func omeDiskInfoSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{

		"id": schema.Int64Attribute{
			MarkdownDescription: "ID",
			Description:         "ID",

			Computed: true,
		},

		"disk_number": schema.StringAttribute{
			MarkdownDescription: "Disk Number",
			Description:         "Disk Number",

			Computed: true,
		},

		"vendor_name": schema.StringAttribute{
			MarkdownDescription: "Vendor Name",
			Description:         "Vendor Name",

			Computed: true,
		},

		"status": schema.Int64Attribute{
			MarkdownDescription: "Status",
			Description:         "Status",

			Computed: true,
		},

		"status_string": schema.StringAttribute{
			MarkdownDescription: "Status String",
			Description:         "Status String",

			Computed: true,
		},

		"model_number": schema.StringAttribute{
			MarkdownDescription: "Model Number",
			Description:         "Model Number",

			Computed: true,
		},

		"serial_number": schema.StringAttribute{
			MarkdownDescription: "Serial Number",
			Description:         "Serial Number",

			Computed: true,
		},

		"sas_address": schema.StringAttribute{
			MarkdownDescription: "Sas Address",
			Description:         "Sas Address",

			Computed: true,
		},

		"revision": schema.StringAttribute{
			MarkdownDescription: "Revision",
			Description:         "Revision",

			Computed: true,
		},

		"manufactured_day": schema.Int64Attribute{
			MarkdownDescription: "Manufactured Day",
			Description:         "Manufactured Day",

			Computed: true,
		},

		"manufactured_week": schema.Int64Attribute{
			MarkdownDescription: "Manufactured Week",
			Description:         "Manufactured Week",

			Computed: true,
		},

		"manufactured_year": schema.Int64Attribute{
			MarkdownDescription: "Manufactured Year",
			Description:         "Manufactured Year",

			Computed: true,
		},

		"encryption_ability": schema.BoolAttribute{
			MarkdownDescription: "Encryption Ability",
			Description:         "Encryption Ability",

			Computed: true,
		},

		"form_factor": schema.StringAttribute{
			MarkdownDescription: "Form Factor",
			Description:         "Form Factor",

			Computed: true,
		},

		"part_number": schema.StringAttribute{
			MarkdownDescription: "Part Number",
			Description:         "Part Number",

			Computed: true,
		},

		"predictive_failure_state": schema.StringAttribute{
			MarkdownDescription: "Predictive Failure State",
			Description:         "Predictive Failure State",

			Computed: true,
		},

		"enclosure_id": schema.StringAttribute{
			MarkdownDescription: "Enclosure ID",
			Description:         "Enclosure ID",

			Computed: true,
		},

		"channel": schema.Int64Attribute{
			MarkdownDescription: "Channel",
			Description:         "Channel",

			Computed: true,
		},

		"size": schema.StringAttribute{
			MarkdownDescription: "Size",
			Description:         "Size",

			Computed: true,
		},

		"free_space": schema.StringAttribute{
			MarkdownDescription: "Free Space",
			Description:         "Free Space",

			Computed: true,
		},

		"used_space": schema.StringAttribute{
			MarkdownDescription: "Used Space",
			Description:         "Used Space",

			Computed: true,
		},

		"bus_type": schema.StringAttribute{
			MarkdownDescription: "Bus Type",
			Description:         "Bus Type",

			Computed: true,
		},

		"slot_number": schema.Int64Attribute{
			MarkdownDescription: "Slot Number",
			Description:         "Slot Number",

			Computed: true,
		},

		"media_type": schema.StringAttribute{
			MarkdownDescription: "Media Type",
			Description:         "Media Type",

			Computed: true,
		},

		"remaining_read_write_endurance": schema.StringAttribute{
			MarkdownDescription: "Remaining Read Write Endurance",
			Description:         "Remaining Read Write Endurance",

			Computed: true,
		},

		"security_state": schema.StringAttribute{
			MarkdownDescription: "Security State",
			Description:         "Security State",

			Computed: true,
		},

		"raid_status": schema.StringAttribute{
			MarkdownDescription: "Raid Status",
			Description:         "Raid Status",

			Computed: true,
		},
	}
}

func omeServerVirtualDiskSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{

		"id": schema.Int64Attribute{
			MarkdownDescription: "ID",
			Description:         "ID",

			Computed: true,
		},

		"raid_controller_id": schema.Int64Attribute{
			MarkdownDescription: "Raid Controller ID",
			Description:         "Raid Controller ID",

			Computed: true,
		},

		"device_id": schema.Int64Attribute{
			MarkdownDescription: "Device ID",
			Description:         "Device ID",

			Computed: true,
		},

		"fqdd": schema.StringAttribute{
			MarkdownDescription: "Fqdd",
			Description:         "Fqdd",

			Computed: true,
		},

		"state": schema.StringAttribute{
			MarkdownDescription: "State",
			Description:         "State",

			Computed: true,
		},

		"rollup_status": schema.Int64Attribute{
			MarkdownDescription: "Rollup Status",
			Description:         "Rollup Status",

			Computed: true,
		},

		"status": schema.Int64Attribute{
			MarkdownDescription: "Status",
			Description:         "Status",

			Computed: true,
		},

		"layout": schema.StringAttribute{
			MarkdownDescription: "Layout",
			Description:         "Layout",

			Computed: true,
		},

		"media_type": schema.StringAttribute{
			MarkdownDescription: "Media Type",
			Description:         "Media Type",

			Computed: true,
		},

		"name": schema.StringAttribute{
			MarkdownDescription: "Name",
			Description:         "Name",

			Computed: true,
		},

		"read_policy": schema.StringAttribute{
			MarkdownDescription: "Read Policy",
			Description:         "Read Policy",

			Computed: true,
		},

		"write_policy": schema.StringAttribute{
			MarkdownDescription: "Write Policy",
			Description:         "Write Policy",

			Computed: true,
		},

		"cache_policy": schema.StringAttribute{
			MarkdownDescription: "Cache Policy",
			Description:         "Cache Policy",

			Computed: true,
		},

		"stripe_size": schema.StringAttribute{
			MarkdownDescription: "Stripe Size",
			Description:         "Stripe Size",

			Computed: true,
		},

		"size": schema.StringAttribute{
			MarkdownDescription: "Size",
			Description:         "Size",

			Computed: true,
		},

		"target_id": schema.Int64Attribute{
			MarkdownDescription: "Target ID",
			Description:         "Target ID",

			Computed: true,
		},

		"lock_status": schema.StringAttribute{
			MarkdownDescription: "Lock Status",
			Description:         "Lock Status",

			Computed: true,
		},
	}
}

func omeRAIDControllerInfoSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{

		"id": schema.Int64Attribute{
			MarkdownDescription: "ID",
			Description:         "ID",

			Computed: true,
		},

		"name": schema.StringAttribute{
			MarkdownDescription: "Name",
			Description:         "Name",

			Computed: true,
		},

		"fqdd": schema.StringAttribute{
			MarkdownDescription: "Fqdd",
			Description:         "Fqdd",

			Computed: true,
		},

		"device_description": schema.StringAttribute{
			MarkdownDescription: "Device Description",
			Description:         "Device Description",

			Computed: true,
		},

		"status": schema.Int64Attribute{
			MarkdownDescription: "Status",
			Description:         "Status",

			Computed: true,
		},

		"status_type": schema.StringAttribute{
			MarkdownDescription: "Status Type",
			Description:         "Status Type",

			Computed: true,
		},

		"rollup_status": schema.Int64Attribute{
			MarkdownDescription: "Rollup Status",
			Description:         "Rollup Status",

			Computed: true,
		},

		"rollup_status_string": schema.StringAttribute{
			MarkdownDescription: "Rollup Status String",
			Description:         "Rollup Status String",

			Computed: true,
		},

		"firmware_version": schema.StringAttribute{
			MarkdownDescription: "Firmware Version",
			Description:         "Firmware Version",

			Computed: true,
		},

		"cache_size_in_mb": schema.Int64Attribute{
			MarkdownDescription: "Cache Size In Mb",
			Description:         "Cache Size In Mb",

			Computed: true,
		},

		"pci_slot": schema.StringAttribute{
			MarkdownDescription: "Pci Slot",
			Description:         "Pci Slot",

			Computed: true,
		},

		"driver_version": schema.StringAttribute{
			MarkdownDescription: "Driver Version",
			Description:         "Driver Version",

			Computed: true,
		},

		"storage_assignment_allowed": schema.StringAttribute{
			MarkdownDescription: "Storage Assignment Allowed",
			Description:         "Storage Assignment Allowed",

			Computed: true,
		},

		"server_virtual_disks": schema.ListNestedAttribute{
			MarkdownDescription: "Server Virtual Disks",
			Description:         "Server Virtual Disks",

			Computed:     true,
			NestedObject: schema.NestedAttributeObject{Attributes: omeServerVirtualDiskSchema()},
		},
	}
}

func omeMemoryInfoSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{

		"id": schema.Int64Attribute{
			MarkdownDescription: "ID",
			Description:         "ID",

			Computed: true,
		},

		"name": schema.StringAttribute{
			MarkdownDescription: "Name",
			Description:         "Name",

			Computed: true,
		},

		"bank_name": schema.StringAttribute{
			MarkdownDescription: "Bank Name",
			Description:         "Bank Name",

			Computed: true,
		},

		"size": schema.Int64Attribute{
			MarkdownDescription: "Size",
			Description:         "Size",

			Computed: true,
		},

		"status": schema.Int64Attribute{
			MarkdownDescription: "Status",
			Description:         "Status",

			Computed: true,
		},

		"manufacturer": schema.StringAttribute{
			MarkdownDescription: "Manufacturer",
			Description:         "Manufacturer",

			Computed: true,
		},

		"part_number": schema.StringAttribute{
			MarkdownDescription: "Part Number",
			Description:         "Part Number",

			Computed: true,
		},

		"serial_number": schema.StringAttribute{
			MarkdownDescription: "Serial Number",
			Description:         "Serial Number",

			Computed: true,
		},

		"type_details": schema.StringAttribute{
			MarkdownDescription: "Type Details",
			Description:         "Type Details",

			Computed: true,
		},

		"manufacturer_date": schema.StringAttribute{
			MarkdownDescription: "Manufacturer Date",
			Description:         "Manufacturer Date",

			Computed: true,
		},

		"speed": schema.Int64Attribute{
			MarkdownDescription: "Speed",
			Description:         "Speed",

			Computed: true,
		},

		"current_operating_speed": schema.Int64Attribute{
			MarkdownDescription: "Current Operating Speed",
			Description:         "Current Operating Speed",

			Computed: true,
		},

		"rank": schema.StringAttribute{
			MarkdownDescription: "Rank",
			Description:         "Rank",

			Computed: true,
		},

		"instance_id": schema.StringAttribute{
			MarkdownDescription: "Instance ID",
			Description:         "Instance ID",

			Computed: true,
		},

		"device_description": schema.StringAttribute{
			MarkdownDescription: "Device Description",
			Description:         "Device Description",

			Computed: true,
		},
	}
}

func omeStorageEnclosureInfoSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{

		"id": schema.Int64Attribute{
			MarkdownDescription: "ID",
			Description:         "ID",

			Computed: true,
		},

		"name": schema.StringAttribute{
			MarkdownDescription: "Name",
			Description:         "Name",

			Computed: true,
		},

		"status": schema.Int64Attribute{
			MarkdownDescription: "Status",
			Description:         "Status",

			Computed: true,
		},

		"status_type": schema.StringAttribute{
			MarkdownDescription: "Status Type String",
			Description:         "Status Type String",

			Computed: true,
		},

		"channel_number": schema.StringAttribute{
			MarkdownDescription: "Channel Number",
			Description:         "Channel Number",

			Computed: true,
		},

		"backplane_part_num": schema.StringAttribute{
			MarkdownDescription: "Backplane Part Num",
			Description:         "Backplane Part Num",

			Computed: true,
		},

		"number_of_fan_packs": schema.Int64Attribute{
			MarkdownDescription: "Number Of Fan Packs",
			Description:         "Number Of Fan Packs",

			Computed: true,
		},

		"version": schema.StringAttribute{
			MarkdownDescription: "Version",
			Description:         "Version",

			Computed: true,
		},

		"rollup_status": schema.Int64Attribute{
			MarkdownDescription: "Rollup Status",
			Description:         "Rollup Status",

			Computed: true,
		},

		"slot_count": schema.Int64Attribute{
			MarkdownDescription: "Slot Count",
			Description:         "Slot Count",

			Computed: true,
		},
	}
}

func omeServerPowerStateSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{

		"id": schema.Int64Attribute{
			MarkdownDescription: "ID",
			Description:         "ID",

			Computed: true,
		},

		"power_state": schema.Int64Attribute{
			MarkdownDescription: "Power State",
			Description:         "Power State",

			Computed: true,
		},
	}
}

func omeDeviceLicenseSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{

		"sold_date": schema.StringAttribute{
			MarkdownDescription: "Sold Date",
			Description:         "Sold Date",

			Computed: true,
		},

		"license_bound": schema.Int64Attribute{
			MarkdownDescription: "License Bound",
			Description:         "License Bound",

			Computed: true,
		},

		"eval_time_remaining": schema.Int64Attribute{
			MarkdownDescription: "Eval Time Remaining",
			Description:         "Eval Time Remaining",

			Computed: true,
		},

		"assigned_devices": schema.StringAttribute{
			MarkdownDescription: "Assigned Devices",
			Description:         "Assigned Devices",

			Computed: true,
		},

		"license_status": schema.Int64Attribute{
			MarkdownDescription: "License Status",
			Description:         "License Status",

			Computed: true,
		},

		"entitlement_id": schema.StringAttribute{
			MarkdownDescription: "Entitlement Id",
			Description:         "Entitlement Id",

			Computed: true,
		},

		"license_description": schema.StringAttribute{
			MarkdownDescription: "License Description",
			Description:         "License Description",

			Computed: true,
		},

		"license_type": schema.SingleNestedAttribute{
			MarkdownDescription: "License Type",
			Description:         "License Type",

			Computed:   true,
			Attributes: omeLicenseTypeSchema(),
		},
	}
}

func omeLicenseTypeSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{

		"name": schema.StringAttribute{
			MarkdownDescription: "Name",
			Description:         "Name",

			Computed: true,
		},

		"license_id": schema.Int64Attribute{
			MarkdownDescription: "License Id",
			Description:         "License Id",

			Computed: true,
		},
	}
}
