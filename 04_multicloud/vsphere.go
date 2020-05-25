package main

import (
	"log"
	"os"
	"strings"
)

func initVsphere() {
	vars = map[string]interface{}{
		"vsphere_username": os.Getenv("VSPHERE_USERNAME"),
		"vsphere_password": os.Getenv("VSPHERE_PASSWORD"),
		"vsphere_server":   os.Getenv("VSPHERE_SERVER"),
	}

	for k, v := range vars {
		if len(v.(string)) == 0 {
			log.Fatalf("Not found a value for %q, use the environment variable %s to set the value", k, strings.ToUpper(k))
		}
	}

	code = string(`
	variable "vsphere_username" {}
	variable "vsphere_password" {}
	variable "vsphere_server"   {}
	provider "vsphere" {
		user                  = "${var.vsphere_username}"
		password              = "${var.vsphere_password}"
		vsphere_server        = "${var.vsphere_server}"
		allow_unverified_ssl  = "true"
	}
	data "vsphere_datacenter" "dc" {
		name = "dc1"
	}
	data "vsphere_datastore" "datastore" {
		name          = "datastore1"
		datacenter_id = "${data.vsphere_datacenter.dc.id}"
	}
	data "vsphere_compute_cluster" "cluster" {
		name          = "cluster1"
		datacenter_id = "${data.vsphere_datacenter.dc.id}"
	}
	data "vsphere_network" "network" {
		name          = "public"
		datacenter_id = "${data.vsphere_datacenter.dc.id}"
	}
	resource "vsphere_virtual_machine" "vm" {
		count 					 = "${var.srv_count}"
		name             = "terranova-vm"
		resource_pool_id = "${data.vsphere_compute_cluster.cluster.resource_pool_id}"
		datastore_id     = "${data.vsphere_datastore.datastore.id}"
		num_cpus 				 = 2
		memory   				 = 1024
		guest_id 				 = "TerranovaLinux64Guest"
		network_interface {
			network_id 		 = "${data.vsphere_network.network.id}"
		}
		disk {
			label 				 = "disk0"
			size  				 = 20
		}
	}
	`)
}
