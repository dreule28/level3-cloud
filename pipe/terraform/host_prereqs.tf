# resource "null_resource" "openstack_host_prereqs" {
#   triggers = {
#     sysctl_conf = sha1(<<-EOF
#       net.ipv4.conf.all.rp_filter=0
#       net.ipv4.conf.default.rp_filter=0
#       net.ipv4.conf.br-ex.rp_filter=0
#     EOF
#     )
#   }

#   provisioner "local-exec" {
#     interpreter = ["/bin/bash", "-lc"]
#     command = <<-BASH
#       set -euo pipefail

#       # Persist across reboots
#       sudo tee /etc/sysctl.d/99-openstack-rpf.conf >/dev/null <<EOF
# net.ipv4.conf.all.rp_filter=0
# net.ipv4.conf.default.rp_filter=0
# net.ipv4.conf.br-ex.rp_filter=0
# EOF

#       # Apply now
#       sudo sysctl --system >/dev/null

#       # Log current values (helps when reading terraform output)
#       sysctl net.ipv4.conf.all.rp_filter net.ipv4.conf.default.rp_filter net.ipv4.conf.br-ex.rp_filter
#     BASH
#   }
# }
