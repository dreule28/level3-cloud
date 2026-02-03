data "openstack_compute_flavor_v2" "small" {
  name = "m1.small"
}

data "openstack_compute_flavor_v2" "medium" {
  name = "m1.medium"
}

resource "openstack_images_image_v2" "jammy" {
  name             = "jammy"
  image_source_url = "https://cloud-images.ubuntu.com/jammy/current/jammy-server-cloudimg-amd64.img"

  container_format = "bare"
  disk_format      = "qcow2"

}
