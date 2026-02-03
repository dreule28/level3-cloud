resource "openstack_networking_secgroup_v2" "ssh" {
  name        = "ssh"
  description = "SSH access for Ansible"
}

resource "openstack_networking_secgroup_rule_v2" "ssh_ingress" {
  security_group_id = openstack_networking_secgroup_v2.ssh.id
  direction         = "ingress"
  ethertype         = "IPv4"
  protocol          = "tcp"
  port_range_min    = 22
  port_range_max    = 22
  remote_ip_prefix  = var.ssh_allowed_cidr
}

resource "openstack_networking_secgroup_v2" "k8s_internal" {
  name        = "k8s-internal"
  description = "Kubernetes node-to-node internal traffic"
}

resource "openstack_networking_secgroup_rule_v2" "k8s_tcp_in" {
  security_group_id = openstack_networking_secgroup_v2.k8s_internal.id
  direction         = "ingress"
  ethertype         = "IPv4"
  protocol          = "tcp"
  remote_group_id   = openstack_networking_secgroup_v2.k8s_internal.id
}

resource "openstack_networking_secgroup_rule_v2" "k8s_udp_in" {
  security_group_id = openstack_networking_secgroup_v2.k8s_internal.id
  direction         = "ingress"
  ethertype         = "IPv4"
  protocol          = "udp"
  remote_group_id   = openstack_networking_secgroup_v2.k8s_internal.id
}

resource "openstack_networking_secgroup_rule_v2" "k8s_icmp_in" {
  security_group_id = openstack_networking_secgroup_v2.k8s_internal.id
  direction         = "ingress"
  ethertype         = "IPv4"
  protocol          = "icmp"
  remote_group_id   = openstack_networking_secgroup_v2.k8s_internal.id
}
