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
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// DiscoveryJobSchema for discovery job schema.
func DiscoveryJobSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{

		"id": schema.StringAttribute{
			MarkdownDescription: "ID of the discovery configuration group",
			Description:         "ID of the discovery configuration group",
			Computed:            true,
		},

		"name": schema.StringAttribute{
			MarkdownDescription: "Name of the discovery configuration job",
			Description:         "Name of the discovery configuration job",
			Required:            true,
			Validators: []validator.String{
				stringvalidator.LengthAtLeast(1),
			},
		},

		"email_recipient": schema.StringAttribute{
			MarkdownDescription: `
				- Enter the email address to which notifications are to be sent about the discovery job status.
				- Configure the SMTP settings to allow sending notifications to an email address.`,
			Description: `
				- Enter the email address to which notifications are to be sent about the discovery job status.
				- Configure the SMTP settings to allow sending notifications to an email address.`,
			Optional: true,
			Validators: []validator.String{
				stringvalidator.LengthAtLeast(1),
			},
		},

		"discovery_config_targets": schema.SetNestedAttribute{
			MarkdownDescription: `
				- Provide the list of discovery targets.
      			- Each discovery target is a set of "network_address_detail", "device_types", and one or more protocol credentials.`,
			Description: `
				- Provide the list of discovery targets.
      			- Each discovery target is a set of "network_address_detail", "device_types", and one or more protocol credentials.`,
			Required:     true,
			NestedObject: schema.NestedAttributeObject{Attributes: DiscoveryConfigTargetsSchema()},
			Validators: []validator.Set{
				setvalidator.SizeAtLeast(1),
			},
		},

		"schedule": schema.StringAttribute{
			MarkdownDescription: "Provides the option to schedule the discovery job. If `RunLater` is selected, then attribute `cron` must be specified. If `RunNow` is selected, then attribute `timeout` must be specified.",
			Description:         "Provides the option to schedule the discovery job. If `RunLater` is selected, then attribute `cron` must be specified. If `RunNow` is selected, then attribute `timeout` must be specified.",
			Required:            true,
			Validators: []validator.String{stringvalidator.OneOf(
				"RunNow",
				"RunLater",
			)},
			PlanModifiers: []planmodifier.String{
				StringDefaultValue(types.StringValue("RunNow")),
			},
		},

		"cron": schema.StringAttribute{
			MarkdownDescription: "Provide a cron expression based on Quartz cron format",
			Description:         "Provide a cron expression based on Quartz cron format",
			Optional:            true,
			Computed:            true,
			Validators: []validator.String{
				stringvalidator.LengthAtLeast(1),
			},
		},

		"timeout": schema.Int64Attribute{
			MarkdownDescription: "Provide a timeout in minute to track the job",
			Description:         "Provide a timeout in minute to track the job",
			Optional:            true,
			Validators: []validator.Int64{
				int64validator.AtLeast(1),
			},
		},

		"ignore_partial_failure": schema.BoolAttribute{
			MarkdownDescription: "Provides the option to ignore partial failures. Partial failures occur when there is a combination of both discovered and undiscovered IPs with Schedule is set to `RunNow`. If `partial_failure` is set `false` then partial_failure is not ignored, and module will error out.If `partial_failure` is set `true` then partial_failure is ignored, and module will not error out.",
			Description:         "Provides the option to ignore partial failures. Partial failures occur when there is a combination of both discovered and undiscovered IPs with Schedule is set to `RunNow`. If `partial_failure` is set `false` then partial_failure is not ignored, and module will error out.If `partial_failure` is set `true` then partial_failure is ignored, and module will not error out.",
			Optional:            true,
		},

		"trap_destination": schema.BoolAttribute{
			MarkdownDescription: `
				- Enable OpenManage Enterprise to receive the incoming SNMP traps from the discovered devices. 
				- This is effective only for servers discovered by using their iDRAC interface.`,
			Description: `
				- Enable OpenManage Enterprise to receive the incoming SNMP traps from the discovered devices. 
				- This is effective only for servers discovered by using their iDRAC interface.`,
			Optional: true,
			Computed: true,
			PlanModifiers: []planmodifier.Bool{
				BoolDefaultValue(types.BoolValue(false)),
			},
		},

		"enable_community_strings": schema.BoolAttribute{
			MarkdownDescription: `
				- Enable the use of SNMP community strings to receive SNMP traps using Application Settings in OpenManage Enterprise. 
				- This option is available only for the discovered iDRAC servers and MX7000 chassis.`,
			Description: `
				- Enable the use of SNMP community strings to receive SNMP traps using Application Settings in OpenManage Enterprise. 
				- This option is available only for the discovered iDRAC servers and MX7000 chassis.`,
			Optional: true,
			Computed: true,
			PlanModifiers: []planmodifier.Bool{
				BoolDefaultValue(types.BoolValue(false)),
			},
		},
		"job_id": schema.Int64Attribute{
			MarkdownDescription: "Discovery Job ID.",
			Description:         "Discovery Job ID.",
			Computed:            true,
		},
		"job_tracking": schema.SingleNestedAttribute{
			MarkdownDescription: "Discovery Job Tracking Info Captured When Schedule is set to `RunNow`",
			Description:         "Discovery Job Tracking Info Captured When Schedule is set to `RunNow`",
			Computed:            true,
			Attributes:          JobSchema(),
		},
	}
}

