package data

// TerraformIBMTmpl is the Terraform code in form of Go template to be applied by
// Terranova after rendere it
const TerraformIBMTmpl = `
{{ if .SSHAccess }}
variable "public_key_file" { }
locals {
  public_key = file(pathexpand(var.public_key_file))
}
{{- end }}

locals {
  port       = {{ .LetsChatPort }}
  prefix     = "terranova-terranova"
}

output "ip_address" {
  value = ibm_is_floating_ip.demo_terranova_floating_ip.address
}

output "entrypoint" {
  value = "http://${ibm_is_floating_ip.demo_terranova_floating_ip.address}:${local.port}/"
}

provider "ibm" {
  generation = 2
  region     = "us-south"
}

resource "ibm_is_vpc" "demo_terranova_vpc" {
  name = "${local.prefix}-vpc"
  tags = [
    "project:demo-terranova",
    "owner:johandry",
  ]
}

resource "ibm_is_subnet" "demo_terranova_subnet" {
  name            = "${local.prefix}-subnet"
  vpc             = ibm_is_vpc.demo_terranova_vpc.id
  zone            = "us-south-1"
  ipv4_cidr_block = "10.240.0.0/24"
}

resource "ibm_is_security_group" "demo_terranova_security_group" {
  name = "${local.prefix}-sg-public"
  vpc  = ibm_is_vpc.demo_terranova_vpc.id
}

resource "ibm_is_security_group_rule" "demo_terranova_security_group_rule_all_outbound" {
  group     = ibm_is_security_group.demo_terranova_security_group.id
  direction = "outbound"
}

resource "ibm_is_security_group_rule" "demo_terranova_security_group_rule_tcp_http" {
  group     = ibm_is_security_group.demo_terranova_security_group.id
  direction = "inbound"
  tcp {
    port_min = local.port
    port_max = local.port
  }
}

{{ if .SSHAccess }}
resource "ibm_is_security_group_rule" "demo_terranova_security_group_rule_tcp_ssh" {
  group     = ibm_is_security_group.demo_terranova_security_group.id
  direction = "inbound"
  tcp {
    port_min = 22
    port_max = 22
  }
}
{{- end }}

resource "ibm_is_floating_ip" "demo_terranova_floating_ip" {
  name   = "${local.prefix}-ip"
  target = ibm_is_instance.demo_terranova_instance.primary_network_interface.0.id
  tags = [
    "project:demo-terranova",
    "owner:johandry",
  ]
}

{{ if .SSHAccess }}
resource "ibm_is_ssh_key" "demo_terranova_key" {
  name       = "${local.prefix}-key"
  public_key = local.public_key
  tags = [
    "project:demo-terranova",
    "owner:johandry",
  ]
}
{{- end }}

resource "ibm_is_instance" "demo_terranova_instance" {
  name    = "${local.prefix}-instance"
  image   = "r006-14140f94-fcc4-11e9-96e7-a72723715315"
  profile = "cx2-2x4"

  primary_network_interface {
    name            = "eth1"
    subnet          = ibm_is_subnet.demo_terranova_subnet.id
    security_groups = [ibm_is_security_group.demo_terranova_security_group.id]
  }

  vpc  = ibm_is_vpc.demo_terranova_vpc.id
  zone = "us-south-1"
  keys = [{{ if .SSHAccess }}ibm_is_ssh_key.demo_terranova_key.id{{- end }}]

  user_data = <<-EOUD
              #!/bin/bash
              echo "Hello World" > index.html
              nohup busybox httpd -f -p ${local.port} &
              EOUD

  tags = [
    "project:demo-terranova",
    "owner:johandry",
  ]
}
`
