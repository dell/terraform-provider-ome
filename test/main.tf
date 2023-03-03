	provider "ome" {
		username = ""
		password = ""
		host = ""
		skipssl = true
	}

	data "ome_vlannetworks_info" "vlans" {
	}