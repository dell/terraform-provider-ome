terraform {
  required_providers {
    ome = {
      version = "0.0.1"
      source  = "registry.terraform.io/dell/ome"
    }
  }
}

provider "ome" {
  username = "username"
  password = "password"
  host     = "yourhost.host.com"
  skipssl  = false
}