// JobSchema for collecting info of job tracking.
func JobSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"job_execution_results": schema.ListAttribute{
			MarkdownDescription: "Provides information about job executions discovered after tracking the job until it timeout.",
			Description:         "Provides information about job executions discovered after tracking the job until it timeout.",
			Computed:            true,
			ElementType:         types.StringType,
		},
		"discovered_ip": schema.ListAttribute{
			MarkdownDescription: "IPs discovered after tracking the job until it timeout.",
			Description:         "IPs discovered after tracking the job until it timeout.",
			Computed:            true,
			ElementType:         types.StringType,
		},
		"undiscovered_ip": schema.ListAttribute{
			MarkdownDescription: "IPs remains undiscovered after tracking the job until it timeout.",
			Description:         "IPs remains undiscovered after tracking the job until it timeout.",
			Computed:            true,
			ElementType:         types.StringType,
		},
	}
}

// DiscoveryConfigTargetsSchema for discovery config target schema.
func DiscoveryConfigTargetsSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{

		"network_address_detail": schema.ListAttribute{
			MarkdownDescription: `
				- "Provide the list of IP addresses, host names, or the range of IP addresses of the devices to be discoveredor included."
         		- "Sample Valid IP Range Formats"
         		- "   192.35.0.0"
         		- "   192.36.0.0-10.36.0.255"
         		- "   192.37.0.0/24"
         		- "   2345:f2b1:f083:135::5500/118"
         		- "   2345:f2b1:f083:135::a500-2607:f2b1:f083:135::a600"
         		- "   hostname.domain.tld"
         		- "   hostname"
         		- "   2345:f2b1:f083:139::22a"
         		- "Sample Invalid IP Range Formats"
         		- "   192.35.0.*"
         		- "   192.36.0.0-255"
         		- "   192.35.0.0/255.255.255.0"
         		- NOTE: The range size for the number of IP addresses is limited to 16,385 (0x4001).
         		- NOTE: Both IPv6 and IPv6 CIDR formats are supported.`,
			Description: `
				- "Provide the list of IP addresses, host names, or the range of IP addresses of the devices to be discoveredor included."
         		- "Sample Valid IP Range Formats"
         		- "   192.35.0.0"
         		- "   192.36.0.0-10.36.0.255"
         		- "   192.37.0.0/24"
         		- "   2345:f2b1:f083:135::5500/118"
         		- "   2345:f2b1:f083:135::a500-2607:f2b1:f083:135::a600"
         		- "   hostname.domain.tld"
         		- "   hostname"
         		- "   2345:f2b1:f083:139::22a"
         		- "Sample Invalid IP Range Formats"
         		- "   192.35.0.*"
         		- "   192.36.0.0-255"
         		- "   192.35.0.0/255.255.255.0"
         		- NOTE: The range size for the number of IP addresses is limited to 16,385 (0x4001).
         		- NOTE: Both IPv6 and IPv6 CIDR formats are supported.`,
			Required:    true,
			ElementType: types.StringType,
			Validators: []validator.List{
				listvalidator.SizeAtLeast(1),
			},
		},

		"device_type": schema.ListAttribute{
			MarkdownDescription: `
				- Provide the type of devices to be discovered.
				- The accepted types are SERVER, CHASSIS, NETWORK SWITCH, and STORAGE.
				- A combination or all of the above can be provided.
				- "Supported protocols for each device type are:"
				- SERVER - "redfish", "snmp", and "ssh".
				- CHASSIS - "redfish".
				- NETWORK SWITCH - "snmp".
				- STORAGE - "snmp".`,
			Description: `
			- Provide the type of devices to be discovered.
			- The accepted types are SERVER, CHASSIS, NETWORK SWITCH, and STORAGE.
			- A combination or all of the above can be provided.
			- "Supported protocols for each device type are:"
			- SERVER - "redfish", "snmp", and "ssh".
			- CHASSIS - "redfish".
			- NETWORK SWITCH - "snmp".
			- STORAGE - "snmp".`,
			Required:    true,
			ElementType: types.StringType,
			Validators: []validator.List{
				listvalidator.SizeAtLeast(1),
				listvalidator.ValueStringsAre(
					stringvalidator.OneOf("SERVER", "CHASSIS", "NETWORK SWITCH", "STORAGE"),
				),
				listvalidator.UniqueValues(),
			},
		},

		"redfish": schema.SingleNestedAttribute{
			MarkdownDescription: "REDFISH protocol",
			Description:         "REDFISH protocol",
			Optional:            true,
			Attributes:          RedfishSchema(),
		},

		"wsman": schema.SingleNestedAttribute{
			MarkdownDescription: "WSMAN protocol",
			Description:         "WSMAN protocol",
			Optional:            true,
			Attributes:          WSMANSchema(),
		},

		"snmp": schema.SingleNestedAttribute{
			MarkdownDescription: "Simple Network Management Protocol (SNMP)",
			Description:         "Simple Network Management Protocol (SNMP)",
			Optional:            true,
			Attributes:          SNMPSchema(),
		},

		"ssh": schema.SingleNestedAttribute{
			MarkdownDescription: "Secure Shell (SSH)",
			Description:         "Secure Shell (SSH)",
			Optional:            true,
			Attributes:          SSHSchema(),
		},
	}
}

