locals {
  bgp_password = "955dB0b81Ef"
  project_id   = "<UUID_of_your_project>"
}

# you need to enable BGP config for the project. If you decide to create new
# project, you can use the bgp_config section to enable BGP.
# resource "equinix_metal_project" "test" {
#   name = "testpro"
#   bgp_config {
#      deployment_type = "local"
#      md5 = local.bgp_password
#      asn = 65000
#   }
# }

resource "equinix_metal_reserved_ip_block" "addr" {
  project_id = local.project_id
  metro      = "ny"
  quantity   = 1
}

resource "equinix_metal_device" "test" {
  hostname         = "terraform-test-bgp-sesh"
  plan             = "c3.small.x86"
  metro            = ["ny"]
  operating_system = "ubuntu_20_04"
  billing_cycle    = "hourly"
  project_id       = local.project_id
}

resource "equinix_metal_bgp_session" "test" {
  device_id      = equinix_metal_device.test.id
  address_family = "ipv4"
}


data "template_file" "interface_lo0" {
  template = <<EOF
auto lo:0
iface lo:0 inet static
   address $${floating_ip}
   netmask $${floating_netmask}
EOF

  vars = {
    floating_ip      = equinix_metal_reserved_ip_block.addr.address
    floating_netmask = equinix_metal_reserved_ip_block.addr.netmask
  }
}

data "template_file" "bird_conf_template" {

  template = <<EOF
filter equinix_metal_bgp {
    if net = $${floating_ip}/$${floating_cidr} then accept;
}
router id $${private_ipv4};
protocol direct {
    interface "lo";
}
protocol kernel {
    scan time 10;
    persist;
    import all;
    export all;
}
protocol device {
    scan time 10;
}
protocol bgp {
    export filter equinix_metal_bgp;
    local as 65000;
    neighbor $${gateway_ip} as 65530;
    password "$${bgp_password;
}
EOF

  vars = {
    floating_ip   = equinix_metal_reserved_ip_block.addr.address
    floating_cidr = equinix_metal_reserved_ip_block.addr.cidr
    private_ipv4  = equinix_metal_device.test.network.2.address
    gateway_ip    = equinix_metal_device.test.network.2.gateway
    bgp_password  = local.bgp_password
  }
}

resource "null_resource" "configure_bird" {

  connection {
    type        = "ssh"
    host        = equinix_metal_device.test.access_public_ipv4
    private_key = file("/home/tomk/keys/tkarasek_key.pem")
    agent       = false
  }

  provisioner "remote-exec" {
    inline = [
      "apt-get install bird",
      "mv /etc/bird/bird.conf /etc/bird/bird.conf.old",
    ]
  }

  triggers = {
    template = data.template_file.bird_conf_template.rendered
    template = data.template_file.interface_lo0.rendered
  }

  provisioner "file" {
    content     = data.template_file.bird_conf_template.rendered
    destination = "/etc/bird/bird.conf"
  }

  provisioner "file" {
    content     = data.template_file.interface_lo0.rendered
    destination = "/etc/network/interfaces.d/lo0"
  }

  provisioner "remote-exec" {
    inline = [
      "sysctl net.ipv4.ip_forward=1",
      "grep /etc/network/interfaces.d /etc/network/interfaces || echo 'source /etc/network/interfaces.d/*' >> /etc/network/interfaces",
      "ifup lo:0",
      "service bird restart",
    ]
  }
}
