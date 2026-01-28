output "inventory_path" {
  value = local_file.ansible_inventory.filename
}

output "node_ssh_ip" {
  value = local.node_ssh_ip
}

output "floating_ips" {
  value = {
    for name, fip in openstack_networking_floatingip_v2.fip :
    name => fip.address
  }
}