// RedfishSchema for redfish protocol schema
func RedfishSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{

		"username": schema.StringAttribute{
			MarkdownDescription: "Provide a username for the protocol.",
			Description:         "Provide a username for the protocol.",
			Required:            true,
			Validators: []validator.String{
				stringvalidator.LengthAtLeast(1),
			},
		},

		"password": schema.StringAttribute{
			MarkdownDescription: "Provide a password for the protocol.",
			Description:         "Provide a password for the protocol.",
			Required:            true,
			Validators: []validator.String{
				stringvalidator.LengthAtLeast(1),
			},
		},

		"port": schema.Int64Attribute{
			MarkdownDescription: "Enter the port number that the job must use to discover the devices.",
			Description:         "Enter the port number that the job must use to discover the devices.",
			Optional:            true,
			Computed:            true,
			PlanModifiers: []planmodifier.Int64{
				Int64DefaultValue(types.Int64Value(443)),
			},
		},

		"retries": schema.Int64Attribute{
			MarkdownDescription: "Enter the number of repeated attempts required to discover a device",
			Description:         "Enter the number of repeated attempts required to discover a device",
			Optional:            true,
			Computed:            true,
			PlanModifiers: []planmodifier.Int64{
				Int64DefaultValue(types.Int64Value(3)),
			},
		},

		"timeout": schema.Int64Attribute{
			MarkdownDescription: "Enter the time in seconds after which a job must stop running.",
			Description:         "Enter the time in seconds after which a job must stop running.",
			Optional:            true,
			Computed:            true,
			PlanModifiers: []planmodifier.Int64{
				Int64DefaultValue(types.Int64Value(60)),
			},
		},

		"cn_check": schema.BoolAttribute{
			MarkdownDescription: "Enable the Common Name (CN) check.",
			Description:         "Enable the Common Name (CN) check.",
			Optional:            true,
			Computed:            true,
			PlanModifiers: []planmodifier.Bool{
				BoolDefaultValue(types.BoolValue(false)),
			},
		},

		"ca_check": schema.BoolAttribute{
			MarkdownDescription: "Enable the Certificate Authority (CA) check.",
			Description:         "Enable the Certificate Authority (CA) check.",
			Optional:            true,
			Computed:            true,
			PlanModifiers: []planmodifier.Bool{
				BoolDefaultValue(types.BoolValue(false)),
			},
		},

		// "domain": schema.StringAttribute{
		// 	MarkdownDescription: "Provide a domain for the protocol.",
		// 	Description:         "Provide a domain for the protocol.",
		// 	Optional:            true,
		// 	Computed:            true,
		// 	Validators: []validator.String{
		// 		stringvalidator.LengthAtLeast(1),
		// 	},
		// },

		// "certificate_data": schema.StringAttribute{
		// 	MarkdownDescription: "Provide certificate data for the CA check.",
		// 	Description:         "Provide certificate data for the CA check.",
		// 	Optional:            true,
		// 	Computed:            true,
		// 	Validators: []validator.String{
		// 		stringvalidator.LengthAtLeast(1),
		// 	},
		// },
	}
}

