data "openstack_compute_flavor_v2" "medium" {
  name = "m1.medium"
}

data "openstack_images_image_v2" "jammy" {
  name        = "jammy"
  most_recent = true
}