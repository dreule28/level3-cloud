locals {
  control_plane_names = [for name, n in var.nodes : name if n.role == "control_plane"]
  worker_names        = [for name, n in var.nodes : name if n.role == "worker"]
}

locals {
  flavor_by_role = {
    control_plane = data.openstack_compute_flavor_v2.medium.id
    worker        = data.openstack_compute_flavor_v2.small.id
  }
}


resource "openstack_compute_instance_v2" "vm" {
  depends_on = [null_resource.openstack_host_prereqs]

  for_each = var.nodes

  name      = each.key
  flavor_id = local.flavor_by_role[each.value.role]
  image_id  = openstack_images_image_v2.jammy.id
  key_pair  = openstack_compute_keypair_v2.this.name

  config_drive = true

  network {
    uuid = openstack_networking_network_v2.tenant.id
  }

  security_groups = [
    openstack_networking_secgroup_v2.ssh.name,
    openstack_networking_secgroup_v2.k8s_internal.name,
  ]
}
