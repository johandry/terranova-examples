package main

import (
	"flag"
	"log"
	"os"
	"strconv"

	"github.com/hashicorp/terraform/builtin/provisioners/file"
	"github.com/johandry/terranova"
	"github.com/terraform-providers/terraform-provider-aws/aws"
)

var (
	code        string
	count       int
	pubKeyFile  string
	privKeyFile string
)

const stateFilename = "aws-ec2-ubuntu.tfstate"

func main() {
	flag.Parse()

	if count < 0 {
		log.Fatalf("count cannot be negative. It has to be '0' to terminate all the creted instances or the desired number of instances")
	}

	platform, err := terranova.NewPlatform(code).
		AddProvider("aws", aws.Provider()).
		AddProvisioner("file", file.Provisioner()).
		Var("count", strconv.Itoa(count)).
		ReadStateFromFile(stateFilename)

	if len(pubKeyFile) != 0 {
		platform.Var("public_key_file", pubKeyFile)
	}
	if len(privKeyFile) != 0 {
		platform.Var("private_key_file", privKeyFile)
	}

	if err != nil {
		if os.IsNotExist(err) {
			log.Printf("[DEBUG] state file %s does not exists", stateFilename)
		} else {
			log.Fatalf("Fail to load the initial state of the platform from file %s. %s", stateFilename, err)
		}
	}

	if err := platform.Apply((count == 0)); err != nil {
		log.Fatalf("Fail to apply the changes to the platform. %s", err)
	}

	if _, err := platform.WriteStateToFile(stateFilename); err != nil {
		log.Fatalf("Fail to save the final state of the platform to file %s. %s", stateFilename, err)
	}

	log.Printf("Check your EC2 instances with the AWS CLI command: `aws ec2 describe-instances --query 'Reservations[*].Instances[*].[InstanceId, PublicIpAddress, State.Name]' --output table`")
}

func init() {
	flag.IntVar(&count, "count", 2, "number of instances to create. Set to '0' to terminate them all.")
	flag.StringVar(&pubKeyFile, "pub", "", "public key file to create the AWS Key Pair")
	flag.StringVar(&privKeyFile, "priv", "", "private key file to connect to the new AWS EC2 instances")

	code = `
  variable "count" 	          { default = 2 }
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
    count         = "${var.count}"
    key_name      = "server_key"

    provisioner "file" {
      content     = "ami used: ${self.ami}"
      destination = "/tmp/file.log"

      connection {
        user        = "ubuntu"
        private_key = "${local.private_key}"
      }
    }
  }
  resource "aws_key_pair" "keypair" {
    key_name    = "server_key"
    public_key  = "${local.public_key}"
  }
`
}
