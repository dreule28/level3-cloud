locals {
  # Floating IP used for SSH
  node_ssh_ip = {
    for name, inst in openstack_compute_instance_v2.vm :
    name => openstack_networking_floatingip_v2.fip[name].address
  }

  # Fixed IP used inside the cluster
  node_k8s_ip = {
    for name, inst in openstack_compute_instance_v2.vm :
    name => inst.network[0].fixed_ip_v4
  }

  control_plane_hosts = {
    for n in local.control_plane_names :
    n => {
      ansible_host = local.node_ssh_ip[n]
      k8s_ip       = local.node_k8s_ip[n]
    }
  }

  worker_hosts = {
    for n in local.worker_names :
    n => {
      ansible_host = local.node_ssh_ip[n]
      k8s_ip       = local.node_k8s_ip[n]
    }
  }
}


resource "local_file" "ansible_inventory" {
  filename = "${path.module}/../ansible/inventory.ini"

  content = templatefile("${path.module}/templates/inventory.tftpl", {
    control_plane            = local.control_plane_hosts
    workers                  = local.worker_hosts
    ansible_user             = var.ansible_user
    ansible_private_key_file = var.ansible_private_key_file
  })
}

