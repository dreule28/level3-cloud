data "openstack_compute_flavor_v2" "small" {
  name = "m1.small"
}

resource "openstack_images_image_v2" "bookworm" {
  name             = "bookworm"
  image_source_url = "https://cloud.debian.org/images/cloud/bookworm/latest/debian-12-genericcloud-amd64.qcow2"

  container_format = "bare"
  disk_format      = "qcow2"

}