// WSMANSchema for wsman protocol schema
func WSMANSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{

		"username": schema.StringAttribute{
			MarkdownDescription: "Provide a username for the protocol.",
			Description:         "Provide a username for the protocol.",
			Required:            true,
			Validators: []validator.String{
				stringvalidator.LengthAtLeast(1),
			},
		},

		"password": schema.StringAttribute{
			MarkdownDescription: "Provide a password for the protocol.",
			Description:         "Provide a password for the protocol.",
			Required:            true,
			Validators: []validator.String{
				stringvalidator.LengthAtLeast(1),
			},
		},

		"port": schema.Int64Attribute{
			MarkdownDescription: "Enter the port number that the job must use to discover the devices.",
			Description:         "Enter the port number that the job must use to discover the devices.",
			Optional:            true,
			Computed:            true,
			PlanModifiers: []planmodifier.Int64{
				Int64DefaultValue(types.Int64Value(443)),
			},
		},

		"retries": schema.Int64Attribute{
			MarkdownDescription: "Enter the number of repeated attempts required to discover a device",
			Description:         "Enter the number of repeated attempts required to discover a device",
			Optional:            true,
			Computed:            true,
			PlanModifiers: []planmodifier.Int64{
				Int64DefaultValue(types.Int64Value(3)),
			},
		},

		"timeout": schema.Int64Attribute{
			MarkdownDescription: "Enter the time in seconds after which a job must stop running.",
			Description:         "Enter the time in seconds after which a job must stop running.",
			Optional:            true,
			Computed:            true,
			PlanModifiers: []planmodifier.Int64{
				Int64DefaultValue(types.Int64Value(60)),
			},
		},

		"cn_check": schema.BoolAttribute{
			MarkdownDescription: "Enable the Common Name (CN) check.",
			Description:         "Enable the Common Name (CN) check.",
			Optional:            true,
			Computed:            true,
			PlanModifiers: []planmodifier.Bool{
				BoolDefaultValue(types.BoolValue(false)),
			},
		},

		"ca_check": schema.BoolAttribute{
			MarkdownDescription: "Enable the Certificate Authority (CA) check.",
			Description:         "Enable the Certificate Authority (CA) check.",
			Optional:            true,
			Computed:            true,
			PlanModifiers: []planmodifier.Bool{
				BoolDefaultValue(types.BoolValue(false)),
			},
		},

		// "domain": schema.StringAttribute{
		// 	MarkdownDescription: "Provide a domain for the protocol.",
		// 	Description:         "Provide a domain for the protocol.",
		// 	Optional:            true,
		// 	Computed:            true,
		// 	Validators: []validator.String{
		// 		stringvalidator.LengthAtLeast(1),
		// 	},
		// },

		// "certificate_data": schema.StringAttribute{
		// 	MarkdownDescription: "Provide certificate data for the CA check.",
		// 	Description:         "Provide certificate data for the CA check.",
		// 	Optional:            true,
		// 	Computed:            true,
		// 	Validators: []validator.String{
		// 		stringvalidator.LengthAtLeast(1),
		// 	},
		// },
	}
}

