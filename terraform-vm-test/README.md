# Terraform VM test (OpenStack)

This folder contains a minimal Terraform configuration that provisions a single OpenStack VM.

## What the files do

- `main.tf`: Pins the Terraform/OpenStack provider versions and declares the `openstack` provider (credentials/region are expected to come from your environment, e.g. `OS_*` variables or a `clouds.yaml`).
- `probe.tf`: Looks up existing OpenStack resources:
	- Flavor `m1.small`
	- Image `cirros-0.6.3-x86_64-disk` (most recent match)
- `instance.tf`: Creates one compute instance named `tf-vm1` using the looked-up flavor/image, attaches it to the `private` network, and uses an existing keypair named `openstack-key`.

## Result

Running `terraform apply` will create one VM (`openstack_compute_instance_v2.vm1`).
