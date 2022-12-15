variable "username" {
  type = string
  description = "stores the username of ome."
}

variable "password" {
  type = string
  description = "stores the password of ome."
}

variable "host" {
  type = string
  description = "stores the host address of ome instance."
}

variable "skipssl" {
  type = bool
  description = "specifies if the ssl verification needs to be skipped."
}

variable "baselinename" {
  type = string
  description = "baseline name"
}