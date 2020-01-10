package main

import (
	"log"
	"os"

	"github.com/johandry/terranova"
	"github.com/johandry/terranova/logger"
	"github.com/terraform-providers/terraform-provider-aws/aws"
)

var code string

const stateFilename = "simple.tfstate"

func main() {
	count := 0
	keyName := "demo"

	log := log.New(os.Stderr, "", log.LstdFlags)
	logMiddleware := logger.NewMiddleware()
	defer logMiddleware.Close()

	platform, err := terranova.NewPlatform(code).
		SetMiddleware(logMiddleware).
		AddProvider("aws", aws.Provider()).
		Var("c", count).
		Var("key_name", keyName).
		PersistStateToFile(stateFilename)

	if err != nil {
		log.Fatalf("Fail to create the platform using state file %s. %s", stateFilename, err)
	}

	terminate := (count == 0)
	if err := platform.Apply(terminate); err != nil {
		log.Fatalf("Fail to apply the changes to the platform. %s", err)
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
