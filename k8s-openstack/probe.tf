data "openstack_compute_flavor_v2" "small" {
  name = "m1.small"
}

data "openstack_images_image_v2" "bookworm" {
  name        = "bookworm"
  most_recent = true
}