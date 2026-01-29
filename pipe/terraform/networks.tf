# External network (Floating IP network)
data "openstack_networking_network_v2" "external" {
  name = var.external_network_name
}

# Tenant/private network (for your VMs)
resource "openstack_networking_network_v2" "tenant" {
  name           = "k8s-net"
  admin_state_up = true
}

# Tenant subnet (private RFC1918 range)
resource "openstack_networking_subnet_v2" "tenant" {
  name            = "k8s-subnet"
  network_id      = openstack_networking_network_v2.tenant.id
  cidr            = "10.10.0.0/24"
  ip_version      = 4
  dns_nameservers = ["1.1.1.1", "8.8.8.8"]
}

# Router (connects tenant network to external network)
resource "openstack_networking_router_v2" "router" {
  name                = "k8s-router"
  external_network_id = data.openstack_networking_network_v2.external.id
}

# Attach tenant subnet to router
resource "openstack_networking_router_interface_v2" "router_interface" {
  router_id = openstack_networking_router_v2.router.id
  subnet_id = openstack_networking_subnet_v2.tenant.id
}
