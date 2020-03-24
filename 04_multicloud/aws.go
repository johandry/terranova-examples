package main

import "os"

func initAWS() {
	vars = map[string]interface{}{}

	if v := os.Getenv("PUB_KEY_FILE"); len(v) != 0 {
		vars["public_key_file"] = v
	}
	if v := os.Getenv("PRIV_KEY_FILE"); len(v) != 0 {
		vars["private_key_file"] = v
	}

	code = string(`
  variable "srv_count" 	      { default = 2 }
  variable "public_key_file"  { default = "~/.ssh/id_rsa.pub" }
  variable "private_key_file" { default = "~/.ssh/id_rsa" }
  locals {
    public_key    = "${file(pathexpand(var.public_key_file))}"
    private_key   = "${file(pathexpand(var.private_key_file))}"
  }
  provider "aws" {
    region        = "us-west-2"
  }
  resource "aws_instance" "server" {
    instance_type = "t2.micro"
    ami           = "ami-6e1a0117"
    count         = "${var.srv_count}"
    key_name      = "server_key"

    provisioner "file" {
      content     = "ami used: ${self.ami}"
      destination = "/tmp/file.log"

      connection {
        user        = "ubuntu"
				private_key = "${local.private_key}"
				host 				= "${self.public_ip}"
      }
    }
  }
  resource "aws_key_pair" "keypair" {
    key_name    = "server_key"
    public_key  = "${local.public_key}"
  }
	`)
}
