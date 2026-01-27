resource "openstack_compute_instance_v2" "server" {
  name      = "server"
  flavor_id = data.openstack_compute_flavor_v2.small.id
  image_id  = data.openstack_images_image_v2.bookworm.id
  key_pair  = "openstack-key"

  network {
    name = "shared"
  }
}

resource "openstack_compute_instance_v2" "node-0" {
  name      = "node-0"
  flavor_id = data.openstack_compute_flavor_v2.small.id
  image_id  = data.openstack_images_image_v2.bookworm.id
  key_pair  = "openstack-key"

  network {
    name = "shared"
  }
}

resource "openstack_compute_instance_v2" "node-1" {
  name      = "node-1"
  flavor_id = data.openstack_compute_flavor_v2.small.id
  image_id  = data.openstack_images_image_v2.bookworm.id
  key_pair  = "openstack-key"

  network {
    name = "shared"
  }
}