// SNMPSchema for SNMP protocol schema
func SNMPSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{

		"community": schema.StringAttribute{
			MarkdownDescription: "Community string for the SNMP protocol.",
			Description:         "Community string for the SNMP protocol.",
			Required:            true,
			Validators: []validator.String{
				stringvalidator.LengthAtLeast(1),
			},
		},

		"port": schema.Int64Attribute{
			MarkdownDescription: "Enter the port number that the job must use to discover the devices.",
			Description:         "Enter the port number that the job must use to discover the devices.",
			Optional:            true,
			Computed:            true,
			PlanModifiers: []planmodifier.Int64{
				Int64DefaultValue(types.Int64Value(161)),
			},
		},

		"retries": schema.Int64Attribute{
			MarkdownDescription: "Enter the number of repeated attempts required to discover a device.",
			Description:         "Enter the number of repeated attempts required to discover a device.",
			Optional:            true,
			Computed:            true,
			PlanModifiers: []planmodifier.Int64{
				Int64DefaultValue(types.Int64Value(3)),
			},
		},

		"timeout": schema.Int64Attribute{
			MarkdownDescription: "Enter the time in seconds after which a job must stop running.",
			Description:         "Enter the time in seconds after which a job must stop running.",
			Optional:            true,
			Computed:            true,
			PlanModifiers: []planmodifier.Int64{
				Int64DefaultValue(types.Int64Value(3)),
			},
		},
	}
}

// SSHSchema for SSH protocol schema
func SSHSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{

		"username": schema.StringAttribute{
			MarkdownDescription: "Provide a username for the protocol.",
			Description:         "Provide a username for the protocol.",
			Required:            true,
			Validators: []validator.String{
				stringvalidator.LengthAtLeast(1),
			},
		},

		"password": schema.StringAttribute{
			MarkdownDescription: "Provide a password for the protocol.",
			Description:         "Provide a password for the protocol.",
			Required:            true,
			Validators: []validator.String{
				stringvalidator.LengthAtLeast(1),
			},
		},

		"port": schema.Int64Attribute{
			MarkdownDescription: "Enter the port number that the job must use to discover the devices.",
			Description:         "Enter the port number that the job must use to discover the devices.",
			Optional:            true,
			Computed:            true,
			PlanModifiers: []planmodifier.Int64{
				Int64DefaultValue(types.Int64Value(22)),
			},
		},

		"retries": schema.Int64Attribute{
			MarkdownDescription: "Enter the number of repeated attempts required to discover a device.",
			Description:         "Enter the number of repeated attempts required to discover a device.",
			Optional:            true,
			Computed:            true,
			PlanModifiers: []planmodifier.Int64{
				Int64DefaultValue(types.Int64Value(3)),
			},
		},

		"timeout": schema.Int64Attribute{
			MarkdownDescription: "Enter the time in seconds after which a job must stop running.",
			Description:         "Enter the time in seconds after which a job must stop running.",
			Optional:            true,
			Computed:            true,
			PlanModifiers: []planmodifier.Int64{
				Int64DefaultValue(types.Int64Value(60)),
			},
		},

		"check_known_hosts": schema.BoolAttribute{
			MarkdownDescription: "Verify the known host key.",
			Description:         "Verify the known host key.",
			Optional:            true,
			Computed:            true,
			PlanModifiers: []planmodifier.Bool{
				BoolDefaultValue(types.BoolValue(false)),
			},
		},

		"is_sudo_user": schema.BoolAttribute{
			MarkdownDescription: "Use the SUDO option",
			Description:         "Use the SUDO option",
			Optional:            true,
			Computed:            true,
			PlanModifiers: []planmodifier.Bool{
				BoolDefaultValue(types.BoolValue(false)),
			},
		},
	}
}
