resource "openstack_networking_floatingip_v2" "fip" {
  for_each = openstack_compute_instance_v2.vm
  pool     = var.floating_ip_pool
}

resource "openstack_compute_floatingip_associate_v2" "fip_assoc" {
  for_each    = openstack_compute_instance_v2.vm
  floating_ip = openstack_networking_floatingip_v2.fip[each.key].address
  instance_id = each.value.id
}
