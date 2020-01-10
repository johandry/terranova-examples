package main

import (
	"log"
	"os"

	"github.com/johandry/terranova"
	mylog "github.com/johandry/terranova-examples/custom-logs/log"
	"github.com/johandry/terranova/logger"
	"github.com/terraform-providers/terraform-provider-aws/aws"
)

var code string

const stateFilename = "terractl.tfstate"

func main() {
	var logMiddleware *logger.Middleware

	count := 0
	logType := "custom" // "default", "terraform", "discard", "custom", "jlog-viper", "jlog-config", "logrus", "logrus-json"

	keyName := "demo"

	platform, err := terranova.NewPlatform(code).
		AddProvider("aws", aws.Provider()).
		Var("c", count).
		Var("key_name", keyName).
		ReadStateFromFile(stateFilename)

	switch logType {
	case "custom":
		logMiddleware = mylog.Custom()
	case "jlog-viper":
		logMiddleware = mylog.JLogViper()
	case "jlog-config":
		logMiddleware = mylog.JLogConfig()
	case "logrus-json":
		logMiddleware = mylog.LogrusJSON(count)
	case "logrus":
		logMiddleware = mylog.Logrus(count)
	case "discard":
		logMiddleware = mylog.Discard()
	case "terraform":
		logMiddleware = mylog.Terraform()
	case "default":
		logMiddleware = mylog.Default()
	}

	if logMiddleware != nil {
		platform.SetMiddleware(logMiddleware)
	}

	log := log.New(os.Stderr, "TERRACTL: ", log.LstdFlags)

	if err != nil {
		if os.IsNotExist(err) {
			log.Printf("[DEBUG] state file %s does not exists", stateFilename)
		} else {
			log.Fatalf("Fail to load the initial state of the platform from file %s. %s", stateFilename, err)
		}
	}

	log.Print("Begining of (de)provisioning of AWS instances")
	terminate := (count == 0)
	if err := platform.Apply(terminate); err != nil {
		log.Fatalf("Fail to apply the changes to the platform. %s", err)
	}

	if _, err := platform.WriteStateToFile(stateFilename); err != nil {
		log.Fatalf("Fail to save the final state of the platform to file %s. %s", stateFilename, err)
	}
	log.Print("End of (de)provisioning of AWS instances")
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
