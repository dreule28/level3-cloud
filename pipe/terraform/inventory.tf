locals {
  node_ssh_ip = {
    for name, inst in openstack_compute_instance_v2.vm :
    name => openstack_networking_floatingip_v2.fip[name].address
  }

  control_plane_ips = { for n in local.control_plane_names : n => local.node_ssh_ip[n] }
  worker_ips        = { for n in local.worker_names : n => local.node_ssh_ip[n] }
}


resource "local_file" "ansible_inventory" {
  filename = "${path.module}/../ansible/inventory.ini"
  content = templatefile("${path.module}/templates/inventory.tftpl", {
    control_plane            = local.control_plane_ips
    workers                  = local.worker_ips
    ansible_user             = var.ansible_user
    ansible_private_key_file = var.ansible_private_key_file
  })
}
