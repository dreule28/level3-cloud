resource "openstack_compute_instance_v2" "ControlPlane" {
  name      = "controlplane"
  flavor_id = data.openstack_compute_flavor_v2.medium.id
  image_id  = data.openstack_images_image_v2.jammy.id
  key_pair  = "openstack-key"

  network {
    name = "shared"
  }
}

resource "openstack_compute_instance_v2" "Worker" {
  name      = "worker"
  flavor_id = data.openstack_compute_flavor_v2.medium.id
  image_id  = data.openstack_images_image_v2.jammy.id
  key_pair  = "openstack-key"

  network {
    name = "shared"
  }
}
