resource "openstack_compute_keypair_v2" "this" {
  name       = var.keypair_name
  public_key = file(pathexpand(var.public_key_path))
}
