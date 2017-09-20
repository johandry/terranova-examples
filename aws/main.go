package main

import (
	"flag"
	"log"
	"os"
	"strconv"

	"github.com/johandry/platformer"
	"github.com/terraform-providers/terraform-provider-aws/aws"
)

const defultStateFilename = "aws-ec2-ubuntu.tfstate"

var (
	code          string
	terminate     bool
	stateFilename string
	count         int
)

func main() {
	flag.Parse()

	p, err := platformer.New("", code)
	if err != nil {
		log.Fatalf("Fail to initialize platformer. %s", err)
	}

	p.AddProvider("aws", aws.Provider())

	p.Var("count", strconv.Itoa(count))

	if p.ExistStateFile(stateFilename) {
		if _, err := p.ReadStateFile(stateFilename); err != nil {
			log.Panicf("Fail to load the initial state of the platform from file %s. %s", stateFilename, err)
		}
	}

	if terminate {
		if err := p.Destroy(); err != nil {
			log.Fatalf("Fail to destroy the platform. %s", err)
		}
	} else {
		if err := p.Create(); err != nil {
			log.Fatalf("Fail to create the platform. %s", err)
		}
	}

	if err := p.WriteStateFile(stateFilename); err != nil {
		log.Fatalf("Fail to save the final state of the platform to file %s. %s", stateFilename, err)
	}
}

func init() {
	flag.BoolVar(&terminate, "terminate", false, "Terminate the EC2 instance")
	flag.BoolVar(&terminate, "t", false, "Terminate the EC2 instance")
	flag.StringVar(&stateFilename, "state", defultStateFilename, "Filename to store the state")

	countStr := os.Getenv("PLATFORMER_AWS_COUNT")
	if len(countStr) > 0 {
		var err error
		count, err = strconv.Atoi(countStr)
		if err != nil {
			log.Printf("Could not convert to int the value of 'PLATFORMER_AWS_COUNT' (%s). %s", countStr, err)
		}
	}

	code = `
  variable "aws_region" {
    description = "The AWS region to create things in."
    default     = "us-west-2"
  }

	variable "count" {
		description = "Total instances to create"
		default 		= 2
	}

  # Ubuntu Precise 12.04 LTS (x64)
  variable "aws_amis" {
    default = {
      "eu-west-1" = "ami-cd0f5cb6"
      "us-east-1" = "ami-10547475"
      "us-west-1" = "ami-09d2fb69"
      "us-west-2" = "ami-6e1a0117"
    }
  }
  provider "aws" {
    region = "${var.aws_region}"
  }
  resource "aws_instance" "web" {
    instance_type = "t2.micro"
    ami           = "${lookup(var.aws_amis, var.aws_region)}"

    # This will create 2 instances
    count = ${var.count}
  }`
}
