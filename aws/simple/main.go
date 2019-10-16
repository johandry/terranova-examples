package main

import (
	"log"
	"os"

	"github.com/johandry/terranova"
	"github.com/terraform-providers/terraform-provider-aws/aws"
)

var code string

const stateFilename = "simple.tfstate"

func main() {
	count := 1
	keyName := "demo"

	platform, err := terranova.NewPlatform(code).
		AddProvider("aws", aws.Provider()).
		Var("c", count).
		Var("key_name", keyName).
		ReadStateFromFile(stateFilename)

	if err != nil {
		if os.IsNotExist(err) {
			log.Printf("[DEBUG] state file %s does not exists", stateFilename)
		} else {
			log.Fatalf("Fail to load the initial state of the platform from file %s. %s", stateFilename, err)
		}
	}

	terminate := (count == 0)
	if err := platform.Apply(terminate); err != nil {
		log.Fatalf("Fail to apply the changes to the platform. %s", err)
	}

	if _, err := platform.WriteStateToFile(stateFilename); err != nil {
		log.Fatalf("Fail to save the final state of the platform to file %s. %s", stateFilename, err)
	}
}

func init() {
	code = `
  variable "c"    { default = 2 }
  variable "key_name" {}
  provider "aws" {
    region        = "us-west-2"
  }
  resource "aws_instance" "server" {
    instance_type = "t2.micro"
    ami           = "ami-6e1a0117"
    count         = "${var.c}"
    key_name      = "${var.key_name}"
  }
`
}
