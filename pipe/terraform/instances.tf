locals {
  control_plane_names = [for name, n in var.nodes : name if n.role == "control_plane"]
  worker_names        = [for name, n in var.nodes : name if n.role == "worker"]
}

resource "openstack_compute_instance_v2" "vm" {
  # depends_on = [null_resource.openstack_host_prereqs]

  for_each = var.nodes

  name      = each.key
  flavor_id = data.openstack_compute_flavor_v2.small.id
  image_id  = data.openstack_images_image_v2.bookworm.id
  key_pair  = openstack_compute_keypair_v2.this.name

  config_drive = true

  network {
    name = var.network_name
  }

  security_groups = [
    openstack_networking_secgroup_v2.ssh.name,
    openstack_networking_secgroup_v2.k8s_internal.name,
  ]
}
