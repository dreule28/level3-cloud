variable "network_name" {
  type    = string
  default = "shared"
}

variable "key_pair" {
  type    = string
  default = "openstack-key"
}

variable "ssh_allowed_cidr" {
  type        = string
  description = "CIDR allowed to SSH into nodes (tighten this!). Example: 1.2.3.4/32"
  default     = "0.0.0.0/0"
}

variable "ansible_user" {
  type        = string
  description = "SSH user Ansible should use (depends on your image)."
  default     = "debian"
}

variable "ansible_private_key_file" {
  type        = string
  description = "Path to SSH private key for Ansible."
  default     = "/opt/stack/.ssh/openstack_ed25519"
}


# The only thing you edit to add/remove nodes:
variable "nodes" {
  description = "Cluster nodes keyed by name."
  type = map(object({
    role = string # control_plane | worker
  }))

  default = {
    server = { role = "control_plane" }
    node-0 = { role = "worker" }
    node-1 = { role = "worker" }
  }
}

variable "floating_ip_pool" {
  type        = string
  description = "Exeternal network name used for FLoating IPs"
  default     = "public"
}

variable "keypair_name" {
  type    = string
  default = "tf-openstack-key"
}

variable "public_key_path" {
  type        = string
  description = "Path to your SSH public key that will be uploaded to OpenStack."
  default     = "/opt/stack/.ssh/openstack_ed25519.pub"
}
