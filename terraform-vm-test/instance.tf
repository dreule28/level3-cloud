resource "openstack_compute_instance_v2" "vm1" {
  name      = "tf-vm1"
  flavor_id = data.openstack_compute_flavor_v2.small.id
  image_id  = data.openstack_images_image_v2.ubuntu.id
  key_pair  = "openstack-key"

  network {
    name = "private"
  }
}