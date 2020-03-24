package main

import (
	"flag"
	"io"
	"io/ioutil"
	"log"
	"os"

	"github.com/hashicorp/terraform/builtin/provisioners/file"
	"github.com/johandry/terranova"
	"github.com/johandry/terranova/logger"
	"github.com/terraform-providers/terraform-provider-aws/aws"
)

var (
	code        string
	count       int
	pubKeyFile  string
	privKeyFile string
	debug       bool
	quiet       bool
)

const stateFilename = "aws-ec2-ubuntu.tfstate"

func main() {
	flag.Parse()

	// the standard log cannot be used because it's hijaked by the logger
	// Middleware. To use the logs, a new instance is required.
	log := log.New(os.Stderr, "", log.LstdFlags)

	if count < 0 {
		log.Fatalf("count cannot be negative. It has to be '0' to terminate all the creted instances or the desired number of instances")
	}

	if debug && quiet {
		log.Fatal("debug mode and quiet mode cannot be set at the same time")
	}

	log.Printf("Starting to (de)provision on AWS")

	logLevel := logger.LogLevelInfo
	if debug {
		logLevel = logger.LogLevelDebug
	}
	var output io.Writer = os.Stderr
	if quiet {
		output = ioutil.Discard
	}
	myLog := logger.NewLog(output, "EC2", logLevel)
	logMiddleware := logger.NewMiddleware(myLog)
	defer logMiddleware.Close()

	platform, err := terranova.NewPlatform(code).
		SetMiddleware(logMiddleware).
		AddProvider("aws", aws.Provider()).
		AddProvisioner("file", file.Provisioner()).
		Var("srv_count", count).
		PersistStateToFile(stateFilename)

	if len(pubKeyFile) != 0 {
		platform.Var("public_key_file", pubKeyFile)
	}
	if len(privKeyFile) != 0 {
		platform.Var("private_key_file", privKeyFile)
	}

	if err != nil {
		log.Fatalf("Fail to create the platform using state file %s. %s", stateFilename, err)
	}

	if err := platform.Apply((count == 0)); err != nil {
		log.Fatalf("Fail to apply the changes to the platform. %s", err)
	}

	log.Printf("Check your EC2 instances with the AWS CLI command: `aws ec2 describe-instances --query 'Reservations[*].Instances[*].[InstanceId, PublicIpAddress, State.Name]' --output table`")
}

func init() {
	flag.IntVar(&count, "count", 2, "number of instances to create. Set to '0' to terminate them all.")
	flag.StringVar(&pubKeyFile, "pub", "", "public key file to create the AWS Key Pair")
	flag.StringVar(&privKeyFile, "priv", "", "private key file to connect to the new AWS EC2 instances")
	flag.BoolVar(&debug, "debug", false, "debug mode, prints also debug output from terraform")
	flag.BoolVar(&quiet, "quiet", false, "quiet/silence mode, do not print any terraform output")

	code = `
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
`
}
