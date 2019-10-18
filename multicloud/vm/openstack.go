package main

import (
	"log"
	"os"
	"strings"
)

func initOpenStack() {
	vars = map[string]interface{}{
		"openstack_username":   os.Getenv("OPENSTACK_USERNAME"),
		"openstack_tenantname": os.Getenv("OPENSTACK_TENANTNAME"),
		"openstack_password":   os.Getenv("OPENSTACK_PASSWORD"),
		"openstack_authurl":    os.Getenv("OPENSTACK_AUTHURL"),
		"openstack_domainname": os.Getenv("OPENSTACK_DOMAINNAME"),
	}

	for k, v := range vars {
		if len(v.(string)) == 0 {
			log.Fatalf("Not found a value for %q, use the environment variable %s to set the value", k, strings.ToUpper(k))
		}
	}

	if v := os.Getenv("OPENSTACK_REGION"); len(v) != 0 {
		vars["openstack_region"] = v
	}

	code = string(`
	variable "openstack_tenantname" 	{}
	variable "openstack_authurl" 			{}
	variable "openstack_username"   	{}
	variable "openstack_password"   	{}
	variable "openstack_domainname"  	{}
	variable "openstack_region"   		{ default = "RegionOne" }
	variable "public_key_file"  { default = "~/.ssh/id_rsa.pub" }
  variable "private_key_file" { default = "~/.ssh/id_rsa" }
  locals {
    public_key    = "${file(pathexpand(var.public_key_file))}"
    private_key   = "${file(pathexpand(var.private_key_file))}"
  }
	provider "openstack" {
		user_name     = "${var.openstack_username}"
		tenant_name   = "${var.openstack_tenantname}"
		password      = "${var.openstack_password}"
		auth_url      = "${var.openstack_authurl}"
		domain_name   = "${var.openstack_domainname}"
		region      	= "${var.openstack_region}"
	}
	resource "openstack_compute_keypair_v2" "keypair" {
		region     = "${var.openstack_region}"
		name       = "terranova-keypair"
		public_key = "${local.public_key}\n"
	}
	resource "openstack_compute_instance_v2" "instance" {
		name            = "terranova_instance"
		image_id        = "TerranovaLinux64ImageId"
		flavor_id       = "3"
		key_pair        = "${openstack_compute_keypair_v2.keypair.id}"
		security_groups = ["default"]
		block_device {
			uuid                  = "TerranovaLinux64ImageId"
			source_type           = "image"
			destination_type      = "local"
			boot_index            = 0
			delete_on_termination = true
		}
		block_device {
			source_type           = "blank"
			destination_type      = "volume"
			volume_size           = 1
			boot_index            = 1
			delete_on_termination = true
		}
	}
	`)
}
